# Проверка доступности инстанса сервиса

Данная операция необходима для определения доступности инстанса сервиса при заведении системы.

## Обзор операции

Данная операция реализуется следующим RPC:

```proto
service ServiceManager {
  rpc DebugAccess(DebugServiceAccessRequest) returns (DebugServiceAccessResponse);
}

// Контракт запроса для rpc DebugAccess
message DebugServiceAccessRequest {
  // URI сервиса
  string uri = 1;
  // Протокол сервиса
  ServiceProtocol protocol = 2;
}

// Контракт ответа для rpc DebugAccess
message DebugServiceAccessResponse {
  // Результат проверки
  DebugAccessResult result = 1;
}
```

С полной структурой данных вы можете ознакомиться в [протофайлах](../../../.proto/service_hypervisor_manager.proto).

## Пример реализации

Дополняем уже имеющийся [шаблон](../create_project/project/main.go) релизацией RPC `DebugAccess`.

В ходе обработки запроса нам необходимо проверить доступность адреса сервиса в сети в завизиости от пяти доступных протоколов:

* TCP
* UDP
* WS
* HTTP
* gRPC

Сделаем проверку следующим образом, если протокол в запросе будет `TCP`, `WS`, `HTTP`, `gPRC`, то будем проверять адрес на доступность по протоколу `TCP`, если протокол в запросе `UDP` то по протоколу `UDP`. Это обусловлено тем, что низкоуровнево протоколы `WS`, `HTTP` и `gPRC` работают по протоколу `TCP`, а значит если этот адрес доступен по `TCP`, то можно пробовать производить его мониторинг. Соответственно и с `UDP`.

На вход нам приходит `oneof` адреса, который может быть либо доменным именем, либо ip-адресом с портом.

Производим проверку в зависимости от протокола:

```golang
switch req.Protocol {
case pb.ServiceProtocol_SERVICE_PROTOCOL_GRPC,
    pb.ServiceProtocol_SERVICE_PROTOCOL_WS,
    pb.ServiceProtocol_SERVICE_PROTOCOL_TCP,
    pb.ServiceProtocol_SERVICE_PROTOCOL_HTTP:
    availability, err = pingTcp(address)

    out.Result.State = determineAvailability(availability, err)

case pb.ServiceProtocol_SERVICE_PROTOCOL_UDP:
    availability, err = pingUdp(address)

    out.Result.State = determineAvailability(availability, err)

default:
    return nil, fmt.Errorf("unsuppoted protocol: %v", req.Protocol)
}
```

С помощью данной функции определяем статус адреса:

```golang
    // Определение статуса сервиса.
    func determineAvailability(availability bool, err error) pb.ServiceAvailableState {
        if err == nil && availability {
            return pb.ServiceAvailableState_SERVICE_AVAILABLE_STATE_AVAILABLE
        } else if err != nil {
            log.Printf("determine service availability: %s", err)
            return pb.ServiceAvailableState_SERVICE_AVAILABLE_STATE_UNAVAILABLE
        } else {
            return pb.ServiceAvailableState_SERVICE_AVAILABLE_STATE_UNAVAILABLE
        }
    }
```

В завершении отправляем полученный результат как ответ по RPC.

Пример готового проекта расположен в папке [project](./project)