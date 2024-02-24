# ДЗ-04: Alertmanager 

## Установка Zabbix

`https://www.zabbix.com/download?zabbix=6.4&os_distribution=debian&os_version=12&components=server_frontend_agent&db=mysql&ws=apache`

## Настройка шаблона (LLD)

Создание шаблона

![template](img/setup1.jpg?raw=true "Template" )

Правило обнаружения

![template](img/setup2.jpg?raw=true "Template" )

![template](img/setup3.jpg?raw=true "Template" )

Препроцесинг при помощи JS:

```
const regex = /(\w+)\[(\w+)\] (\d+)/g;
var match;
var lines = [];
while ((match = regex.exec(value)) !== null) {
    const name = match[2];
    const val = match[3];
    lines.push({ name, value: Number(val) });
}
return JSON.stringify(lines);
```

![template](img/setup4.jpg?raw=true "Template" )

Прототип  итема

![template](img/setup5.jpg?raw=true "Template" )

![template](img/setup6.jpg?raw=true "Template" )

Тайкой же JS скрипт как выше + берем значение нужной метрики.

Триггер

![trigger](img/setup7.jpg?raw=true "Trigger" )

Алерт

У пользователя должны быть права на чтение host’a.

![alert](img/setup8.jpg?raw=true "Alert" )

![alert](img/setup9.jpg?raw=true "Alert" )

![alert](img/setup10.jpg?raw=true "Alert" )

![alert](img/setup11.jpg?raw=true "Alert" )

Хост

![host](img/setup12.jpg?raw=true "Host" )

## Dashboard

![dashboard](img/dashboard.jpg?raw=true "Dashboard" )

## Telegram

![telega](img/telega.jpg?raw=true "Telega" )