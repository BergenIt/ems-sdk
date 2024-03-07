# Обновление прошивок BMC

После формирования скелета проекта можно перейти к наполнению логикой операции обновления прошивок BMC модуля расширения Bmc.

В данном разделе описана логика сбора информации в рамках этой системной операции.

## Обзор операции

Данная операция отвечает за процесс по обновлению BMC с транспортировкой по протоколу `Redfish`.

Для реализации операции RPC будет иметь следующую сигнатуру:

* **`rpc BmcFirmwareUpdate(BmcFirmwareUpdateRequest) returns (BmcFirmwareUpdateResponse)`**

На вход модуль получает следующий запрос в формате Protocol buffers:

Тип `BmcFirmwareUpdateRequest`:

* `device`:
  * **Тип параметра:** `DeviceContent`
  * **Описание:** Данные об устройстве.
* `firmware_url`:
  * **Тип параметра:** `string`
  * **Описание:** Ссылка на файл прошивки.
* `mode`:
  * **Тип параметра:** `BmcFirmwareUpdateMode`
  * **Описание:** Режим обновления прошивки BMC.

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

Перечисление `BmcFirmwareUpdateMode`:

* `BMC_FIRMWARE_UPDATE_MODE_UNSPECIFIED`:
  * **Описание:** Невалидное значение.
* `BMC_FIRMWARE_UPDATE_MODE_API`:
  * **Описание:** Обновление прошивки BMC через Redfish.
* `BMC_FIRMWARE_UPDATE_MODE_RAW_IPMI`:
  * **Описание:** Обновление прошивки BMC через IPMI.

Данная структура запроса является общей для реализации операции обновления прошивки BMC по разным протоколам, поэтому может содержать большее количество полей, чем поддерживает Bmc.

Для корректной работы сбора данных, устройство должно иметь хотя бы одно действительное подключение с протоколом Redfish (поле **`protocol`** списка **`credentials`** со значением **`CONNECTOR_PROTOCOL_REDFISH`**), остальные устройства модуль должен игнорировать.

Пример данных запроса для обновления прошивки BMC:

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

Тип `BmcFirmwareUpdateResponse`:

* `result`:
  * **Тип параметра:** `OperationResult`
  * **Описание:** Результат выполнения операции управления.

Тип `OperationResult`:

* `device_id`:
  * **Тип параметра:** `string`
  * **Описание:** Идентификатор устройства, на котором происходила операция.
* `state`:
  * **Тип параметра:** `OperationState`
  * **Описание:** Тип результата выполнения операции.
* `output`:
  * **Тип параметра:** `string`
  * **Описание:** Текстовое описание результата выполнения операции.

Перечисление `OperationState`:

* `OPERATION_STATE_UNSPECIFIED`:
  * **Описание:** Невалидное значение.
* `OPERATION_STATE_SUCCESS`:
  * **Описание:** Операция завершена успешно.
* `OPERATION_STATE_FAILED`:
  * **Описание:** Выполнение операции завершилось с ошибкой.

С полной структурой данных вы можете ознакомиться в [протофайлах](../../../.proto).

## Пример реализации

Реализация операции будет производиться на устройстве вендора `Huawei` со следующими характеристикам:

* Название ОС: **`-`**
* Версия: **`-`**
* Процессор: **`Intel(R) Xeon(R) Silver 4208 CPU @ 2.10GHz x 2`**
* Объем памяти дисков: **`-`**
* Суммарно установлено памяти (ОЗУ): **`24 GB`**
* Количество плашек ОЗУ: **`3`**

Для реализации операции необходимо иметь возможность подключения к удаленному хосту по протоколу Redish, для этого необходимо сформировать набор данных подключения.

