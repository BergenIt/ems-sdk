# Сохранение конфигурации коммутатора

Данная операция необходима для сохранения текущей конфигурации данного коммутатора

## Обзор операции

Данная операция реализуется следующим RPC:

```proto
// Описание сервиса для мониторинга доступности и управления конфигурациями коммутаторов
service NetworkManager {
  // процедура для сохранения настроек коммутатора
  rpc CreateConfig(CreateNetworkConfigRequest) returns (CreateNetworkConfigResponse);
}

message CreateNetworkConfigRequest {
  DeviceContent device = 1;
}

message CreateNetworkConfigResponse {
  OperationResult result = 1;
}
```

С полной структурой данных вы можете ознакомиться в [протофайлах](../../../.proto/service_network_manager.proto).

## Пример реализации

Реализация операции будет производиться на Microtik RouterOS.

Дополняем уже имеющийся [шаблон](../create_project/project/main.go) релизацией RPC `CreateConfig`.

Во входящем запросе мы получаем креды и адрес для подключения к ОС коммутатора. Для начала нам нужно извлечь их.
Создадим удобную структуру для их хранения:

```golang
type sshConnInfo struct {
	addr  string
	login string
	pass  string
	port  int32
}
```

Теперь создадим функцию для извлечения этих кредов:

```golang
func extractSSHConnInfo(connectors []*pb.DeviceConnector) (sshConnInfo, error) {
	for _, conn := range connectors {
		for _, creds := range conn.Credentials {
			if creds.Protocol == pb.ConnectorProtocol_CONNECTOR_PROTOCOL_SSH {
				res := sshConnInfo{
					addr:  conn.Address,
					login: creds.Login,
					pass:  creds.Password,
					port:  creds.Port,
				}
				return res, nil
			}
		}
	}
	return sshConnInfo{}, fmt.Errorf("ssh credentials were not found")
}
```

Теперь мы можем использовать эту функцию для извлечения кредов и сознания клиента:

```golang
    info, err := extractSSHConnInfo(req.Device.Connectors)
	if err != nil {
		return nil, err
	}

	cfg := &ssh.ClientConfig{
		User: info.login,
		HostKeyCallback: ssh.HostKeyCallback(
			func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			}),
		Auth: []ssh.AuthMethod{
			ssh.Password(info.pass),
		},
	}

	conn, err := ssh.Dial("tcp", info.addr, cfg)
	if err != nil {
		return nil, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
```

Далее мы соединяем stdout ssh-сессии с байтовым буфером, и запускаем удаленную команду:

```golang
	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(createConfigCmd); err != nil {
		return nil, err
	}
```

В завершении, скомпонуем ответ для RPC, используя полученные данные:

```golang
	res := pb.CreateNetworkConfigResponse{
		Result: &pb.OperationResult{
			DeviceId: req.Device.DeviceId,
			State:    pb.OperationState_OPERATION_STATE_SUCCESS,
			Output:   b.String(),
		},
	}

	return &res, nil
```
