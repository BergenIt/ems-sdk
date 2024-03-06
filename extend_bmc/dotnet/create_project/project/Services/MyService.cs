using Grpc.Core;

namespace BmcManager.Services
{
    public class MyService : MyProtoService.MyProtoServiceBase
    {
        public override Task<PingReply> SendPing(PingRequest request, ServerCallContext context)
        {
            return Task.FromResult(new PingReply { RespondString = "Respond string"});
        }
    }
}
