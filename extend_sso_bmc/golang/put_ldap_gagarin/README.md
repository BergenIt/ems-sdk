# Настройка LDAP-авторизации в BMC на примере производителя Gagarin

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

Для активации работы LDAP по защищенному соединению необходимо получить CA. Для этого необходимо зайти на ВМ sdk и выполнить команду `docker exec -it ems-traefik-1 sh -c "wget -O - --no-check-certificate https://acme:443/roots.pem 2> /dev/null" > roots.pem`. Данный сертификат необходимо загрузить в проект через волюм в docker-compose.yaml слудующим образом:

```yaml
volumes:
    - roots.pem:roots.pem
```

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

Далее выполним установку настроек в LDAP BMC по адресу `/redfish/v1/AccountService/` методом `PATCH` и выстроим тело запроса в зависимости от режима установки:

```golang
// Включение/выключение лдапа на бмс Gagarin
func putSettings(client *RedfishClient, state pb.SsoState, ssoDn, ssoPassword string) error {
    body := createLDAPSetBody(ssoDn, ssoPassword, state)
    if err := client.PatchData(ACCOUNT_SERVICE_PAGE, body); err != nil {
        return err
    }

    return nil
}

func createLDAPSetBody(ssoDn, ssoPassword string, state pb.SsoState) string {
    req := ""

    if state == pb.SsoState_SSO_STATE_ACTIVE {
        req = strings.Replace(enableSSORequestGagarin, "@SsoDn", ssoDn, 1)
        req = strings.Replace(req, "@SsoPassword", ssoPassword, 1)
        req = strings.Replace(req, "@SsoAddress", SSO_HOST+":389", 1)
        req = strings.Replace(req, "@SsoRootDn", BASE_DN, 1)
        req = strings.Replace(req, "@SsoGroupAttribute", GROUP_FORMAT, 1)
        req = strings.Replace(req, "@SsoUserNameAttribute", NAME_FORMAT, 1)
    } else {
        req = disableSSORequestGagarin
    }

    return req
}
```

Тело шаблона запроса активации LDAP к BMC имеет следующий вид (значения с символом '@' необходимо заменить):

```json
{
    "LDAP": {
        "ServiceEnabled": true,
        "Authentication": {
            "Username": "@SsoDn",
            "Password": "@SsoPassword",
            "AuthenticationType": "UsernameAndPassword"
        },
        "ServiceAddresses": [
            "ldap://@SsoAddress"
        ],
        "RemoteRoleMapping": [
            {
                "LocalRole": "Operator",
                "RemoteGroup": "accounts"
            }
        ],
        "LDAPService": {
            "SearchSettings": {
                "BaseDistinguishedNames": [
                "@SsoRootDn"
                ],
                "GroupsAttribute": "gidNumber",
                "UsernameAttribute": "cn"
            }
        }
    }
}
```

Тело шаблона запроса деактивации LDAP к BMC имеет следующий вид:

```json
{"LDAP":{"ServiceEnabled":false}}
```

В завершении отправляем полученный результат как ответ по RPC.

Пример готового проекта расположен в папке [project](./project)
