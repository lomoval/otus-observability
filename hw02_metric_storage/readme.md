# ДЗ-02: Хранилище метрик – VictoriaMetrics

## Установка

[https://www.dmosk.ru/miniinstruktions.php?mini=victoriametrics](https://www.dmosk.ru/miniinstruktions.php?mini=victoriametrics)
```
wget https://github.com/VictoriaMetrics/VictoriaMetrics/releases/download/v1.96.0/victoria-metrics-linux-amd64-v1.96.0.tar.gz
tar zxf victoria-metrics-linux-amd64-*.tar.gz -C /usr/local/bin/
useradd -r -c 'VictoriaMetrics TSDB Service' victoriametrics
mkdir -p /var/lib/victoriametrics /run/victoriametrics 
chown victoriametrics:victoriametrics /var/lib/victoriametrics /run/victoriametrics
```

## Настройка сервиса

`vi /etc/systemd/system/victoriametrics.service`
```
[Unit]
Description=VictoriaMetrics
After=network.target
[Service]
Type=simple
User=victoriametrics
PIDFile=/run/victoriametrics/victoriametrics.pid
ExecStart=/usr/local/bin/victoria-metrics-prod -storageDataPath /var/lib/victoriametrics -retentionPeriod 2w
ExecStop=/bin/kill -s SIGTERM $MAINPID
StartLimitBurst=5
StartLimitInterval=0
Restart=on-failure
RestartSec=1
[Install]
WantedBy=multi-user.target
```
```
systemctl daemon-reload
systemctl enable victoriametrics --now
systemctl status victoriametrics
curl 127.0.0.1:8428
```

##  Настройка prometheus

`vi /etc/prometheus/prometheus.yml`

```
global:
  external_labels:
    site: prod
remote_write:
  - url: http://127.0.0.1:8428/api/v1/write
    queue_config:
      max_samples_per_send: 10000
      capacity: 20000
      max_shards: 30
```

```
systemctl restart prometheus
http://127.0.0.1:8428/vmui/
```




