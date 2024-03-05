# Название сценария
Пример реализации системной операции "Получить данные о нагрузке на CPU". 
Cсылка на документацию(https://docs.bergen.tech/ems/release-documents/latest/#/specifications/ds/host-domain/linux-manager/README?id=%d0%9f%d0%be%d0%bb%d1%83%d1%87%d0%b8%d1%82%d1%8c-%d0%b4%d0%b0%d0%bd%d0%bd%d1%8b%d0%b5-%d0%be-%d0%bd%d0%b0%d0%b3%d1%80%d1%83%d0%b7%d0%ba%d0%b5-%d0%bd%d0%b0-cpu)
## Обзор операции
Данная операция отвечает за получение данных о нагрузке на CPU устройства, либо от агента, если он установлен на устройстве.

Для выполнения этой операции сервис на вход должен получить следующий запрос:

* `CollectLinuxCpuUtilisationRequest`:
  * `device`:
    * **Тип параметра:** `DeviceContent`
    * **Описание:** Данные по 1 устройству.

Тип `DeviceContent`:

  * `device_id`:
    * **Тип параметра:** `string`
    * **Описание:** Идентификатор устройства.
  * `model_name`:
    * **Тип параметра:** `string`
    * **Описание:** Название модели устройства.
  * `vendor_name`:
    * **Тип параметра:** `string`
    * **Описание:** Название вендора устройства.
  * `connectors`:
    * **Тип параметра:** `RepeatedField<DeviceConnector>`
    * **Описание:** Идентификатор сетевого интерфейса.

Тип `DeviceConnector`:

* `device_network_id`:
  * **Тип параметра:** `string`
  * **Описание:** Идентификатор сетевого интерфейса устройства.
* `address`:
  * **Тип параметра:** `string`
  * **Описание:** IP/FQDN адрес устройства.
* `mac`:
  * **Тип параметра:** `string`
  * **Описание:** MAC-адрес устройства.
* `credentials`:  
  * **Тип параметра:** `RepeatedField<Credential>`
  * **Описание:** Учетные данные подключения.

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
    * **Тип параметра:** `uint32`
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
  * `privateKey`:
    * **Тип параметра:** `string`
    * **Описание:** Private key (только для SNMP).
  * `security_level`:
    * **Тип параметра:** `string`
    * **Описание:** Уровень безопасности.
	
Enum `ConnectorProtocol`:
* `ConnectorProtocol` (0):
	* **Описание:** Невалидное значение.
* `CONNECTOR_PROTOCOL_IPMI` (1):
	* **Описание:** Ipmi протокол для проверки подключения.
* `CONNECTOR_PROTOCOL_REDFISH` (2):
	* **Описание:** Redfish протокол для проверки подключения.
* `CONNECTOR_PROTOCOL_SNMP` (3):
	* **Описание:** Snmp протокол для проверки подключения.
* `CONNECTOR_PROTOCOL_SSH` (4):
	* **Описание:** Ssh протокол для проверки подключения.
* `CONNECTOR_PROTOCOL_WMI` (5):
	* **Описание:** Wmi протокол для проверки подключения.


Из даных запроса, сервис интересуют следующие поля:
*	DeviceContent->connectors - список возможэных сетевых подключений
*	DeviceContent->connectors->address - Адрес сервера
*	DeviceContent->connectors->credentials - Список данных, необходимых для подключения
*	DeviceContent->connectors->protocol - Протокол подключения. Обязательно должен быть CONNECTOR_PROTOCOL_SSH. Другие протоколы должны игнорироваться
*	DeviceContent->connectors->login - Логин пользователя, под которым будет выполнятсья команда
*	DeviceContent->connectors->password - Пароль пользователя, под которым будет выполнятсья команда
*	DeviceContent->connectors->port - Порт подключения

Используя полученные данные, сервис выполняет на оборудовании следующую команду:
*	top -b -n1 -1 -p0 -w 400", "top -b -n1 -w 400

Во время вызова сначала производится подключение к NATS, опрос его на наличие Jetstream, принадлежащих агентам целевого устройства.
Если агент найден, то запрос выполняется через него. Если не найден или агент вернул ошибку, то запрос выполняется по протоколу SSH.
Затем обрабатывает полученные результаты и возвращает ответ на grpc запрос в виде следующего объекта:

Тип `DeviceCpuUtilisation`:

* `device_identity`:
  * **Тип параметра:** `DeviceDataIdentity`
  * **Описание:** Описание источника сбора данных.
* `unit_utilistaions`:
  * **Тип параметра:** `map<int32, CpuUnitUtilisation`
  * **Описание:** Данные о статусе каждого из процессоров/ядер. Ключ - id процессора.
* `summary_utilisation`:
  * **Тип параметра:** `CpuUnitUtilisation`
  * **Описание:** Суммарные метрики потребления CPU.

Тип `CpuUnitUtilisation`:

* `total_using`:
  * **Тип параметра:** `int32`
  * **Описание:** Общий процент использования процессора.
* `idle_time`:
  * **Тип параметра:** `int32`
  * **Описание:** Время простоя процессора, выраженное в процентах.
* `user_using`:
  * **Тип параметра:** `google.protobuf.Int32Value`
  * **Описание:** Процент использования процессора пользовательскими процессами.
* `system_using`:
  * **Тип параметра:** `google.protobuf.Int32Value`
  * **Описание:** Процент использования процессора системными процессами.
* `nice_value_using`:
  * **Тип параметра:** `google.protobuf.Int32Value`
  * **Описание:** Процент времени, в течение которого CPU выполнял процессы, выставленные пользователем вручную (nice).
* `io_waiting`:
  * **Тип параметра:** `google.protobuf.Int32Value`
  * **Описание:** Процент времени, потраченного на ожидание ввода-вывода.
* `hw_service_interrupts`:
  * **Тип параметра:** `google.protobuf.Int32Value`
  * **Описание:** Процент времени, потраченного на обработку аппаратных прерываний.
* `soft_service_interrupts`:
  * **Тип параметра:** `google.protobuf.Int32Value`
  * **Описание:** Процент времени, потраченного на обработку программных (системных) прерываний.
* `steal_time`:
  * **Тип параметра:** `google.protobuf.Int32Value`
  * **Описание:** Процент времени, потраченного на выполнение задач в виртуальной машине (виртуализация).


## Пример реализации
В качестве примера реализуем простейшее приложение, реализующее только систумную операцию "Получить данные о нагрузке на CPU".
В качестве опрашиваемого устройства подойдет любое устройство на linux, у которого открыт протокол ssh. Вполне подойдет виртуалка на WSL.
*	Создайте проект из шаблона "Служба ASP.NET Core gRPC"
*	Удалите демонстрационные файлы из директории Protos и скопируйте туда proto файлы из директории project\SshExample\SshExample\Protos. Это уже подготовленные grpc объекты, описанные выше
*	В директории service создайте файл SshCommandCaller.cs
*	Добавьте в этот файл следующие константы:

		private static readonly string[] _commands = { "top -b -n1 -1 -p0 -w 400", "top -b -n1 -w 400" };
		private static string TopCpuPrefix = "%cpu";
		
	_commands - список shell команд, которые будет запрашивать сервис у оборудования.
	TopCpuPrefix - префикс для более удобной обработки результатов
*	Добавьте в SshCommandCaller.cs следующий record

		public record HandleResult(int Success, string Stdout, string Stderr);
		
	Этот объект будет содержать результаты выполнения ssh запроса

*	Добавьте метод:

		public static HandleResult CallSsh(CollectLinuxCpuUtilisationRequest request, string command)
		{
			DeviceConnector connection = request.Device.Connectors.First();
			Credential creds=connection.Credentials.First();

			using (SshClient client = new SshClient(connection.Address, creds.Port, creds.Login, creds.Password))
			{
				client.Connect();
				SshCommand cmdRes = client.RunCommand(command);
				if (cmdRes.ExitStatus == 0)
				{
					client.Disconnect();
					return new HandleResult(cmdRes.ExitStatus,
					cmdRes.Result, cmdRes.Error);
				}
				client.Disconnect();
				return new HandleResult(cmdRes.ExitStatus,
					cmdRes.Result, cmdRes.Error);
			}
		}
		
	Этот метод не учитывает случай, что пользователь может не передать данных для подсоединения	или у устройства их может быть несколько.
	Этот метод просто берет первые возможные данные для подсоединения, подключается к устройству, выполняет команду и собирает 
	объект HandleResult.
	
*	Добавьте метод:
	
		string response = "";
		foreach (HandleResult res in results)
		{
			if (res.Success == 0)
			{
				response += res.Stdout;
			}
		}

		CollectLinuxCpuUtilisationResponse statDeviceCpu = new()
		{
			CpuUtilisation = new()
			{
				DeviceIdentity = new DeviceDataIdentity()
				{
					DeviceId = request.Device.DeviceId,
					Source = ServiceSource.LinuxManager
				},
				SummaryUtilisation =new()
			}
		};

		IEnumerable<string> rows = response.ToLowerInvariant()
			.Split('\n', StringSplitOptions.TrimEntries)
			.Where(s => s.StartsWith(TopCpuPrefix))
			.SelectMany(c => c
				.Split(TopCpuPrefix, StringSplitOptions.TrimEntries)
				.Where(d => !string.IsNullOrWhiteSpace(d))
				.Select(s => TopCpuPrefix + s));

		foreach (string item in rows)
		{
			int processorIdEnd = item.IndexOf(' ', TopCpuPrefix.Length);

			if (processorIdEnd == -1)
			{
				continue;
			}

			string? strProcessorId = item[TopCpuPrefix.Length..processorIdEnd];

			int processorId;

			if (strProcessorId == "(s):")
			{
				processorId = -1;
			}
			else if (!int.TryParse(strProcessorId, out processorId))
			{
				continue;
			}

			foreach (string entiry in item[(processorIdEnd + 3)..].Split(',', StringSplitOptions.TrimEntries))
			{

				string[] parts = entiry.Split(' ', 2, StringSplitOptions.TrimEntries);

				if (parts.Length != 2)
				{
					continue;
				}

				string key = parts[1];
				string strValue = parts[0];

				if (float.TryParse(strValue, out float intValue))
				{
					switch (key)
					{
						case "us": statDeviceCpu.CpuUtilisation.SummaryUtilisation.UserUsing = (int)Math.Ceiling(intValue); break;
						case "sy": statDeviceCpu.CpuUtilisation.SummaryUtilisation.SystemUsing = (int)Math.Ceiling(intValue); break;
						case "ni": statDeviceCpu.CpuUtilisation.SummaryUtilisation.NiceValueUsing = (int)Math.Ceiling(intValue); break;
						case "id": statDeviceCpu.CpuUtilisation.SummaryUtilisation.IdleTime = (int)Math.Ceiling(intValue); break;
						case "wa": statDeviceCpu.CpuUtilisation.SummaryUtilisation.IoWaiting = (int)Math.Ceiling(intValue); break;
						case "hi": statDeviceCpu.CpuUtilisation.SummaryUtilisation.HwServiceInterrupts = (int)Math.Ceiling(intValue); break;
						case "si": statDeviceCpu.CpuUtilisation.SummaryUtilisation.SoftServiceInterrupts = (int)Math.Ceiling(intValue); break;
						case "st": statDeviceCpu.CpuUtilisation.SummaryUtilisation.StealTime = (int)Math.Ceiling(intValue); break;
						default: break;
					}
				}
			}
		}
		return statDeviceCpu;
		
	Этот метод должен заниматься анализом полученных в результате выполнения команд данных и сборкой ответа для пользователя.
	На вход метода должны поступать примерно такие данные:
	
		top - 17:55:32 up 1 min,  1 user,  load average: 0.01, 0.01, 0.00
		Tasks:   1 total,   1 running,   0 sleeping,   0 stopped,   0 zombie
		%Cpu0  :  0.0 us,  0.0 sy,  0.0 ni,100.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st     
		%Cpu1  :  0.0 us,  0.0 sy,  0.0 ni,100.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st

		...
		   
		%Cpu19 :  0.0 us,  0.0 sy,  0.0 ni,100.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
		MiB Mem :   7806.4 total,   6860.3 free,    509.2 used,    436.8 buff/cache
		MiB Swap:   2048.0 total,   2048.0 free,      0.0 used.   7067.7 avail Mem

			PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND
			697 vvsurje+  20   0    7656   3236   2872 R   0.0   0.0   0:00.00 top
			
		top - 17:55:52 up 2 min,  1 user,  load average: 0.01, 0.00, 0.00
		Tasks:  35 total,   1 running,  34 sleeping,   0 stopped,   0 zombie
		%Cpu(s):  0.0 us,  0.0 sy,  0.0 ni,100.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
		MiB Mem :   7806.4 total,   6858.4 free,    511.2 used,    436.9 buff/cache
		MiB Swap:   2048.0 total,   2048.0 free,      0.0 used.   7065.8 avail Mem

			PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND
			690 root      20   0   43696  37724   9956 S  13.3   0.5   0:00.75 python3
			  1 root      20   0  166004  11392   8352 S   0.0   0.1   0:00.21 systemd	
			  ...
			  
	Метод объединяет оба результата в одну строку. Затем разделяет ее на под-строки по символу '\n' (перевод строки).
	Затем отфильтровывает все лишнее, оставляя только строки, начинающиеся с %Cpu.
	Каждая строка - это статистика отдельного процессора. %Cpu(s) - общие, сводные, данные по всем процессорам.
	Далее метод разделяет строки на блоки по символу ',' и обрабатывает каждый блок отдельно.
	Результатом работы метода является объект CollectLinuxCpuUtilisationResponse со сводными данными по всем процессорам.
	
*	Добавьте метод:

		public static CollectLinuxCpuUtilisationResponse GetCpuUtilisation(CollectLinuxCpuUtilisationRequest request)
		{
			HandleResult[] responses= _commands.Select(cmd=> CallSsh(request,cmd)).ToArray();
			return ProcessResponse(responses, request);
		}
		
	Это управляющий метод, который сначала вызовет на устройстве команды, а затем передаст результаты в обработку.	
	
*	В файле GreeterService.cs (который был создан из шаблона) удалите все содержимое и добавьте следующее:

	namespace SshExample.Services;

	public class LinuxManagerService : LinuxManager.LinuxManagerBase
	{
		public override Task<CollectLinuxCpuUtilisationResponse> CollectCpuUtilisation(CollectLinuxCpuUtilisationRequest request, ServerCallContext context)
		{
			return Task.FromResult(SshCommandCaller.GetCpuUtilisation(request));
		}
	}

	Это обработчик Grpc запросов, который передает уаравление в класс SshCommandCaller и метод GetCpuUtilisation.
	
*	После сборки и запуска приложения, можно проверять результат работы. Для этого нужно вызвать метод CollectCpuUtilisation, например, в Postman.