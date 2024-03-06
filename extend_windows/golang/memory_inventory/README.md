# Сбор инвентарных данных по ОЗУ

## Обзор операции

Данная операция отвечает за сбор информации по инвентарным данным ОЗУ (модель, вендор, номер партии, слот устройства). Подключение к удаленной машине происходит по протоколу WMI. Как результат на этом хосте выполняется команда для получения данных.

Для реализации операции RPC будет иметь следующую сигнатуру:

 *  **`rpc CollectMemory(CollectWindowsMemoryRequest) returns (CollectWindowsMemoryResponse)`**

На вход должен получить следующий запрос в формате Protocol buffers:

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

Перечисление `ConnectorProtocol`:

  * `CONNECTOR_PROTOCOL_UNSPECIFIED`:
    * **Описание:** Невалидное значение.
  * `CONNECTOR_PROTOCOL_IPMI`:
    * **Описание:** Ipmi протокол для проверки подключения.
  * `CONNECTOR_PROTOCOL_REDFISH`:
    * **Описание:** Redfish протокол для проверки подключения.
  * `CONNECTOR_PROTOCOL_SNMP`:
    * **Описание:** Snmp протокол для проверки подключения.
  * `CONNECTOR_PROTOCOL_SSH`:
    * **Описание:** Ssh протокол для проверки подключения.
  * `CONNECTOR_PROTOCOL_WMI`:
    * **Описание:** Wmi протокол для проверки подключения.

Данная структура является общей для реализации операции сбора инвентарных данных по ОЗУ по разным протоколам, поэтому может содержать большее количество полей (в зависимости от протокола).

Для корректной работы сбора данных, устройство должно иметь хотя бы одно действительное подключение с протоколом WMI (поле **`protocol`** списка **`credentials`** со значением **`CONNECTOR_PROTOCOL_WMI`**).

Пример заполнения структуры запроса для получения данных по ОЗУ в формате Json по WMI:

