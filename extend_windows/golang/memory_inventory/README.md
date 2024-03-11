# Сбор инвентарных данных по ОЗУ

После формирования скелета проекта можно перейти к наполнению логикой операции сбора инвентарных данных по ОЗУ модуля расширения Windows.

В данном разделе описана логика сбора информации в рамках этой системной операции.

## Обзор операции

Данная операция отвечает за сбор информации по инвентарным данным ОЗУ (модель, вендор, номер партии, слот устройства). Подключение к удаленной машине происходит по протоколу WMI. Как результат на этом хосте выполняется команда для получения данных.

Для реализации операции RPC будет иметь следующую сигнатуру:

* **`rpc CollectMemory(CollectWindowsMemoryRequest) returns (CollectWindowsMemoryResponse)`**

На вход модуль получает следующий запрос в формате Protocol buffers:

Тип `CollectWindowsMemoryRequest`:

* `device`:
  * **Тип параметра:** `DeviceContent`
  * **Описание:** Данные об устройстве.

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

Данная структура является общей для реализации операции сбора инвентарных данных по ОЗУ по разным протоколам, поэтому может содержать большее количество полей, чем поддерживает windows.

Для корректной работы сбора данных, устройство должно иметь хотя бы одно действительное подключение с протоколом WMI (поле **`protocol`** списка **`credentials`** со значением **`CONNECTOR_PROTOCOL_WMI`**), остальные устройства модуль должен игнорировать.

Пример данных запроса для получения данных по ОЗУ:

```protobuf
{
  "device": {
    "device_id": "test",
    "connectors": [
      {
        "address": "10.1.18.34",
        "credentials": [
          {
            "login": "Administrator",
            "password": "1qaz@WSX",
            "port": 5985,
            "protocol": "CONNECTOR_PROTOCOL_WMI"
          }
        ]
      }
    ]
  }
}
```

В качестве cообщения-ответа используется модель:

Тип `CollectWindowsMemoryResponse`:

* `memory`:
  * **Тип параметра:** `DeviceMemory`
  * **Описание:** Инвентарные данные по ОЗУ.

Тип `DeviceMemory`:

* `device_identity`:
  * **Тип параметра:** `DeviceDataIdentity`
  * **Описание:** Описание источника сбора данных.
* `memories`:
  * **Тип параметра:** `repeated MemoryCard`
  * **Описание:**  Плашки оперативной памяти устройства.

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

Тип `MemoryCard`:

* `memory_type`:
  * **Тип параметра:** `MemoryType`
  * **Описание:** Тип памяти.
* `memory_device_type`:
  * **Тип параметра:** `MemoryDeviceType`
  * **Описание:** Тип устройства памяти.
* `base_module_type`:
  * **Тип параметра:** `BaseModuleType`
  * **Описание:** Тип базового модуля памяти.
* `vendor`:
  * **Тип параметра:** `string`
  * **Описание:** Производитель или идентификатор поставщика (VendorID).
* `size`:
  * **Тип параметра:** `int32`
  * **Описание:** Емкость памяти в мегабайтах (CapacityMiB).
* `part_number`:
  * **Тип параметра:** `string`
  * **Описание:** Номер детали (PartNumber).
* `serial_number`:
  * **Тип параметра:** `string`
  * **Описание:** Серийный номер.
* `firmware_revision`:
  * **Тип параметра:** `string`
  * **Описание:** Ревизия прошивки.
* `slot`:
  * **Тип параметра:** `int32`
  * **Описание:** Слот расположения памяти (MemoryLocation.Slot).
* `state`:
  * **Тип параметра:** `MemoryState`
  * **Описание:** Состояние памяти.
* `socket`:
  * **Тип параметра:** `int32`
  * **Описание:** Разъем, к которому подключена память (MemoryLocation.Socket).
* `speed_mhz`:
  * **Тип параметра:** `int32`
  * **Описание:** Скорость работы в мегагерцах (OperatingSpeedMhz).
* `location`:
  * **Тип параметра:** `string`
  * **Описание:** Расположение устройства (DeviceLocator).

С полной структурой данных вы можете ознакомиться в [протофайлах](../../../.proto).

Так как сообщения для операций идут по Protocol Buffer, необходимо иметь возможность скомпилировать protobuf-файлы под конкретный язык программирования.

Подробнее об этом:

* <https://github.com/protocolbuffers/protobuf?tab=readme-ov-file>

## Пример реализации

Реализация операции будет производиться на виртуальной машине, запущенной c помощью программы по виртуализации **Proxmox** со следующими характеристикам:

* Название ОС: **`Microsoft Windows Server 2019 Standard`**
* Версия: **`10.0.17763 Build 17763`**
* Процессор: **`Common KVM processor, 2095 Mhz, 2 Core(s), 2 Logical Processor(s)`**
* Объем памяти дисков: **`100,00 GB`**
* Суммарно установлено памяти (ОЗУ): **`4,00 GB`**
* Количество плашек ОЗУ: **`1`**

