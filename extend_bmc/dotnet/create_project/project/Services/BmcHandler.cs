using Grpc.Core;

using ToolCluster.V4;

namespace BmcHandler.Services
{
    public class BmcHandlerService : ToolCluster.V4.BmcManager.BmcManagerBase
    {
        public override async Task<CollectBmcLedStateResponse> CollectLedState(CollectBmcLedStateRequest request, ServerCallContext context)
        {
            throw new RpcException(new Status(StatusCode.Unimplemented, ""));
        }

        public override async Task<BmcFirmwareUpdateResponse> BmcFirmwareUpdate(BmcFirmwareUpdateRequest request, ServerCallContext context)
        {
            throw new RpcException(new Status(StatusCode.Unimplemented, ""));
        }
    }
}