Изучив [документацию](https://support.huawei.com/enterprise/en/doc/EDOC1100105856/745aefb6/managing-chassis-resources) с официального сайта вендора **`Huawei`** становится понятно каким образом хранятся данные по **`Redfish`**.

URI домашней страницы службы Redfish (другое название корневой адрес службы) можно получить с помощью следующего URI: `https://<IP-адрес:порт>/redfish/v1`

Запрос прошивки устройства по **`Redfish`** описан на официальной [странице](https://support.huawei.com/enterprise/en/doc/EDOC1100177343/e9b8568b/upgrading-firmware) вендора **`Huawei`**

Исходя из специфики операции необходимо нам необходимо иметь файл прошивки BMC для Redfish Huawei, который можно найти на официальной сайте вендора.

Также для того, чтобы целевое устройство могло получить файл прошивки, он должен быть расположен на SFTP сервере (в поставке EMS специально развернут SFTP сервер).

Конфигурация для выполнения операции выглядит следующим образом:
```yaml
  - vendor: Huawei
    model: 2288H V5
    simple_update:
      url: /redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate
      request: '{"ImageURI": "{filePath}", "TransferProtocol": "{protocol}"}'
```

Сам процесс обновления состоит создания сеанса пользователя, обновления прошивки и закрытия сессии:

```go
func (r *redfishService) huaweiUpdateFirmware(
	filePath string,
	host Host,
	cfg RedfishCFG,
	redfish *rfish.RedfishAdapter,
) error {
	if !strings.HasPrefix(filePath, SFTP) {
		return ErrHuaweiWrongFilePath
	}

	h, id, err := redfish.CreateSession(host.Headers, host.User, host.Password)
	if err != nil {
		return fmt.Errorf("CreateSession failed: %w", err)
	}

	_, ok := h[AuthHeader]
	if !ok {
		return ErrEmptyAuthToken
	}

	if err := r.updateFirmware(h, filePath, "@odata.id", "TaskState", cfg, redfish); err != nil {
		return fmt.Errorf("updateFirmware failed: %w", err)
	}

	if err := redfish.DeleteSession(h, id); err != nil {
		return fmt.Errorf("DeleteSession failed: X-Auth-Token: %s: Id: %q", h[AuthHeader], id)
	}

	return nil
}
```

Процесс создания сессии Redfish с устройством:

```go
func (ra *RedfishAdapter) CreateSession(headers map[string]string, username, password string) (map[string]string, string, error) {
	p := fmt.Sprintf("{\"UserName\": \"%s\", \"Password\": \"%s\"}", username, password)
	payload := make(map[string]any)
	if err := json.Unmarshal([]byte(p), &payload); err != nil {
		return nil, "", err
	}
	sessionURL := "/redfish/v1/SessionService/Sessions"
	resp, err := ra.client.PostWithHeaders(sessionURL, payload, headers)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	var result = make(map[string]any)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", err
	}

	var id string
	x, ok := result["Id"]
	if ok {
		id = x.(string)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, "", ErrCreateSession
	}

	authToken := resp.Header.Get(AuthHeader)
	if authToken == "" {
		return nil, "", ErrGetAuthToken
	}

	headers[AuthHeader] = authToken

	return headers, id, nil
}
```

Отправка Post-запроса на обновление прошивки имеет такой вид:

```go
func (ra *RedfishAdapter) PostWithHeaders(headers map[string]string, updateURL, body string) (*http.Response, error) {
	payload := make(map[string]any)
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return nil, fmt.Errorf("unmarshal payload: %w", err)
	}

	resp, err := ra.client.PostWithHeaders(updateURL, payload, headers)
	if err != nil {
		return nil, fmt.Errorf("error when try request to update: %w", err)
	}

	return resp, nil
}
```

После запуска процесса обновления прошивки BMC необходимо проконтролировать, чтобы процесс обновления прошёл обновления прошёл успешно. Для этого написана следующая функция:

```go
func (r *redfishService) monitoringUpdFirmware(taskMonitor, state string, redfish *rfish.RedfishAdapter) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	var getTaskErr error
	for {
		select {
		default:
			// каждый раз делаем запрос по пути мониторинга для получения свежей информации
			resp, err := redfish.GetTaskMonitor(taskMonitor)
			if err != nil {
				if getTaskErr == nil {
					getTaskErr = &RedfishRequestErr{Origin: err}
				}
				continue
			}

			switch resp.StatusCode {
			case http.StatusNoContent, http.StatusOK:
				if resp.ContentLength <= 0 {
					return nil
				}
			case http.StatusAccepted:
				continue
			}

			var result map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return fmt.Errorf("unmarshal response from task service: %w", err)
			}
			resp.Body.Close()

			states, ok := result[state].(string)
			if !ok {
				return fmt.Errorf(
					"failed assert states: states value: [%s] status: [%v] content lenght: [%v]",
					states,
					resp.StatusCode,
					resp.ContentLength,
				)
			}
			states = strings.ToLower(states)
			switch states {
			case "completed", "ok":
				return nil
			case "exception", "error", "critical":
				return fmt.Errorf("failed task for updating source: [%s]", result)
			}

		case <-ctx.Done():
			return fmt.Errorf("%w: %s", ctx.Err(), getTaskErr.Error())
		}
	}
}
```

Пример готового проекта расположен в папке [project](./project)
