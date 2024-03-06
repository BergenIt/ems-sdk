# Сбор перечня виртуальных машин

Данная операция необходима для определения списка гипервизоров, заведенных в системе.

## Обзор операции

Данная операция реализуется следующим RPC:

```proto
// Сервис HypervisorManager реализует логику взаимодействия и сбора информации с гипервизоров
service HypervisorManager  {
  // Получение списка виртуальных машин гипервизора
  rpc CollectVirtialMachinesList(CollectVirtialMachinesListRequest) returns (CollectVirtialMachinesListResponse);
}

// Контракт запроса для rpc CollectVirtialMachinesList
message CollectVirtialMachinesListRequest {
  // Информация о гипервизоре
  HypervisorContent hypervisor = 1;
}

// Контракт ответа для rpc CollectVirtialMachinesList
message CollectVirtialMachinesListResponse {
  // Спискок виртуальных машин
  VirtualMachines virtual_machines = 1;
}
```

С полной структурой данных вы можете ознакомиться в [протофайлах](../../../.proto/service_hypervisor_manager.proto).

## Пример реализации

Реализация операции будет производиться на ESXI версии 7.???????? запущеной на оборудовании `Gagar>n oracul gen. 1`.

Дополняем уже имеющийся [шаблон](../create_project/project/main.go) релизацией RPC `CollectVirtialMachinesList`.

В ходе обработки запроса на языке golang для гипервизора `ESXI` мы воспользуемся библиотекой [govmomi](https://github.com/vmware/govmomi) для взаимодействия с VMware vSphere APIs и сбора информации с гипервизора.

> КАК Я ПОНЯЛ ОТКУДА МНЕ ПОЛУЧИТЬ ДАННЫЕ (ссылка на пример, че нибудь такое)

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

Пример готового проекта расположен в папке [project](./project)
