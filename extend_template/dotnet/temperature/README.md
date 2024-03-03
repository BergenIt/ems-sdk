# Название сценария
Пример реализации системной операции "Получению температуры по протоколу SNMP". ссылка на документацию(https://docs.bergen.tech/ems/release-documents/latest/#/specifications/ds/host-domain/template-manager/README?id=%d0%9f%d0%be%d0%bb%d1%83%d1%87%d0%b8%d1%82%d1%8c-%d1%82%d0%b5%d0%bc%d0%bf%d0%b5%d1%80%d0%b0%d1%82%d1%83%d1%80%d1%83-%d1%83%d1%81%d1%82%d1%80%d0%be%d0%b9%d1%81%d1%82%d0%b2%d0%b0)
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

Варианты `SystemMetric`:

* SYSTEM_METRIC_UNSPECIFIED = 0;
* SYSTEM_METRIC_POWER_STATE = 1;
* SYSTEM_METRIC_BOOT_GET = 2;
* SYSTEM_METRIC_MODEL = 3;
* SYSTEM_METRIC_VENDOR = 4;
* SYSTEM_METRIC_SERVER_LED = 5;
* SYSTEM_METRIC_POWER_USAGE = 6;
* SYSTEM_METRIC_FIRMWARE_BOOT_SOURCE_GET = 7;
* SYSTEM_METRIC_SERIAL_NUMBER = 8;
* SYSTEM_METRIC_BIOS_VERSION = 9;
* SYSTEM_METRIC_BMC_VERSION = 10;
* SYSTEM_METRIC_OPERATION_SYSTEM = 11;
* SYSTEM_METRIC_HOSTINFO = 12;
* SYSTEM_METRIC_PROCESSOR = 13;
* SYSTEM_METRIC_DEVICE_TEMPERATURE = 14;

Класс `DeviceContent`:

* `device_id`:
    * Тип параметра: `string`
    * Описание: Идентификатор устройства.
* `model_name`:
    * Тип параметра: `string`
    * Описание: Название модели устройства.
* `vendor_name`:
    * Тип параметра: `string`
    * Описание: Название вендора устройства.
* `connectors`:
    * Тип параметра: `RepeatedField<DeviceConnector>`
    * Описание: Идентификатор сетевого интерфейса.
Класс `DeviceConnector`:

* `device_network_id`:
    * Тип параметра: `string`
    * Описание: Идентификатор сетевого интерфейса устройства.
* `address`:
    * Тип параметра: `string`
    * Описание: IP/FQDN адрес устройства.
* `mac`:
    * Тип параметра: `string`
    * Описание: MAC-адрес устройства.
* `credentials`:
    * Тип параметра: `RepeatedField<Credential>`
    * Описание: Учетные данные подключения.
Класс `Credential`:

* `protocol`:
    * Тип параметра: ConnectorProtocol (Enum)
    * Описание: Протокол подключения.
* `login`:
    * Тип параметра: `string`
    * Описание: Логин для подключения.
* `password`:
    * Тип параметра: `string`
    * Описание: Пароль для подключения.
* `port`:
    * Тип параметра: `int32`
    * Описание: Порт подключения.
* `cipher`:
    * Тип параметра: `int32`
    * Описание: Шифрование (только для IPMI).
* `version`:
    * Тип параметра: `uint32`
    * Описание: Версия протокола (только для SNMP).
* `community`:
    * Тип параметра: `string`
    * Описание: Community слово (только для SNMP).
* `security_name`:
    * Тип параметра: `string`
    * Описание: Security name (только для SNMP).
* `context`:
    * Тип параметра: `string`
    * Описание: Контекст подключения (только для SNMP).
* `auth_protocol`:
    * Тип параметра: `string`
    * Описание: Auth protocol (только для SNMP).
* `auth_key`:
    * Тип параметра: `string`
    * Описание: Auth key (только для SNMP).
* `private_protocol`:
    * Тип параметра: `string`
    * Описание: Private protocol (только для SNMP).
* `private_key`:
    * Тип параметра: `string`
    * Описание: Private key (только для SNMP).
* `security_level`:
    * Тип параметра: `string`
    * Описание: Уровень безопасности.

Варианты `ConnectorProtocol`:

* CONNECTOR_PROTOCOL_UNSPECIFIED - Невалидное значение.
* CONNECTOR_PROTOCOL_IPMI - Ipmi протокол для проверки подключения.
* CONNECTOR_PROTOCOL_REDFISH - Redfish протокол для проверки подключения.
* CONNECTOR_PROTOCOL_SNMP - Snmp протокол для проверки подключения.
* CONNECTOR_PROTOCOL_SSH - Ssh протокол для проверки подключения.
* CONNECTOR_PROTOCOL_WMI - Wmi протокол для проверки подключения.

Из тела полученного запроса сервис берет данные необходимые для отправки SNMP-запрос к оборудованию. Общий вид используемого SNMP-запроса выглядит следующим образом:
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
Варианты `ServiceSource`:

* SERVICE_SOURCE_UNSPECIFIED - Невалидное значение.
* SERVICE_SOURCE_BMC_MANAGER - Реализация управления и сбора с BMC.
* SERVICE_SOURCE_LINUX_MANAGER - Реализация управления и сбора с linux-хостов.
* SERVICE_SOURCE_WINDOWS_MANAGER - Реализация управления и сбора с windows-хостов.
* SERVICE_SOURCE_HYPERVISOR_MANAGER - Реализация управления и сбора с гипервизоров.
* SERVICE_SOURCE_TEMPLATE_MANAGER - Реализация шаблонного мониторинга.

## Пример реализации


> Описываем на примере чего будем реализовывать (желательно брать что-то не похожее на наше стандартное окружение)

> Описываем откуда мы знаем как это реализовывать (например ссылка на вендорское описание апи)

> Вкладываем снипетами куски кода прям сюда и описываем что это и зачем это и почему так

> Указываем что пример итогового проекта находится в папке project
