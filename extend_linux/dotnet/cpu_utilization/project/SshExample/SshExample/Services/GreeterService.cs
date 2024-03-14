using Grpc.Core;
using Infrastructure.Service;
using ToolCluster.V4;


namespace SshExample.Services
{
    public class LinuxManagerService : LinuxManager.LinuxManagerBase
    {
        public override Task<CollectLinuxCpuUtilizationResponse> CollectCpuUtilization(CollectLinuxCpuUtilizationRequest request, ServerCallContext context)
        {
            Console.WriteLine("cpu request");
            Console.WriteLine("------");
            Console.WriteLine(request.Device.DeviceId);
            Console.WriteLine("++++++");
            return Task.FromResult(SshCommandCaller.GetCpuUtilisation(request));
        }
    }
}
