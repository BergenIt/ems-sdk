using ToolCluster.V4;
using Grpc.Core;
using SnmpExample.Snmp;

namespace SnmpExample.Services
{
    public class MyService : TemplateManager.TemplateManagerBase
    {
        public override Task<CollectTemplateTemperatureResponse> CollectTemperature(CollectTemplateTemperatureRequest request, ServerCallContext context)
        {
            {
                string deviceId = request.Device.DeviceId;
                double? temperature = null;

                IEnumerable<SnmpCredential> connectCreds = request.Device.Connectors
                    .SelectMany(connector => connector.Credentials.Where(credentials => credentials.Protocol == ConnectorProtocol.Snmp)
                    .Select(credentials => new SnmpCredential(
                        connector.Address,
                        credentials.Login,
                        credentials.Password,
                        credentials.Port,
                        credentials.Version,
                        credentials.SecurityName,
                        credentials.SecurityLevel,
                        credentials.Community,
                        credentials.Context,
                        credentials.AuthProtocol,
                        credentials.AuthKey,
                        credentials.PrivateProtocol,
                        credentials.PrivateKey)
                ));

                IEnumerable<string> templates = request.MetricTemplates.Where(template => template.SystemMetric == SystemMetric.DeviceTemperature).Select(template => template.Template);

                foreach (string template in templates)
                {
                    foreach(SnmpCredential cred in connectCreds)
                    {
                        string respond = SnmpClient.SendRequest(cred, template, 161, 10000);
                        if (!respond.StartsWith("Error"))                        
                            if (double.TryParse(respond, out double temp))                            
                                if (temp < 2000 && temp > 0)
                                {
                                    temperature = temp;
                                    break;
                                }
                    }
                    if (temperature != null) break;
                }

                return Task.FromResult(new CollectTemplateTemperatureResponse()
                {
                    Temperature = new DeviceTemperature()
                    {
                        DeviceIdentity = new DeviceDataIdentity()
                        {
                            DeviceId = deviceId,
                            Source = ServiceSource.TemplateManager
                        },
                        Temperature = temperature 
                    }
                });
            }
        }
    }
}
