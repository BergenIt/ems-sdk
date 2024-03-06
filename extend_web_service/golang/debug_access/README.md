# Название сценария

Данный документ описывает подход к разработке модуля расширения для RPC `DebugAccess`.

## Обзор операции

Данная операция необходима для проверки доступности инстанса сервиса при заведении в систему в зависимости от протокола.

Proto-контракт для данной операции находится [здесь](project/proto/service_service_manager.proto).

## Пример реализации

Реализацию можно посмотреть в директории [project](project/main.go).

Дополняем уже имеющийся [шаблон](../create_project/project/main.go) релизацией RPC `DebugAccess`.

В ходе обработки запроса нам необходимо проверить доступность адреса сервиса в сети в завизиости от пяти доступных протоколов:

* TCP
* UDP
* WS
* HTTP
* gRPC

Сделаем проверку следующим образом, если протокол в запросе будет `TCP`, `WS`, `HTTP`, `gPRC`, то будем проверять адрес на доступность по протоколу `TCP`, если протокол в запросе `UDP` то по протоколу `UDP`. Это обусловлено тем, что низкоуровнево протоколы `WS`, `HTTP` и `gPRC` работают по протоколу `TCP`, а значит если этот адрес доступен по `TCP`, то можно пробовать производить его мониторинг. Соответственно и с `UDP`.

На вход нам приходит `oneof` адреса, который может быть либо доменным именем, либо ip-адресом с портом.

Составляем адрес в зависимости от `oneof` поля Address:

```golang
address := ""
switch f := req.Address.(type) {
case *pb.DebugServiceAccessRequest_AddressPort:
    address = fmt.Sprintf("%s:%d", f.AddressPort.Address, f.AddressPort.Port)
case *pb.DebugServiceAccessRequest_Uri:
    address = f.Uri
}
```

Далее производим проверку в зависимости от протокола:

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
