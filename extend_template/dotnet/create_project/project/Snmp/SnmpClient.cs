using System.Net;
using Lextm.SharpSnmpLib;
using Lextm.SharpSnmpLib.Messaging;

namespace SnmpExample.Snmp
{
    public record SnmpCredential(
       string Ip,
       string Login,
       int Port,
       int Version,
       string SecurityName,
       string SecurityLevel,
       string? Community,
       string? Context,
       string? AuthProtocol,
       string? AuthKey,
       string? PrivateProtocol,
       string? PrivateKey);

    public class SnmpClient
    {
        public static string SendRequest(SnmpCredential credential, string oidTemplate, int port, int timeout)
        {
            string result;
            IPEndPoint endpoint = new(IPAddress.Parse(credential.Ip), port);
            OctetString community = new(credential.Community);
            ObjectIdentifier oid = new(oidTemplate);
            VersionCode versionCode = (credential.Version == 1 || credential.Version == 2) ? VersionCode.V2 : VersionCode.V3;

            string resultGet = "Null";
            try
            {
                GetRequestMessage message = new(0, versionCode, community, new List<Variable> { new(oid) });

                ISnmpMessage response = message.GetResponse(timeout, endpoint);
                if (response.Pdu().ErrorStatus.ToInt32() == 0)
                {
                    resultGet = response.Pdu().Variables.FirstOrDefault().Data.ToString();
                }
            }
            catch (Exception ex)
            {
                result = "Error: " + ex.Message;
            }

            if (resultGet != "NoSuchObject" && resultGet != "Null")
            {
                return resultGet;
            }

            List<Variable> resultGetBulk = new();
            try
            {
                GetBulkRequestMessage message = new(0, versionCode, community, 0, 1, new List<Variable> { new(oid) });

                ISnmpMessage response = message.GetResponse(timeout, endpoint);
                if (response.Pdu().ErrorStatus.ToInt32() == 0)
                {
                    resultGetBulk = response.Pdu().Variables.ToList();
                }
            }
            catch (Exception ex)
            {
                return "Error: " + ex.Message;
            }

            return resultGetBulk.Count == 0 ? "Error: No valid oid" : resultGetBulk.First().Data.ToString();
        }
    }
}
