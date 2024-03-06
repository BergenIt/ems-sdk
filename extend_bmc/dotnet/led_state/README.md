# Cбор статуса LED

После формирования скелета проекта можно перейти к наполнению логикой операции сбора статуса LED модуля расширения Bmc.

В данном разделе описана логика сбора информации в рамках этой системной операции.

## Обзор операции

Данная операция отвечает за сбор данных по текущему статуса LED устройства - включен, выключен или моргает с использованием протокола `Redfish`.

Для реализации операции RPC будет иметь следующую сигнатуру:

* **`rpc CollectLedState(CollectBmcLedStateRequest) returns (CollectBmcLedStateResponse)`**

На вход модуль получает следующий запрос в формате Protocol buffers:

Тип `CollectBmcLedStateRequest`:

* `device`:
  * **Тип параметра:** `DeviceContent`
  * **Описание:** Данные об устройстве.
* `metric_templates`:
  * **Тип параметра:** `repeated SystemMetricTemplate`
  * **Описание:** Список шаблонов по сбору метрик.

Тип `DeviceContent`:

* `device_id`:
  * **Тип параметра:** `string`
  * **Описание:** Идентификатор устройства.
* `model_name`:
  * **Тип параметра:** `string`
  * **Описание:** Модель устройства.
* `vendor_name`:
  * **Тип параметра:** `string`
  * **Описание:** Вендор устройства.
* `connectors`:
  * **Тип параметра:** `repeated DeviceConnector`
  * **Описание:** Список интерфейсов подключения к устройству.

Тип `DeviceConnector`:

* `device_network_id`:
  * **Тип параметра:** `string`
  * **Описание:** Идентификатор сетевого интерфейса устройства.
* `address`:
  * **Тип параметра:** `string`
  * **Описание:** Адрес подключения (ip/fqdn).
* `mac`:
  * **Тип параметра:** `string`
  * **Описание:** MAC-адрес устройства.
* `credentials`:
  * **Тип параметра:** `repeated Credential`
  * **Описание:** Список данных подключения к устройству.

Тип `Credential`:

* `protocol`:
  * **Тип параметра:** `ConnectorProtocol`
  * **Описание:** Протокол подключения.
* `login`:
  * **Тип параметра:** `string`
  * **Описание:** Логин для подключения.
* `password`:
  * **Тип параметра:** `string`
  * **Описание:** Пароль для подключения.
* `port`:
  * **Тип параметра:** `int32`
  * **Описание:** Порт подключения.
* `cipher`:
  * **Тип параметра:** `int32`
  * **Описание:** Шифрование (только для IPMI).
* `version`:
  * **Тип параметра:** `int32`
  * **Описание:** Версия протокола (только для SNMP).
* `community`:
  * **Тип параметра:** `string`
  * **Описание:** Community слово (только для SNMP).
* `security_name`:
  * **Тип параметра:** `string`
  * **Описание:** Security name (только для SNMP).
* `context`:
  * **Тип параметра:** `string`
  * **Описание:** Контекст подключения (только для SNMP).
* `auth_protocol`:
  * **Тип параметра:** `string`
  * **Описание:** Auth protocol (только для SNMP).
* `auth_key`:
  * **Тип параметра:** `string`
  * **Описание:** Auth key (только для SNMP).
* `private_protocol`:
  * **Тип параметра:** `string`
  * **Описание:** Private protocol (только для SNMP).
* `private_key`:
  * **Тип параметра:** `string`
  * **Описание:** Private key (только для SNMP).
* `security_level`:
  * **Тип параметра:** `string`
  * **Описание:** Уровень безопасности.

Тип `SystemMetricTemplate`:

* `system_metric_template_id`:
  * **Тип параметра:** `string`
  * **Описание:** Идентификатор шаблона системной метрики.
* `template`:
  * **Тип параметра:** `string`
  * **Описание:** Шаблон для сбора данных системной метрики (путь к данным).
* `system_metric`:
  * **Тип параметра:** `SystemMetric`
  * **Описание:** Тип системной метрики.

Перечисление `SystemMetric`:

* `SYSTEM_METRIC_POWER_UNSPECIFIED`:
  * **Описание:** Инвалид.
* `SYSTEM_METRIC_POWER_STATE`:
  * **Описание:** Метрика текущего состояния питания.
* `SYSTEM_METRIC_BOOT_GET`:
  * **Описание:** Метрика получения текущего загрузочного носителя.
* `SYSTEM_METRIC_MODEL`:
  * **Описание:** Метрика получения модели устройства.
* `SYSTEM_METRIC_VENDOR`:
  * **Описание:** Метрика получения вендора.
* `SYSTEM_METRIC_SERVER_LED`:
  * **Описание:** Метрика получения текущего значения LED сервера.
* `SYSTEM_METRIC_POWER_USAGE`:
  * **Описание:** Метрика получения текущего энергопотребления.
* `SYSTEM_METRIC_FIRMWARE_BOOT_SOURCE_GET`:
  * **Описание:** Метрика получения текущего режима загрузки BIOS/UEFI.
* `SYSTEM_METRIC_SERIAL_NUMBER`:
  * **Описание:** Метрика получения серийного номера устройства.
* `SYSTEM_METRIC_BIOS_VERSION`:
  * **Описание:** Метрика получения текущей версии прошивки BIOS/UEFI.
