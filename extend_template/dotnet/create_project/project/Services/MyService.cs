using Grpc.Core;

namespace SnmpExample.Services
{
    public class MyService : MyProtoService.MyProtoServiceBase
    {
        public override Task<PingReply> SendPing(PingRequest request, ServerCallContext context)
        {
            return Task.FromResult(new PingReply { RespondString = "Respond string"});
        }
    }
}
