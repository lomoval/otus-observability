# ДЗ-03: Grafana - dashboard

## Установка

[https://grafana.com/docs/grafana/latest/setup-grafana/installation/debian](https://grafana.com/docs/grafana/latest/setup-grafana/installation/debian/)

`/etc/grafana/grafana.ini`
Set port = 3000

## Alerting

Telegram bot не заработал на версии 10.3.1, хотя вроде как поправили (ошибка: https://github.com/grafana/grafana/issues/73068)

По почте на `mail.ru`:

`/etc/grafana/grafana.ini`

```
[smtp]
enabled = true

host = smtp.mail.ru:465
user =tmp2023@list.ru
password =*********
;cert_file =
;key_file =
;skip_verify = false
from_address =tmp2023@list.ru
from_name =tmp2023@list.ru
```