* `SYSTEM_METRIC_BMC_VERSION`:
  * **Описание:** Метрика получения текущей версии прошивки BMC.
* `SYSTEM_METRIC_OPERATION_SYSTEM`:
  * **Описание:** Метрика получения текущей операционной системы.
* `SYSTEM_METRIC_HOSTINFO`:
  * **Описание:** Метрика получения текущей информации о хостах оборудования.
* `SYSTEM_METRIC_PROCESSOR`:
  * **Описание:** Метрика получения текущей информации о процессорах устройства.
* `SYSTEM_METRIC_DEVICE_TEMPERATURE`:
  * **Описание:** Метрика получения текущей информации о температуре устройства.

Данная структура запроса является общей для реализации операции сбора статуса LED по разным протоколам, поэтому может содержать большее количество полей, чем поддерживает Bmc.

Для корректной работы сбора данных, устройство должно иметь хотя бы одно действительное подключение с протоколом Redfish (поле **`protocol`** списка **`credentials`** со значением **`CONNECTOR_PROTOCOL_REDFISH`**), остальные устройства модуль должен игнорировать.

Пример данных запроса для сбора статуса LED:

```protobuf
{
  "device": {
    "device_id": "test",
    "connectors": [
      {
        "address": "192.168.1.55",
        "credentials": [
          {
            "login": "root",
            "password": "1qaz@WSX",
            "port": 443,
            "protocol": "CONNECTOR_PROTOCOL_REDFISH"
          }
        ]
      }
    ]
  }
}
```

В качестве cообщения-ответа в RPC используется модель:

Тип `CollectBmcLedStateResponse`:

* `led`:
  * **Тип параметра:** `DeviceLed`
  * **Описание:** Текущий статус Led устройства.

Тип `DeviceLed`:

* `device_identity`:
  * **Тип параметра:** `DeviceDataIdentity`
  * **Описание:** Описание источника сбора данных.
* `state`:
  * **Тип параметра:** `LedState`
  * **Описание:** Текущий статус Led.

Тип `DeviceDataIdentity`:

* `device_id`:
  * **Тип параметра:** `string`
  * **Описание:** Идентификатор устройства.
* `source`:
  * **Тип параметра:** `ServiceSource`
  * **Описание:** Идентификатор rpc, с которого были собраны данные.

Перечисление `ServiceSource`:

* `SERVICE_SOURCE_UNSPECIFIED`:
  * **Описание:** Невалидное значение.
* `SERVICE_SOURCE_BMC_MANAGER`:
  * **Описание:** Реализация управления и сбора с BMC.
* `SERVICE_SOURCE_LINUX_MANAGER`:
  * **Описание:** Реализация управления и сбора с linux-хостов.
* `SERVICE_SOURCE_WINDOWS_MANAGER`:
  * **Описание:** Реализация управления и сбора с windows-хостов.
* `SERVICE_SOURCE_HYPERVISOR_MANAGER`:
  * **Описание:** Реализация управления и сбора с гипервизоров.
* `SERVICE_SOURCE_TEMPLATE_MANAGER`:
  * **Описание:** Реализация шаблонного мониторинга.

Перечисление `LedState`:

* `LED_STATE_UNSPECIFIED`:
  * **Описание:** Невалидное значение.
* `LED_STATE_UNKNOWN`:
  * **Описание:** Неизвестно.
* `LED_STATE_OFF`:
  * **Описание:** Led выключен.
* `LED_STATE_ON`:
  * **Описание:** Led включен.
* `LED_STATE_BLINK`:
  * **Описание:** Led моргает.

С полной структурой данных вы можете ознакомиться в [протофайлах](../../../.proto).

## Пример реализации

Реализация операции будет производиться на устройстве вендора `Huawei` со следующими характеристикам:

**TODO: дополнить информацией**

* Название ОС: **`N/A`**
* Версия: **`N/A`**
* Процессор: **`N/A`**
* Объем памяти дисков: **`N/A`**
* Суммарно установлено памяти (ОЗУ): **`N/A`**
* Количество плашек ОЗУ: **`N/A`**

Для реализации сбора данных необходимо иметь возможность подключения к удаленному хосту по протоколу Redish, для этого необходимо сформировать набор данных подключения.

Изучив [документацию](https://support.huawei.com/enterprise/en/doc/EDOC1100105856/745aefb6/managing-chassis-resources) с официального сайта вендора **`Huawei`** становится понятно каким образом хранятся данные по **`Redfish`**.

URI домашней страницы службы Redfish (другое название корневой адрес службы) можно получить с помощью следующего URI: `https://<IP-адрес:порт>/redfish/v1`

Для получения данных статусу Led будет использоваться содержимое страницы `Chassis`, а итоговый адрес получения данных буде выглядит так - `https://<IP-адрес:порт>/redfish/v1/Chassis`

С этой содержимого этой страницы будет интересовать поля **`IndicatorLED`**
Подробнее по ссылке:

* <https://support.huawei.com/enterprise/en/doc/EDOC1100105856/745aefb6/managing-chassis-resources>

Процесс сбора статуса Led, который состоит из получения содержимого страницы:

```csharp
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
```

Процесс отправки запроса получения содержимого страницы:

```csharp
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
```

Процесс получения поля из содержимого страницы:

```csharp
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
```

Пример готового проекта расположен в папке [project](./project)