Для реализации сбора данных необходимо иметь возможность подключения к удаленному хосту по протоколу WinRM, для этого необходимо настроить подключение на удаленной машине.

Подробнее по ссылке:

* <https://learn.microsoft.com/ru-ru/windows/win32/winrm/installation-and-configuration-for-windows-remote-management>

Имея возможность подключения к удаленной машине с ОС Windows мы можем отправлять консольные команды на запуск определенных утилит. В основом это утилита **`WMIC`**.

Подробнее по ссылкам:

* <https://learn.microsoft.com/ru-ru/windows/win32/wmisdk/wmi-start-page>
* <https://learn.microsoft.com/ru-ru/windows/win32/wmisdk/wmic>

Изучив [документацию](https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/computer-system-hardware-classes) с официального сайта майкрософт становится понятно, что для того, чтобы получить информацию с помощью утилиты WMIC нужно понимать, из какого системного класса мы хотим получить данные.

Также для оптимизации сбора данных было решено получать не весь список полей, а только те, который нам необходимы.

По итогу можно сформировать следующий шаблон команды для выполнения команд - `WMIC PATH @Class GET @Fields /format:list`, где:

* **@Class** - Имя класса WMI
* **@Fields** - Список атрибутов, которые необходимо получить

Используемые WMI классы и атрибуты для сбора инвентарных данных по ОЗУ:

* **[Win32_PhysicalMemory](https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-physicalmemory)** - **`BankLabel, TypeDetail, MemoryType, SMBIOSMemoryType, FormFactor, Manufacturer, Capacity, PartNumber, SerialNumber, Status, Speed, DeviceLocator`**

Пример вывода команды **`WMIC PATH Win32_PhysicalMemory get /format:list`** без фильтра по полям:

```out
    Attributes=0
    BankLabel=
    Capacity=4294967296
    Caption=Physical Memory
    ConfiguredClockSpeed=0
    ConfiguredVoltage=0
    CreationClassName=Win32_PhysicalMemory
    DataWidth=
    Description=Physical Memory
    DeviceLocator=DIMM 0
    FormFactor=8
    HotSwappable=
    InstallDate=
    InterleaveDataDepth=
    InterleavePosition=
    Manufacturer=QEMU
    MaxVoltage=0
    MemoryType=9
    MinVoltage=0
    Model=
    Name=Physical Memory
    OtherIdentifyingInfo=
    PartNumber=
    PositionInRow=
    PoweredOn=
    Removable=
    Replaceable=
    SerialNumber=
    SKU=
    SMBIOSMemoryType=7
    Speed=
    Status=
    Tag=Physical Memory 0
    TotalWidth=
    TypeDetail=2
    Version=
```

Запрос к удаленному устройству по WMI выглядит следующим образом:

```go
func SendWinRMCommand(
 ctx context.Context,
 ip, login, pass string,
 port int,
 cmd string,
) (string, error) {
 client, err := newWinRMClient(ip, login, pass, port)
 if err != nil {
  return "", fmt.Errorf("create WMI client error: %s", err)
 }

 var stdout, stderr bytes.Buffer
 _, err = client.RunWithContext(ctx, "chcp 866 | "+cmd, &stdout, &stderr)
 if err != nil {
  return "", fmt.Errorf("cmd [%s] WinRM error: %s", cmd, err)
 }

 stderrStr := stderr.String()
 stdoutStr := stdout.String()
 if stderrStr != "" {
  return "", fmt.Errorf("cmd [%s] error: %s", cmd, stderrStr)
 }

 reader := transform.NewReader(bytes.NewReader([]byte(stdoutStr)), charmap.CodePage866.NewDecoder())
 d, err := io.ReadAll(reader)
 if err != nil {
  return "", fmt.Errorf("encoding stdout [%s] error [%s]", stdoutStr, err)
 }

 return string(d), nil
}
```

Реализуем функцию для преобразования полученных данных в требуемый EMS формат:

```go
func parseRAMInvInfo(stdout string) []*pb.MemoryCard {
  var ramInvInfos []*pb.MemoryCard

  metricInfos := handleMetricStdout(stdout)
  for _, metricInfo := range metricInfos {
    ramInvInf := &pb.MemoryCard{}
    var (
      memoryDeviceType uint16
      smBIOSMemoryType uint32
    )

    for key, value := range metricInfo {
      var i64 int64

      switch key {
        case "BankLabel":
          i64, _ = parseStrToInt64(key, strings.ReplaceAll(value, "BANK ", ""))
          ramInvInf.Slot = int32(i64)
        case "Manufacturer":
          ramInvInf.Vendor = value
        case "PartNumber":
          ramInvInf.PartNumber = value

        // полный switch case можно увидеть в папке project с проектом
      
        default:
          continue
      }
    }

    ramInvInf.MemoryDeviceType = memDeviceTypeInternal(memoryDeviceType, smBIOSMemoryType)
    if ramInvInf != (&pb.MemoryCard{}) {
      ramInvInfos = append(ramInvInfos, ramInvInf)
    }
  }

  return ramInvInfos
}
```

Пример готового проекта расположен в папке [project](./project)
