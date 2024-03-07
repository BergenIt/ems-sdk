using Renci.SshNet;
using ToolCluster.V4;

namespace Infrastructure.Service;

public class SshCommandCaller
{
    private static readonly string[] _commands = { "top -b -n1 -1 -p0 -w 400", "top -b -n1 -w 400" };
    private static string TopCpuPrefix = "%cpu";

    public record HandleResult(int Success, string Stdout, string Stderr);
    public static CollectLinuxCpuUtilizationResponse GetCpuUtilisation(CollectLinuxCpuUtilizationRequest request)
    {
        HandleResult[] responses = _commands.Select(cmd => CallSsh(request, cmd)).ToArray();
        return ProcessResponse(responses, request);
    }
    public static HandleResult CallSsh(CollectLinuxCpuUtilizationRequest request, string command)
    {
        DeviceConnector connection = request.Device.Connectors.First();
        Credential creds = connection.Credentials.First();

        using (SshClient client = new SshClient(connection.Address, creds.Port, creds.Login, creds.Password))
        {
            client.Connect();
            SshCommand cmdRes = client.RunCommand(command);
            if (cmdRes.ExitStatus == 0)
            {
                client.Disconnect();
                return new HandleResult(cmdRes.ExitStatus,
                cmdRes.Result, cmdRes.Error);
            }
            client.Disconnect();
            return new HandleResult(cmdRes.ExitStatus,
                cmdRes.Result, cmdRes.Error);
        }
    }
    public static CollectLinuxCpuUtilizationResponse ProcessResponse(HandleResult[] results, CollectLinuxCpuUtilizationRequest request)
    {
        string response = "";
        foreach (HandleResult res in results)
        {
            if (res.Success == 0)
            {
                response += res.Stdout;
            }
        }

        CollectLinuxCpuUtilizationResponse statDeviceCpu = new()
        {
            CpuUtilization = new()
            {
                DeviceIdentity = new DeviceDataIdentity()
                {
                    DeviceId = request.Device.DeviceId,
                    Source = ServiceSource.LinuxManager
                },
                SummaryUtilization = new()
            }
        };

        IEnumerable<string> rows = response.ToLowerInvariant()
            .Split('\n', StringSplitOptions.TrimEntries)
            .Where(s => s.StartsWith(TopCpuPrefix))
            .SelectMany(c => c
                .Split(TopCpuPrefix, StringSplitOptions.TrimEntries)
                .Where(d => !string.IsNullOrWhiteSpace(d))
                .Select(s => TopCpuPrefix + s));

        foreach (string item in rows)
        {
            int processorIdEnd = item.IndexOf(' ', TopCpuPrefix.Length);

            if (processorIdEnd == -1)
            {
                continue;
            }

            string? strProcessorId = item[TopCpuPrefix.Length..processorIdEnd];

            int processorId;

            if (strProcessorId == "(s):")
            {
                processorId = -1;
            }
            else if (!int.TryParse(strProcessorId, out processorId))
            {
                continue;
            }
            CpuUnitUtilization cpu = new();

            foreach (string entiry in item[(processorIdEnd + 3)..].Split(',', StringSplitOptions.TrimEntries))
            {
                string[] parts = entiry.Split(' ', 2, StringSplitOptions.TrimEntries);

                if (parts.Length != 2)
                {
                    continue;
                }

                string key = parts[1];
                string strValue = parts[0];

                if (float.TryParse(strValue, out float intValue))
                {
                    switch (key)
                    {
                        case "us": cpu.UserUsing = (int)Math.Ceiling(intValue); break;
                        case "sy": cpu.SystemUsing = (int)Math.Ceiling(intValue); break;
                        case "ni": cpu.NiceValueUsing = (int)Math.Ceiling(intValue); break;
                        case "id": cpu.IdleTime = (int)Math.Ceiling(intValue); break;
                        case "wa": cpu.IoWaiting = (int)Math.Ceiling(intValue); break;
                        case "hi": cpu.HwServiceInterrupts = (int)Math.Ceiling(intValue); break;
                        case "si": cpu.SoftServiceInterrupts = (int)Math.Ceiling(intValue); break;
                        case "st": cpu.StealTime = (int)Math.Ceiling(intValue); break;
                        default: break;
                    }
                }
               
            }
            if (processorId == -1)
            {
                statDeviceCpu.CpuUtilization.SummaryUtilization = cpu;
            }
            else
            {
                statDeviceCpu.CpuUtilization.UnitUtilistaions.Add(processorId, cpu);
            }
        }
        return statDeviceCpu;
    }

}