```json
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

В качестве cообщения-ответа предполагается модель:

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

Тип `MemoryType`:

  * `MEMORY_TYPE_UNSPECIFIED`:
    * **Описание:** Невалидное значение.
  * `MEMORY_TYPE_DRAM`:
    * **Описание:** Dynamic Random Access Memory (DRAM) представляет собой тип оперативной памяти, который используется для временного хранения данных, к которым процессор имеет быстрый доступ.
  * `MEMORY_TYPE_NVDIMM_N`:
    * **Описание:** Non-Volatile Dual In-line Memory Module (NVDIMM_N) - это модуль памяти, который объединяет в себе характеристики оперативной и постоянной памяти.
  * `MEMORY_TYPE_NVDIMM_F`:
    * **Описание:** Non-Volatile Dual In-line Memory Module (NVDIMM_F) - второй тип NVDIMM, предоставляющий функциональность хранения данных при выключенном устройстве.
  * `MEMORY_TYPE_NVDIMM_P`:
    * **Описание:** Non-Volatile Dual In-line Memory Module (NVDIMM_P) - третий тип NVDIMM, который сочетает в себе характеристики постоянной памяти и дополнительной энергонезависимой памяти.
  * `MEMORY_TYPE_INTEL_OPTANE`:
    * **Описание:** Intel Optane DC Persistent Memory - технология памяти, разработанная Intel, которая сочетает в себе характеристики оперативной и постоянной памяти.
  * `MEMORY_TYPE_FPRAM`:
    * **Описание:** Fast-paged RAM - тип оперативной памяти с быстрым доступом.
  * `MEMORY_TYPE_SRAM`:
    * **Описание:** Static Random Access Memory (SRAM) - статическая оперативная память, которая сохраняет свое состояние до отключения питания.
  * `MEMORY_TYPE_S_DRAM`:
    * **Описание:** Synchronous DRAM (S-DRAM) - вид DRAM, работающий с синхронизацией с системным тактовым сигналом.
  * `MEMORY_TYPE_PSRAM`:
    * **Описание:** Pseudo-static RAM (PSRAM) - тип оперативной памяти, комбинирующий свойства SRAM и DRAM.
  * `MEMORY_TYPE_RAMBUS`:
    * **Описание:** Rambus DRAM (RDRAM) - высокоскоростной тип DRAM, разработанный Rambus Inc.
  * `MEMORY_TYPE_CMOS`:
    * **Описание:** Complementary Metal-Oxide-Semiconductor (CMOS) - технология производства полупроводников, используемая для создания интегральных микросхем.
  * `MEMORY_TYPE_EDO_RAM`:
    * **Описание:** Extended Data Output RAM (EDO RAM) - улучшенный тип DRAM, предоставляющий более быстрый доступ к данным.
  * `MEMORY_TYPE_WIN_DRAM`:
    * **Описание:** Window DRAM - DRAM с использованием окон для адресации памяти.
  * `MEMORY_TYPE_CACHE_DRAM`:
    * **Описание:** Cache DRAM - DRAM, используемая как кэш для быстрого доступа к данным.
  * `MEMORY_TYPE_NVRAM`:
    * **Описание:** Non-Volatile RAM (NVRAM) - тип оперативной памяти, сохраняющей данные при отключении питания.

Тип `MemoryDeviceType`:

  * `MEMORY_DEVICE_TYPE_UNSPECIFIED`:
    * **Описание:** Невалидное значение.
  * `MEMORY_DEVICE_TYPE_DDR`:
    * **Описание:** Double Data Rate (DDR) - тип DRAM, использующий двойное количество данных на такт.
  * `MEMORY_DEVICE_TYPE_LPDDR4_SDRAM`:
    * **Описание:** Low Power DDR4 SDRAM (LPDDR4 SDRAM) - энергоэффективный тип DDR4 SDRAM.
  * `MEMORY_DEVICE_TYPE_DDR_SDRAM`:
    * **Описание:** DDR Synchronous Dynamic Random Access Memory (DDR SDRAM) - синхронный динамический RAM.
  * `MEMORY_DEVICE_TYPE_DDR2`:
    * **Описание:** Второе поколение Double Data Rate (DDR2) SDRAM.
  * `MEMORY_DEVICE_TYPE_DR3_SDRAM`:
    * **Описание:** Третье поколение Double Data Rate (DDR3) SDRAM.
  * `MEMORY_DEVICE_TYPE_ROM`:
    * **Описание:** Read-Only Memory (ROM) - память только для чтения.
  * `MEMORY_DEVICE_TYPE_DDR3`:
    * **Описание:** Четвертое поколение Double Data Rate (DDR4) SDRAM.
  * `MEMORY_DEVICE_TYPE_LPDDR3_SDRAM`:
    * **Описание:** Low Power DDR3 SDRAM (LPDDR3 SDRAM) - энергоэффективный тип DDR3 SDRAM.
  * `MEMORY_DEVICE_TYPE_SDRAM`:
    * **Описание:** Synchronous Dynamic Random Access Memory (SDRAM) - синхронная динамическая оперативная память.
  * `MEMORY_DEVICE_TYPE_DDR4`:
    * **Описание:** Четвертое поколение Double Data Rate (DDR4) SDRAM.
  * `MEMORY_DEVICE_TYPE_DDR2_SDRAM`:
    * **Описание:** Второе поколение Double Data Rate (DDR2) SDRAM.
  * `MEMORY_DEVICE_TYPE_EDO`:
    * **Описание:** Extended Data Output (EDO) RAM.
  * `MEMORY_DEVICE_TYPE_DDR5`:
    * **Описание:** Пятое поколение Double Data Rate (DDR5) SDRAM.
  * `MEMORY_DEVICE_TYPE_DDR2_SDRAM_FB_DIMM`:
    * **Описание:** DDR2 SDRAM Fully Buffered DIMM.
  * `MEMORY_DEVICE_TYPE_FAST_PAGE_MODE`:
    * **Описание:** Fast-page mode.
  * `MEMORY_DEVICE_TYPE_DDR4_SDRAM`:
    * **Описание:** Четвертое поколение Double Data Rate (DDR4) SDRAM.
  * `MEMORY_DEVICE_TYPE_DDR2_SDRAM_FB_DIMM_PROBE`:
    * **Описание:** DDR2 SDRAM Fully Buffered DIMM Probe.
  * `MEMORY_DEVICE_TYPE_PIPELINED_NIBBLE`:
    * **Описание:** Pipelined Nibble.
  * `MEMORY_DEVICE_TYPE_DDR4_E_SDRAM`:
    * **Описание:** DDR4 E SDRAM.
  * `MEMORY_DEVICE_TYPE_DDR_SGRAM`:
    * **Описание:** DDR SGRAM.
  * `MEMORY_DEVICE_TYPE_LOGICAL`:
    * **Описание:** Логическая.
  * `MEMORY_DEVICE_TYPE_CDRAM`:
    * **Описание:** Cached RAM.
  * `MEMORY_DEVICE_TYPE_EDRAM`:
    * **Описание:** Extended Data RAM.
  * `MEMORY_DEVICE_TYPE_VRAM`:
    * **Описание:** Video RAM.
  * `MEMORY_DEVICE_TYPE_RAM`:
    * **Описание:** RAM.
  * `MEMORY_DEVICE_TYPE_EEPROM`:
    * **Описание:** Electrically Erasable Programmable Read-Only Memory (EEPROM) - электрически стираемая программируемая ПЗУ.
  * `MEMORY_DEVICE_TYPE_FEPROM`:
    * **Описание:** Flash EEPROM.
  * `MEMORY_DEVICE_TYPE_EPROM`:
    * **Описание:** Erasable Programmable Read-Only Memory (EPROM) - программируемая ПЗУ.

Тип `BaseModuleType`:

  * `BASE_MODULE_TYPE_UNSPECIFIED`:
    * **Описание:** Невалидное значение.
  * `BASE_MODULE_TYPE_RDIMM`:
    * **Описание:** Registered DIMM (RDIMM) - Тип DIMM-модуля, который использует регистры для буферизации адресов и команд памяти.
  * `BASE_MODULE_TYPE_UDIMM`:
    * **Описание:** Unbuffered DIMM (UDIMM) - DIMM-модуль без использования регистров для буферизации адресов и команд.
  * `BASE_MODULE_TYPE_SO_DIMM`:
    * **Описание:** Small Outline DIMM (SO-DIMM) - Компактный формат DIMM, обычно используемый в ноутбуках и других мобильных устройствах.
  * `BASE_MODULE_TYPE_LRDIMM`:
    * **Описание:** Load-Reduced DIMM (LRDIMM) - Тип DIMM-модуля, который использует буфер для уменьшения нагрузки на память.
  * `BASE_MODULE_TYPE_MINI_RDIMM`:
    * **Описание:** Mini Registered DIMM - Мини-версия Registered DIMM.
  * `BASE_MODULE_TYPE_MINI_UDIMM`:
    * **Описание:** Mini Unbuffered DIMM - Мини-версия Unbuffered DIMM.
  * `BASE_MODULE_TYPE_SO_RDIMM_72b`:
    * **Описание:** Small Outline Registered DIMM (72-bit) - Компактный формат Registered DIMM с 72-битной шириной шины данных.
  * `BASE_MODULE_TYPE_SO_UDIMM_72b`:
    * **Описание:** Small Outline Unbuffered DIMM (72-bit) - Компактный формат Unbuffered DIMM с 72-битной шириной шины данных.
  * `BASE_MODULE_TYPE_SO_DIMM_16b`:
    * **Описание:** Small Outline DIMM (16-bit) - Компактный формат DIMM с 16-битной шириной шины данных.
  * `BASE_MODULE_TYPE_SO_DIMM_32b`:
    * **Описание:** Small Outline DIMM (32-bit) - Компактный формат DIMM с 32-битной шириной шины данных.
  * `BASE_MODULE_TYPE_SIP`:
    * **Описание:** Single Inline Package (SIP) - Устройство в одном корпусе.
  * `BASE_MODULE_TYPE_DIP`:
    * **Описание:** Dual Inline Package (DIP) - Устройство в двухстрочном корпусе.
  * `BASE_MODULE_TYPE_ZIP`:
    * **Описание:** Zigzag In-line Package (ZIP) - Устройство с корпусом Zigzag In-line Package.
  * `BASE_MODULE_TYPE_SOJ`:
    * **Описание:** Small Outline J-lead (SOJ) - Компактный формат с корпусом Small Outline J-lead.
  * `BASE_MODULE_TYPE_PROPRIETARY`:
    * **Описание:** Проприетарный тип.
  * `BASE_MODULE_TYPE_SIMM`:
    * **Описание:** Single Inline Memory Module (SIMM) - Модуль памяти в одной линии.
  * `BASE_MODULE_TYPE_DIMM`:
    * **Описание:** Dual Inline Memory Module (DIMM) - Модуль памяти в двухстрочном форм-факторе.
  * `BASE_MODULE_TYPE_TSOP`:
    * **Описание:** Thin Small Outline Package (TSOP) - Тонкий компактный корпус с выводами.
  * `BASE_MODULE_TYPE_PGA`:
    * **Описание:** Pin Grid Array (PGA) - Массив выводов в форме сетки.
  * `BASE_MODULE_TYPE_RIMM`:
    * **Описание:** Rambus Inline Memory Module (RIMM) - Модуль памяти для технологии Rambus.
  * `BASE_MODULE_TYPE_SRIMM`:
    * **Описание:** Single Rambus Inline Memory Module (SRIMM) - Одиночный модуль памяти для технологии Rambus.
  * `BASE_MODULE_TYPE_SMD`:
    * **Описание:** Surface Mount Device (SMD) - Устройство для монтажа на поверхности.
  * `BASE_MODULE_TYPE_SSMP`:
    * **Описание:** Shrink Small Outline Package (SSMP) - Компактный корпус с уменьшенными размерами.
  * `BASE_MODULE_TYPE_QFP`:
    * **Описание:** Quad Flat Package (QFP) - Четырехсторонний корпус.
  * `BASE_MODULE_TYPE_TQFP`:
    * **Описание:** Thin Quad Flat Package (TQFP) - Тонкий четырехсторонний корпус
  * `BASE_MODULE_TYPE_SOIC`:
    * **Описание:** Small Outline Integrated Circuit (SOIC) - Компактный корпус для интегральных схем.
  * `BASE_MODULE_TYPE_LCC`:
    * **Описание:** Leadless Chip Carrier (LCC) - Корпус для монтажа чипов без выводов.
  * `BASE_MODULE_TYPE_PLCC`:
    * **Описание:** Plastic Leaded Chip Carrier (PLCC) - Корпус для монтажа чипов с пластиковыми выводами.
  * `BASE_MODULE_TYPE_BGA`:
    * **Описание:** Ball Grid Array (BGA) - Массив шаров для монтажа на поверхности.
  * `BASE_MODULE_TYPE_FPBGA`:
    * **Описание:** Fine-Pitch Ball Grid Array (FPBGA) - Массив шаров с мелким шагом для монтажа на поверхности.
  * `BASE_MODULE_TYPE_LGA`:
    * **Описание:** Land Grid Array (LGA) - Массив выводов для монтажа на поверхности.

Тип `MemoryState`:

  * `MEMORY_STATE_UNSPECIFIED`:
    * **Описание:** Невалидное значение.
  * `MEMORY_STATE_UNKNOWN`:
    * **Описание:** Состояние памяти неизвестно. Нельзя точно определить текущее состояние.
  * `MEMORY_STATE_OK`:
    * **Описание:** Память в нормальном состоянии. Нет проблем или ошибок.
  * `MEMORY_STATE_WARNING`:
    * **Описание:** Предупреждение относительно состояния памяти. Возможны проблемы, но они не критические.
  * `MEMORY_STATE_CRITICAL`:
    * **Описание:** Критическое состояние памяти. Присутствуют серьезные проблемы или ошибки, которые требуют внимания и решения.

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

Так как сообщения для операций идут по Protocol Buffer, необходимо иметь возможность скомпилировать protobuf-файлы под конкретный язык программирования.

Подробнее об этом:

- https://github.com/protocolbuffers/protobuf?tab=readme-ov-file

Для решения этой задачи был сформирован Makefile, который находится в корне проекта. При его выполнении происходит создание необходимых файлов.

## Пример реализации

Реализация операции будет производиться на виртуальной машине, запущенной c помощью программы по виртуализации **Proxmox** со следующими характеристикам:

- Название ОС:	**`Microsoft Windows Server 2019 Standard`**
- Версия:	**`10.0.17763 Build 17763`**
- Процессор:	**`Common KVM processor, 2095 Mhz, 2 Core(s), 2 Logical Processor(s)`**
- Объем памяти дисков: **`100,00 GB`**
- Суммарно установлено памяти (ОЗУ):	**`4,00 GB`**
- Количество плашек ОЗУ: **`1`**

Для реализации сбора данных необходимо иметь возможность подключения к удаленному хосту по протоколу WinRM, для этого необходимо настроить подключение на удаленной машине.

Подробнее по ссылке:
- https://learn.microsoft.com/ru-ru/windows/win32/winrm/installation-and-configuration-for-windows-remote-management

Имея возможность подключения к удаленной машине с ОС Windows мы можем отправлять консольные команды на запуск определенных утилит. В основом это утилита **`WMIC`**.

**`WMIC`** (Windows Management Instrumentation Command-line) — мощная утилита командной строки, которая используется для работы с WMI (Windows Management Instrumentation) в операционных системах Windows. Она предоставляет различные возможности для управления системой и сбора информации о компьютере.

Подробнее по ссылкам:
 - https://learn.microsoft.com/ru-ru/windows/win32/wmisdk/wmi-start-page
 - https://learn.microsoft.com/ru-ru/windows/win32/wmisdk/wmic

Шаблон команды для выполнения - `WMIC PATH @Class GET @Fields /format:list`, где:
 - **@Class** - Имя класса WMI
 - **@Fields** - Список атрибутов, которые необходимо получить

Используемые WMI классы и атрибуты для сбора инвентарных данных по ОЗУ:
 - **[Win32_PhysicalMemory](https://learn.microsoft.com/en-us/windows/win32/cimwin32prov/win32-physicalmemory)** - **`BankLabel, TypeDetail, MemoryType, SMBIOSMemoryType, FormFactor, Manufacturer, Capacity, PartNumber, SerialNumber, Status, Speed, DeviceLocator`**

После формирования команды 

Пример вывода команды **`WMIC PATH Win32_PhysicalMemory get /format:list`** без фильтра по полям:

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

Полученные данные парсятся и возвращаются в качестве ответа по данной RPC.

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

        // ...
      
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