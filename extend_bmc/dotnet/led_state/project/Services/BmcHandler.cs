using Grpc.Core;

using ToolCluster.V4;

namespace BmcHandler.Services
{
    public class BmcHandlerService : BmcManager.BmcManagerBase
    {
        // RPC по сбору статуса LED.
        public override async Task<CollectBmcLedStateResponse> CollectLedState(CollectBmcLedStateRequest request, ServerCallContext context)
        {
            throw new RpcException(new Status(StatusCode.Unimplemented, ""));
        }
    }
}
