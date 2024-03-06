# Название сценария

## Обзор операции

Данная операция необходима для определения списка гипервизоров, заведенных в системе.

Proto-контракт для данной операции находится [здесь](project/proto/service_hypervisor_manager.proto).

## Пример реализации

Реализацию можно посмотреть в директории [project](project/main.go).

Дополняем уже имеющийся [шаблон](../create_project/project/main.go) релизацией RPC `CollectVirtialMachinesList`.

В ходе обработки запроса на языке golang для гипервизора `ESXI` мы воспользуемся библиотекой [govmomi](https://github.com/vmware/govmomi) для взаимодействия с VMware vSphere APIs и сбора информации с гипервизора.

Во входящем запросе мы получаем креды и адрес для подключения ESXI. Используем их и создаем клиента:

```golang
address := req.Hypervisor.GetAddress()

// Парсинг адреса.
u, err := soap.ParseURL(address)
if err != nil {
    return nil, fmt.Errorf("parse URL [%s]: %s", address, err)
}

// Установка кредов.
u.User = url.UserPassword(req.Hypervisor.GetLogin(), req.Hypervisor.GetPassword())

// Создание клиента.
client, err := govmomi.NewClient(context.Background(), u, true)
if err != nil {
    return nil, fmt.Errorf("create client: %s", err)
}
```

Далее собираем информацию с гипервизора:

```golang
// Сбор данных с гипервизора, используя библиотеку govmomi.
virtualMachines := []mo.VirtualMachine{}
m := view.NewManager(client.Client)
v, err := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
if err != nil {
    return nil, fmt.Errorf("create container view: %s", err)
}
defer v.Destroy(ctx)

if err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"guest", "config", "summary"}, virtualMachines); err != nil {
    return nil, fmt.Errorf("retrieve: %s", err)
}
```

Далее производим парсинг полученной информации и заполнение выходной структуры:

```golang
for _, vm := range virtualMachines {
    machine := pb.VirtualMachine{
        Name: vm.Summary.Config.Name,
    }

    for _, ni := range vm.Guest.Net {
        if ni.Network != "" {

            machine.Networks = append(machine.Networks, &pb.VirtualMachineNetwork{
                Ips: getIPv4(ni.IpAddress, regexp),
                Mac: ni.MacAddress,
            })
        }
    }

    resp.VirtualMachines.VirtualMachines = append(resp.VirtualMachines.VirtualMachines, &machine)
}
```

В завершении отправляем полученный результат как ответ по RPC.
