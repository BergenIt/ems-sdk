# Network/switch manager

## Общее описание модуля расширения

Network/switch manager - модуль расширения, предназначеный для обеспечения функционала мониторинга оборудования по протоколу `ICMP` и сбора информации о хостовых устройствах, а также создания и развертывания резервных копий конфигураций коммутаторв.  

К обеспечению функционала относится: 

* Сбор ICMP-статуса доступности сетевых интерфейсов и масок подсетей;
* Создание резервной копии конфигурации коммутатора;
* Восстановление конфигруации коммутатора из резервной копии.

## Разработка собственного network модуля расширения

* [Создание проекта](./golang/create_project/README.md)
* [Реализация операции 'Cоздание резервной конфигурации'](./golang/backup_switch/README.md)
* [Развертывание модуля расширения](./golang/deploy/README.md)
