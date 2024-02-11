# ДЗ-04: Alertmanager 

## Установка alertmanager

`wget https://github.com/prometheus/alertmanager/releases/download/v0.26.0/alertmanager-0.26.0.linux-amd64.tar.gz`

```
tar xvf alertmanager*.tar.gz
cd alertmanager*/
mv alertmanager amtool /usr/local/bin/
mkdir /etc/alertmanager
mkdir /var/lib/prometheus/alertmanager
mv alertmanager.yml /etc/alertmanager/alertmanager.yml
cd ..
rm -rf alertmanager*/
```

```
useradd --no-create-home --shell /bin/false alertmanager;
chown -R alertmanager:alertmanager /etc/alertmanager /var/lib/prometheus/alertmanager
chown alertmanager:alertmanager /usr/local/bin/{alertmanager,amtool}
```

## Настройка сервиса

`/etc/systemd/system/alertmanager.service`

```
[Unit]
Description=Alertmanager Service
After=network.target
 
[Service]
EnvironmentFile=-/etc/default/alertmanager
User=alertmanager
Group=alertmanager
Type=simple
ExecStart=/usr/local/bin/alertmanager \
--config.file=/etc/alertmanager/alertmanager.yml \
--storage.path=/var/lib/prometheus/alertmanager \
--cluster.advertise-address="127.0.0.1:9093"\
$ALERTMANAGER_OPTS
ExecReload=/bin/kill -HUP $MAINPID
Restart=on-failure
 
[Install]
WantedBy=multi-user.target
```

```
systemctl daemon-reload;
systemctl enable alertmanager;
systemctl start alertmanager;
systemctl status alertmanager;
```

## Настройка Prometheus

`/etc/prometheus/prometheus.yml`

```
alerting:
 alertmanagers:
   - static_configs:
       - targets:
          - localhost:9093
rule_files:
  - "rules.yml"
```

`/etc/prometheus/rules.yml`

```
groups:
- name: joomla
  rules:
  - alert: JoomlaInstanceDown
    expr: probe_success == 0
    for: 15s
    annotations:
      title: 'Joomla instance {{ $labels.instance }} down'
      description: 'Joomla {{ $labels.instance }} of job {{ $labels.job }} has been down for more than 1 minute.'
    labels:
      severity: 'critical'
- name: env
  rules:
  - alert: LowMemory
    expr: (100 - ((avg_over_time(node_memory_MemAvailable_bytes[5m]) * 100) / avg_over_time(node_memory_MemTotal_bytes[5m]))) > 60
    for: 15s
    annotations:
      title: 'Low memory {{ $labels.instance }}'
    labels:
      severity: 'warning'
```

`Настройка Alertmanager`

`/etc/alertmanager/alertmanager.yml`
 
```
route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 30s
  receiver: 'default'
  routes:
    - match:
        severity: critical
      receiver: 'telegram'
    - match:
        severity: warning
      receiver: 'web.hook'
receivers:
  - name: 'default'
    webhook_configs:
    - url: 'http://localhost:1010'
  - name: 'web.hook'
    webhook_configs:
    - url: 'https://api.telegram.org/bot***/sendMessage?chat_id=-***&amp;text=LowMemory'
  - name: 'telegram'
    telegram_configs:
    - bot_token: "***"
      api_url: "https://api.telegram.org"
      chat_id: -1002062818733
      parse_mode: ""
      send_resolved: true
      message: "🚨 Alertmanager 🚨\nPANIC\n🔺 Alertname: {{ .GroupLabels.alertname}}\n🔺 Severity: {{ .CommonLabels.severity }}\n📌 {{ range .Alerts }}{{ .Annotations.summary }}\n{{ end }}
```
![alerts](img/alerts.jpg?raw=true "Alerts" )