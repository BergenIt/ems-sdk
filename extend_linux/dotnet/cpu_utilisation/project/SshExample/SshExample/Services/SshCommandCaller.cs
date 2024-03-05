using Renci.SshNet;
using ToolCluster.V4;

namespace Infrastructure.Service;

public class SshCommandCaller 
{
    private static readonly string[] _commands = { "top -b -n1 -1 -p0 -w 400", "top -b -n1 -w 400" };
    private static string TopCpuPrefix = "%cpu";

    public record HandleResult(int Success, string Stdout, string Stderr);
    public static CollectLinuxCpuUtilisationResponse GetCpuUtilisation(CollectLinuxCpuUtilisationRequest request)
    {
        HandleResult[] responses= _commands.Select(cmd=> CallSsh(request,cmd)).ToArray();
        return ProcessResponse(responses, request);
    }
    public static HandleResult CallSsh(CollectLinuxCpuUtilisationRequest request, string command)
    {
        DeviceConnector connection = request.Device.Connectors.First();
        Credential creds=connection.Credentials.First();

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
    public static CollectLinuxCpuUtilisationResponse ProcessResponse(HandleResult[] results, CollectLinuxCpuUtilisationRequest request)
    {
        string response = "";
        foreach (HandleResult res in results)
        {
            if (res.Success == 0)
            {
                response += res.Stdout;
            }
        }

        CollectLinuxCpuUtilisationResponse statDeviceCpu = new()
        {
            CpuUtilisation = new()
            {
                DeviceIdentity = new DeviceDataIdentity()
                {
                    DeviceId = request.Device.DeviceId,
                    Source = ServiceSource.LinuxManager
                },
                SummaryUtilisation =new()
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
                        case "us": statDeviceCpu.CpuUtilisation.SummaryUtilisation.UserUsing = (int)Math.Ceiling(intValue); break;
                        case "sy": statDeviceCpu.CpuUtilisation.SummaryUtilisation.SystemUsing = (int)Math.Ceiling(intValue); break;
                        case "ni": statDeviceCpu.CpuUtilisation.SummaryUtilisation.NiceValueUsing = (int)Math.Ceiling(intValue); break;
                        case "id": statDeviceCpu.CpuUtilisation.SummaryUtilisation.IdleTime = (int)Math.Ceiling(intValue); break;
                        case "wa": statDeviceCpu.CpuUtilisation.SummaryUtilisation.IoWaiting = (int)Math.Ceiling(intValue); break;
                        case "hi": statDeviceCpu.CpuUtilisation.SummaryUtilisation.HwServiceInterrupts = (int)Math.Ceiling(intValue); break;
                        case "si": statDeviceCpu.CpuUtilisation.SummaryUtilisation.SoftServiceInterrupts = (int)Math.Ceiling(intValue); break;
                        case "st": statDeviceCpu.CpuUtilisation.SummaryUtilisation.StealTime = (int)Math.Ceiling(intValue); break;
                        default: break;
                    }
                }
            }
        }
        return statDeviceCpu;
    }

}
