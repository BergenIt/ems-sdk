# EMS SDK

## Описание решения

Gagarin EMS - это система управления, инвентаризации и мониторинга ИТ-инфраструктуры корпоративного класса с возможностью расширения и адаптации функциональных возможностей под потребности конечного пользователя.  
Относится к классу систем управления и контроля групп однотипных сетевых или вычислительных элементов.

Предназначена для решения следующих задач:

  - Автоматизация задач управления серверным оборудованием;
  - Автоматизация задач управления сетевым оборудованием;
  - Мониторинг доступности серверов, устройств хранения и сетевых коммутаторов;
  - Мониторинг виртуальных машин, гипервизоров, систем и сервисов;
  - Мониторинг компонентов серверов;
  - Автоматизация развёртывания операционных систем и программного обеспечения;
  - Предоставление информации для планирования модернизации оборудования ИТ-инфраструктуры;
  - Администрирование географически распределенной ИТ-инфраструктуры.

Является системой управления и мониторинга для работы в мульти-вендорных окружениях. С помощью открытых протоколов, система EMS может взаимодействовать с широким спектром оборудования, позволяя администраторам собирать и отображать инвентарную информацию, информацию по здоровью и загруженности объектов ИТ-инфраструктуры. Возможность расширения функциональных возможностей системы в рамках поддерживаемого перечня метрик и операций управления конечным пользователем, средствами SDK. Помимо сбора и отображения информации об оборудовании, система EMS позволяет производить операции управления над серверами, включая управление питанием, настройками а позволяя администраторам устанавливать ОС и прикладное ПО, работая одновременно с любым числом оборудования. Гранулярное управления правами и доступом к оборудованию, детальное журналирование действий администраторов, функция двойного подтверждения операций над серверами а также возможность использования технологии единого входа на серверы, делает управление серверами при помощи EMS безопасным и защищённым от несанкционированного использования. 

Поставляемый SDK включает в себя набор инструментов, посредством которых разработчик может адаптировать EMS под текущие потребности и расширять перечень производителей оборудования для обеспечения поддержки полного фукнционала системы самостоятельно добавляя необходимые драйверы, скрипты и настройки. Таким образом, вы сможете подключать к EMS SDK любое оборудование, которое поддерживает стандартные протоколы, такие как Redfish, SNMP, IPMI, SSH и т.д. К инструментам для разработки програмного относятся:

  - Примеры кода (С#, Golang);
  - Документация для разработчика модулей расширения.

## Стандартные модули расширение системы 

Описание и документация стандартных модулей расширения доступны по ссылкам: 

  - [BMC manager](extend_bmc/README.md);
  - [Hypervisor manager](extend_hypervisor/README.md);
  - [Linux manager](extend_linux/README.md);
  - [Network/switch manager](extend_network_switch/README.md);
  - [OS installer](extend_os_install/README.md);
  - [SNMP/Syslog manager](extend_snmp/README.md);
  - [SSO center](extend_sso_bmc/README.md);
  - [Web-service manager](extend_web_service/README.md);
  - [Windows manager](extend_windows/README.md);

## UJM - killme

* Создание модуля расширения для работы с BMC
  * Создание проекта
  * Реализация операции сбора LED статуса
  * Реализация операции прошивки BMC
  * Развертывание модуля расширения

* Создание модуля расширения для работы с Linux
  * Создание проекта
  * Реализация операции сбора серийного номера
  * Развертывание модуля расширения

* Создание модуля расширения для работы с Windows
  * Создание проекта
  * Реализация операции сбора серийного номера
  * Развертывание модуля расширения

* Создание модуля расширения для работы с SNMP-шаблонами
  * Создание проекта
  * Реализация операции сбора серийного номера
  * Развертывание модуля расширения

* Создание модуля расширения для работы с SSO BMC
  * Создание проекта
  * Реализация операции подключения BMC в SSO
  * Развертывание модуля расширения

* Создание модуля расширения для работы с гипервизором
  * Создание проекта
  * Реализация операции сбора списка виртуальных машин
  * Развертывание модуля расширения

* Создание модуля расширения для работы с веб-сервисом
  * Создание проекта
  * Реализация операции проверки статуса доступности сервиса (на примере postgres)
  * Развертывание модуля расширения

* Создание модуля расширения для работы с подсетями и коммутаторами
  * Создание проекта
  * Реализация операции создания резервной копии конфигурации коммутатора
  * Развертывание модуля расширения

> TOBE - 4.1.0

* Траблшутинг (назвать нормально + в мин. виде заложить в текущие общие ридми)
* Создание модуля расширения для работы с установки ОС
