using Grpc.Core;
using Infrastructure.Service;
using ToolCluster.V4;


namespace SshExample.Services
{
    public class LinuxManagerService : LinuxManager.LinuxManagerBase
    {
        public override Task<CollectLinuxCpuUtilizationResponse> CollectCpuUtilization(CollectLinuxCpuUtilizationRequest request, ServerCallContext context)
        {
            return Task.FromResult(SshCommandCaller.GetCpuUtilisation(request));
        }
    }
}