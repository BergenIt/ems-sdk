# Сбор температуры устройств

После формирования скелета проекта можно перейти к наполнению логикой операции сбора температуры устройств модуля расширения template.

В данном разделе описана логика сбора информации в рамках этой системной операции.

## Обзор операции

Данная операция отвечает за получение данных по операционной системе устройства по протоколу SNMP. Выполняется после получения сервисом следующего GRPC-запроса.

Запрос отправляемый сервису должен соответствовать `CollectTemplateTemperatureRequest`:

* `device`
  * Тип параметра: `DeviceContent`
  * Описание Данные по 1 устройству.

* `metric_templates`:
  * Тип параметра: `RepeatedField<SystemMetricTemplate>`
  * Описание: Идентификатор сетевого интерфейса.

Класс `SystemMetricTemplate`:

* `system_metric_template_id`:
  * Тип параметра: `string`
  * Описание: Идентификатор шаблона системной метрики.
* `template`:
  * Тип параметра: `string`
  * Описание: Шаблон для сбора данных системной метрики (путь к данным).
* `system_metric`:
  * Тип параметра: `SystemMetric (Enum)`
  * Описание: Тип системной метрики.

Данная структура является общей для реализации операции сбора температуры устройств по разным протоколам, поэтому может содержать большее количество полей, чем поддерживает template.

Для корректной работы сбора данных, устройство должно иметь хотя бы одно действительное подключение с протоколом SNMP (поле **`protocol`** списка **`credentials`** со значением **`CONNECTOR_PROTOCOL_SNMP`**), остальные устройства модуль должен игнорировать.

Из тела полученного запроса сервис берет данные необходимые для отправки SNMP-запроса к оборудованию. Общий вид используемого SNMP-запроса выглядит следующим образом:

```bash
snmpbulkwalk -OentU -v version -l security_name -u Login -a auth_protocol -A auth_key -x private_protocol -X private_key -c community -n context address oid
```

Составляющие запроса:

* `snmpbulkwalk` Команда для получения данных
* `-OentU` Параметры вывода информации
* `-v version` Версия протокола
* `-l security_name` Уровень аутентификации: noAuthNoPriv / authNoPriv / authPriv
* `-u Login` Логин безопасности (только для версии 3)
* `-a auth_protocol` Используемый протокол аутентификации (только для версии 3)
* `-A auth_key` Ключ для аутентификации (только для версии 3)
* `-x private_protocol` Используемый алгоритм шифрования (только для версии 3)
* `-X private_key` Ключ для шифрования (только для версии 3)
* `-c community` Community слово (пароль для версий 1 и 2)
* `-n context` Контекст подключения (только для версии 3)
* `address` IP адрес устройства
* `oid` Идентификатор (путь) к запрашиваемому атрибуту. oid для данной операции получается из SystemMetricTemplate.template

Полученные от устройства данные сервис отправляет GRPC-ответом.

Ответ сервиса должен соответствовать классу CollectTemplateTemperatureResponse:

* `temperature`:
  * Тип параметра: `DeviceTemperature`
  * Описание: Температура устройства.

Тип `DeviceTemperature`:

* `device_identity`:
  * Тип параметра: `DeviceDataIdentity`
  * Описание: Описание источника сбора данных.
* `temperature`:
  * Тип параметра: `google.protobuf.DoubleValue`
  * Описание: Идентификатор сетевого устройства, по которому было собрано значение.
Тип `DeviceDataIdentity`:

* `device_id`:
  * Тип параметра: `string`
  * Описание: Идентификатор устройства.
* `source`:
  * Тип параметра: `ServiceSource (Enum)`
  * Описание: Идентификатор rpc, с которого были собраны данные.

## Пример реализации

Реализация операции будет производиться на ???? со следующими характеристикам:

* ????

