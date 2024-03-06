using Grpc.Core;

using Newtonsoft.Json;
using Newtonsoft.Json.Linq;

using System.Collections.Immutable;
using System.Text;

using ToolCluster.V4;

namespace BmcHandler.Services
{
    public class BmcHandlerService : BmcManager.BmcManagerBase
    {
        private const string OffTag = "Off";

        // RPC по сбору статуса LED.
        public override async Task<CollectBmcLedStateResponse> CollectLedState(CollectBmcLedStateRequest request, ServerCallContext context)
        {
            if (request.Device is null)
            {
                return new();
            }

            CollectBmcLedStateResponse resp = new()
            {
                Led = new()
                {
                    DeviceIdentity = new()
                    {
                        DeviceId = request.Device.DeviceId,
                        Source = ServiceSource.BmcManager,
                    },
                    State = LedState.Unknown,
                }
            };

            foreach (DeviceConnector connector in request.Device.Connectors)
            {
                Credential? credential = connector.Credentials.FirstOrDefault(c => c.Protocol is ConnectorProtocol.Redfish);
                if (credential is null)
                {
                    continue;
                }

                resp.Led.State = await GetRedfishLedState(credential, connector.Address);
                if (resp.Led.State is not LedState.Unspecified && resp.Led.State is not LedState.Unknown)
                {
                    break;
                }
            }

            return resp;
        }

        private static async Task<LedState> GetRedfishLedState(Credential creds, string address)
        {
            string payloadHeader = $"authorization:Basic {Convert.ToBase64String(Encoding.UTF8.GetBytes(creds.Login + ':' + creds.Password))}";

            ImmutableArray<string> systemMembers = await GetMembers("/redfish/v1/Chassis/", payloadHeader, address, creds.Port);
            if (!systemMembers.Any())
            {
                return LedState.Unknown;
            }

            string path = systemMembers[0];

            string? ledState = await TryGetStringPropertyFromPage(path, "IndicatorLED", payloadHeader, address, creds.Port);

            return GetLedState(ledState ?? string.Empty);
        }

        private static LedState GetLedState(string strLedState)
        {
            return strLedState switch
            {
                "Blinking" => LedState.Blink,
                "Lit" => LedState.On,
                "Off" => LedState.Off,
                _ => LedState.Unknown,
            };
        }

        private static async Task<ImmutableArray<string>> GetMembers(string path, string authHeader, string ip, int port)
        {
            string? data = await TryGetJsonPage(path, authHeader, ip, port);

            if (data is null)
            {
                return [];
            }

            MemberList? memberList = default;

            try
            {
                memberList = JsonConvert.DeserializeObject<MemberList>(data);
            }
            catch
            {
                return [];
            }

            ImmutableArray<string> result = memberList
                ?.Members
                ?.Select(m => m.OdataId)
                .Where(d => !string.IsNullOrWhiteSpace(d))
                .OfType<string>()
                .ToImmutableArray() ?? [];

            return result;
        }

        private static async Task<string?> TryGetJsonPage(string path, string authHeader, string ip, int port)
        {
            try
            {
                string uri = (port == 80 || port / 100 == 50 ? "http" : "https") + "://" + ip + ':' + port + path;

                HttpRequestMessage httpRequestMessage = new()
                {
                    RequestUri = new(uri),
                };

                string[] headers = authHeader.Split(':', 2);

                if (headers.Length == 2)
                {
                    httpRequestMessage.Headers.Add(headers[0], headers[1]);
                }

                HttpClient httpClient = new(new HttpClientHandler()
                {
                    ServerCertificateCustomValidationCallback = (message, cert, chain, errors) => true,
                    AllowAutoRedirect = false,
                })
                {
                    Timeout = TimeSpan.FromSeconds(100)
                };

                using HttpResponseMessage httpResponseMessage = await httpClient.SendAsync(httpRequestMessage);

                string result = await httpResponseMessage.Content.ReadAsStringAsync();

                if (httpResponseMessage.IsSuccessStatusCode)
                {
                    return result;
                }
            }
            catch
            {
                return null;
            }

            return null;
        }

        private static async Task<string?> TryGetStringPropertyFromPage(string path, string propertyName, string authHeader, string ip, int port)
        {
            string? data = await TryGetJsonPage(path, authHeader, ip, port);

            if (data is null)
            {
                return null;
            }

            try
            {
                string? value = null;
                string[] props = propertyName.Split("__");

                foreach (string p in props)
                {
                    JObject page = JObject.Parse(data);
                    JToken? obj = page.GetValue(p);

                    if (obj is not null)
                    {
                        data = obj.ToString();
                        value = data;
                    }
                }

                return value?.Trim();
            }
            catch
            {
                return null;
            }
        }

        private sealed record MemberList(List<Member>? Members);
        private sealed class Member { [JsonProperty("@odata.id")] public string? OdataId { get; set; } }
    }
}
