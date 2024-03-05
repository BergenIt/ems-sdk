# BMC manager

## Общее описание модуля расширения

BMC manager - стандартный модуль расширения, предназначеный для обеспечения функционала мониторинга и управления оборудования (в том числе и на BMC устройства).

1) Мониторинг и сбор данных включает в себя:

   * Сбор справочной информации устройства (модель, вендор, серийный номер, ОС, сетевые интерфейсы);
   * Сбор информации о дисках устройства;
   * Сбор статистической информации о оперативной памяти устройства;
   * Сбор инвентарной информации о оперативной памяти устройства;
   * Сбор статистической информации о процессорах устройства;
   * Сбор инвентарной информации о процессорах устройства;
   * Сбор информации и оповещение о событиях устройства;
   * Сбор инвентарной информации устройства;
   * Сбор статуса доступности устройства;
   * Сбор статуса питания устройства;
   * Сбор значения энергопотребления устройства;
   * Сбор значений сенсоров устройства;
   * Сбор информации о версиях прошивок bmc/bios;
   * Сбор текущего статуса BMC устройства;
   * Сбор текущего перечня PCI устройств оборудования;
   * Сбор значения температуры всех устройств;
   * Получение sel логов (system event log).

2) Управление включает в себя:

   * Включение оборудования;
   * Soft-включение оборудования;
   * Выключение оборудования;
   * Soft-выключение оборудования;
   * Перезагрузка оборудования;
   * Soft-перезагрузка оборудования;
   * Эмуляция нажатия кнопки (PowerPushButton);
   * Апаратное прерывание (nmi);
   * Перезагрузка BMC;
   * Soft-перезагрузка BMC;
   * Смена режима загрузки BIOS/UEFI;
   * Включение Led подсветки оборудования;
   * Выключение Led подсветки оборудования;
   * Смена загрузочного носителя;
   * Обновление прошивок  `BMC` и `BIOS/UEFI`;
   * Установка максимального энергопотребления оборудования.

## Описание структуры модуля расширения

Где в какой что, по сути тут мы описывает верхнеурово сценарий работы с этим типом модуля расширения

## Операции модуля расширения

Операция 1 + ссылка на ее доку (если описана)


> Сделал проект

> Реализовал по нашим инструкциям нужные n-функций модуля расширения

> Сдеплоил его