Для отправки Snmp запросов будем использовать библиотеку `Lextm.SharpSnmpLib` - [github](https://github.com/lextudio/sharpsnmplib?tab=readme-ov-file)

### Реализация отправки Snmp запросов

Создадим класс c данными необходимыми для подключения:

```csharp
public record SnmpCredential(
    string Ip
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
```

Реализуем отправку запросов:

```csharp
string SendRequest(SnmpCredential credential, string oidTemplate, int port, int timeout)
{
    string result;

    // Переменные используемые библиотекой Lextm.SharpSnmpLib
    IPEndPoint endpoint = new(IPAddress.Parse(credential.Ip), port);
    OctetString community = new(credential.Community);
    ObjectIdentifier oid = new(oidTemplate);
    VersionCode versionCode = (credential.Version == 1 || credential.Version == 2) ? VersionCode.V2 : VersionCode.V3;

    // Попытка получить данные по конкретному пути (реализация функции Get)
    // Например - При запросе 1.2.3.4 вернет 1.2.3.4
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

    //Возвращаем если получен валидный ответ
    if (resultGet != "NoSuchObject" && resultGet != "Null")
    {
        return resultGet;
    }

    // Если по конкретному пути получить данные не вышло пробуем получить следующий по списку (реализация функции Walk)
    // Например - При запросе 1.2.3.4 вернет 1.2.3.4 или 1.2.3.4.0
    
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

    //При возврате пустого массива возвращаем ошибку
    return resultGetBulk.Count == 0 ? "Error: No valid oid" : resultGetBulk.First().Data.ToString();
}
```

### Реализация RPC

Добавим в проект необходимые протофайлы:

1. service_template_manager.proto
2. shared_common.proto
3. shared_device.proto
4. shared_device_available
5. shared_device_initial.proto
6. shared_device_operation_system.proto
7. shared_device_power_usage.proto
8. shared_device_temperature.proto
9. shared_device_template.proto

Полный список протофайлов вы можете найти в [директории](../../../.proto).

Реализуем требуемое rpc:

```csharp
public class MyService : TemplateManager.TemplateManagerBase
{
    public override Task<CollectTemplateTemperatureResponse> CollectTemperature(CollectTemplateTemperatureRequest request, ServerCallContext context)
    {
        // Переменные необходимые для ответа
        string deviceId = request.Device.DeviceId;
        double? temperature = null;

        // Получение данных для подключений (их может быть несколько)
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

        // Получение путей к данным (их может быть несколько)
        IEnumerable<string> templates = request.MetricTemplates.Where(template => template.SystemMetric == SystemMetric.DeviceTemperature).Select(template => template.Template);

        foreach (string template in templates)
        {
            foreach(SnmpCredential cred in connectCreds)
            {
                string respond = SnmpClient.SendRequest(cred, template, 161, 10000);
                if (!respond.StartsWith("Error"))        
                    //Валидация данных и остановка опроса в случае успеха
                    if (double.TryParse(respond, out double temp))                            
                        if (temp < 2000 && temp > 0)
                        {
                            temperature = temp;
                            break;
                        }
            }

            if (temperature != null) break;
        }

        // Создание и отправка ответа
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
```

## Проверка реализации

Проверить работу сервиса можно через Postman, подробнее об этом описано здесь:

* <https://learning.postman.com/docs/sending-requests/grpc/first-grpc-request/>

Пример тела запроса для проверки:

```json
{
    "device": {
        "connectors": [
            {
                "credentials": [
                    {
                        "port": 161,
                        "security_level": "security_level",
                        "version": 2,
                        "security_name": "security_name",
                        "protocol": "CONNECTOR_PROTOCOL_SNMP",
                        "private_key": "private_key",
                        "context": "context",
                        "login": "login",
                        "community": "public",
                        "auth_protocol": "auth_protocol",
                        "auth_key": "auth_key",
                        "password": "password",
                        "cipher": 0,
                        "private_protocol": "private_protocol"
                    }
                ],
                "address": "10.1.18.23"
            }
        ]
    },
    "metric_templates": [
        {
            "system_metric": 14,
            "template": ".1.3.6.1.2.1.2.2.1.3.1"
        }
    ]
}
```

Пример тела ответа:

```json
"temperature": {
        "device_identity": {
            "device_id": "labore enim fugiat",
            "access_object_id": "",
            "source": "SERVICE_SOURCE_TEMPLATE_MANAGER"
        },
        "temperature": {
            "value": 6
        }
    }
```

Пример готового проекта расположен в папке [project](./project)
