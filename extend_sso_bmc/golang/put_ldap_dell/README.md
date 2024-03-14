# Настройка LDAP-авторизации в BMC на примере производителя DELL

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

Далее выполним установку настроек в LDAP BMC по адресу `/redfish/v1/AccountService/` методом `PATCH` и выстроим тело запроса в зависимости от режима установки:

```golang
// Включение/выключение лдапа на бмс Dell
func putSettings(client *RedfishClient, state pb.SsoState) error {
    body := createLDAPManageBody(state)

    if err := client.PatchData(ACCOUNT_SERVICE_PAGE, body); err != nil {
        return err
    }

    return nil
}

func createLDAPManageBody(state pb.SsoState) string {
    req := ""

    if state == pb.SsoState_SSO_STATE_ACTIVE {
        req = strings.Replace(enableLDAP, "@BaseDN", BASE_DN, 1)
        req = strings.Replace(req, "@SsoHost", SSO_HOST, 1)
    } else {
        req = disableLDAP
    }

    return req
}
```

Тело шаблона запроса активации LDAP к BMC имеет следующий вид (значения с символом '@' необходимо заменить):

```json
{
    "LDAP": {
        "LDAPService": {
            "SearchSettings": {
                "BaseDistinguishedNames": [
                    "@BaseDN"
                ],
                "GroupNameAttribute": "cn"
            }
        },
        "RemoteRoleMapping": [
            {
                "RemoteGroup": "cn=accounts,ou=groups,dc=bergen,dc=ems",
                "LocalRole": "Operator"
            }
        ],
        "ServiceAddresses": [
            "@SsoHost"
        ],
        "ServiceEnabled": true
    }
}
```

Тело шаблона запроса деактивации LDAP к BMC имеет следующий вид:

```json
{"LDAP":{"ServiceEnabled":false}}
```

Далее, при активации LDAP необходимо установить дополнительные параметры методом `PATCH`, адрес запроса необходимо сформировать следующим образом:

```golang
// Установка настроек лдапа в бмс Dell (работает только при включении лдапа)
func setLDAPAttrsDell(client *RedfishClient, ssoDn, ssoPassword string) error {
    b, err := client.GetPage(MANAGERS_PAGE)
    if err != nil {
        return fmt.Errorf("get managers page error: %s", err)
    }

    managers := &ManagersDell{}
    if err := json.Unmarshal(b, managers); err != nil {
        return fmt.Errorf("unmarshal managers page error: %s", err)
    }

    if len(managers.Members) <= 0 {
        return fmt.Errorf("managers not found")
    }

    body := createLDAPAttrsBodyDell(ssoDn, ssoPassword)
    attributesEndpoint := managers.Members[0].OdataID + "/Attributes"

    if err := client.PatchData(attributesEndpoint, body); err != nil {
        return err
    }

    return nil
}
```

Тело запроса установки дополнительных атрибутов выглядит следующим образом (значения с символом '@' необходимо заменить):

```json
{
    "Attributes": {
        "LDAP.1.Port": @SsoPort,
        "LDAP.1.CertValidationEnable": "Disabled",
        "LDAP.1.BindDN": "@BindDN",
        "LDAP.1.SearchFilter": "objectclass=*",
        "LDAP.1.BindPassword": "@BindPassword"
    }
}
```

Далее необходимо загрузить сертификат в LDAP для авторизации с защищенным соединением методом `POST` по адресу `/redfish/v1/Managers/iDRAC.Embedded.1/Oem/Dell/DelliDRACCardService/Actions/DelliDRACCardService.ImportSSLCertificate`. Файл должен быть загружен в корень рабочей директории проекта через волюм в docker-compose.yaml.

```golang
func loadLDAPCA(client *RedfishClient) error {
    ca, err := loadCA()
    if err != nil {
        return fmt.Errorf("load CA: %s", err)
    }

    if err := client.PostData(SET_CA_PAGE, createLoadCABody(ca)); err != nil {
        return fmt.Errorf("post CA error: %s", err)
    }

    return nil
}

func loadCA() (string, error) {
    ca, err := os.ReadFile(CA_LOCAL_PATH)
    if err != nil {
        return "", fmt.Errorf("read body error: %s", err)
    }

    return string(ca), nil
}
```

В завершении отправляем полученный результат как ответ по RPC.

Пример готового проекта расположен в папке [project](./project)
