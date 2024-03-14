# Настройка LDAP-авторизации в BMC на примере производителя Huawei

Данная операция необходима для настройки LDAP-авторизации в BMC.

## Обзор операции

Данная операция реализуется следующим RPC:

```proto
service SsoCenter {
  rpc PutSettings(PutSsoSettingsRequest) returns (PutSsoSettingsResponse);
}

message PutSsoSettingsRequest {
  // Информация о оборудовании
  DeviceContent device = 1;
  // Статус для установки
  SsoState target_state = 2;
  // DN для LDAP
  string sso_dn = 3;
  // Пароль для бинда в LDAP
  string sso_password = 4;
}

message PutSsoSettingsResponse {
  // Результат выполнения операции
  OperationResult result = 1;
}
```

С полной структурой данных вы можете ознакомиться в [протофайлах](../../../.proto/service_sso_center.proto).

## Пример реализации

Дополняем уже имеющийся [шаблон](../create_project/project/main.go) релизацией RPC `PutSettings`.

Для начала нам необходимо спарсить данные для подключения к BMC по протоколу REDFISH из входящих данных:

```golang
func findCreds(in []*pb.DeviceConnector, protocol pb.ConnectorProtocol) (*pb.Credential, string, error) {
    for _, connector := range in {
        for _, creds := range connector.Credentials {
            if creds.Protocol == protocol {
                if creds.Login == "" || creds.Password == "" {
                    return nil, connector.Address, fmt.Errorf("login or password can not be empty")
                }

                return creds, "", nil
            }
        }
    }

    return nil, "", fmt.Errorf("creds not found")
}
```

Далее необходимо создать инстанс клиента к REDFISH для удобной работы с протколом:

```golang
// Создание редфиш клиента
redfishClient := newRedfishClient(creds.Login, creds.Password, address, creds.Port)
```

Далее необходимо получить сессию для взаимодействия с BMC по эндпоинту `/redfish/v1/SessionService/Sessions/` и обеспечить ее закрытие по окончанию операции:

``` golang
headers, err := redfishClient.GetAllHeadersFromPage(
        SessionsPageEndpoint,
        http.MethodPost,
        strings.NewReader(fmt.Sprintf(`{"UserName": "%s","Password": "%s"}`, redfishClient.Username,        redfishClient.Password)),
    )
if err != nil {
    return nil, fmt.Errorf("get headers error: %s", err)
}

location := headers.Get("Location")
if location == "" {
    return nil, fmt.Errorf("empty location header")
}

defer func() {
    // закрытие сессии
    if err := redfishClient.DeleteByPage(location); err != nil {
        log.Printf("failed close huawei session: %s", err)
    }
}()
```

Далее необходимо получить хедер `X-Auth-Token` из хедеров ответа по запросу сессии:

```golang
xauth := headers.Get("X-Auth-Token")
if xauth == "" {
    return nil, fmt.Errorf("empty xauth header")
}
```

Далее необходимо изменить базовые настройки LDAP, для этого нам понадобится хедер `ETag` страницы настройки LDAP и хедер `X-Auth-Token` сессии. Эндпоинт для страницы настроек LDAP `/redfish/v1/AccountService/LdapService/`:

```goalng
etag, err := client.GetHeaderFromPage(LDAPPageEndpoint, "ETag", http.MethodGet, nil)
if err != nil {
    return fmt.Errorf("get etag error: %s", err)
}

var headers map[string]string = map[string]string{
    "If-Match":     etag,
    "X-Auth-Token": xauth,
}

if err := client.PatchDataWithHeaders(LDAPPageEndpoint, createLDAPManageBody(state), headers); err != nil {
    return fmt.Errorf("run redfish request error: %s", err)
}
```

При активации LDAP авторизации тело запроса будет иметь следующий вид:

```json
{ "LdapServiceEnabled": true }
```

При деактивации LDAP авторизации тело запроса будет иметь следующий вид:

```json
{ "LdapServiceEnabled": false }
```

Так же при активации LDAP необходимо установить дополнительные параметры по эндпоинту `/redfish/v1/AccountService/LdapService/LdapControllers/1` с помощью метода `PATCH`:

```golang
// Установка настроек лдапа в бмс Huawei (работает только при включении лдапа)
func setLDAPSettings(client *RedfishClient, xauth string) error {
    etag, err := client.GetHeaderFromPage(LDAPSettingsPageEndpoint, "ETag", http.MethodGet, nil)
    if err != nil {
        return fmt.Errorf("get etag error: %s", err)
    }

    var headers map[string]string = map[string]string{
        "If-Match":     etag,
        "X-Auth-Token": xauth,
    }

    if err := client.PatchDataWithHeaders(LDAPSettingsPageEndpoint, createLDAPSettingsBody(), headers); err != nil {
        return fmt.Errorf("run redfish request error: %s", err)
    }

    return nil
}
```

Тело запроса для дополнительных настроек LDAP при активации имеет следующий вид:

```json
{
    "LdapServerAddress": "@LDAPAddress",
    "LdapPort": @LDAPPort,
    "UserDomain": ",DC=bergen.ems",
    "BindDN": "cn=nerd,dc=bergen,dc=ems",
    "BindPassword": "0penBmc",
    "CertificateVerificationEnabled ": false,
    "CertificateVerificationLevel": "Demand",
    "LdapGroups": [
        {
            "GroupName": "accounts",
            "GroupRole": "Operator",
            "GroupFolder": "OU=groups",
            "GroupLoginRule": [],
            "GroupLoginInterface": [
                "Web",
                "SSH",
                "Redfish"
            ]
        }
    ]
}
```

Так же при активации LDAP необходимо загрузить CA для защищенного взаимодействия между BMC и LDAP по эндпоинту `/redfish/v1/AccountService/LdapService/LdapControllers/1/Actions/HwLdapController.ImportCert` с помощью метода `POST`:

```golang
func loadLDAPCA(client *RedfishClient, xauth string) error {
    var headers map[string]string = map[string]string{
        "X-Auth-Token": xauth,
    }

    ca, err := loadCa()
    if err != nil {
        return fmt.Errorf("load CA: %s", err)
    }

    if err := client.PostDataWithHeaders(LDAPSetCAPageEndpoint, createLDAPSetCABody(string(ca)), headers); err != nil {
        return fmt.Errorf("run redfish request error: %s", err)
    }

    return nil
}

func loadCa() (string, error) {
    req, _ := http.NewRequest(http.MethodGet, "https://traefik:7071/roots.pem", nil)
    c := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
            },
        },
    }

    response, err := c.Do(req)
    if err != nil {
        return "", fmt.Errorf("http GET: %s", err)
    }
    defer response.Body.Close()

    ca, err := io.ReadAll(response.Body)
    if err != nil {
        return "", fmt.Errorf("read response body: %s", err)
    }

    return string(ca), nil
}
```

Тело запроса для установки CA имеет следующий вид:

```json
{"Type":"text","Content": "@Cert"}
```

В завершении отправляем полученный результат как ответ по RPC.

Пример готового проекта расположен в папке [project](./project)
