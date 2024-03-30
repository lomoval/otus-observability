# ДЗ-07: ELK-beats
 
## Установка

https://www.elastic.co/guide/en/beats/heartbeat/current/heartbeat-installation-configuration.html
https://www.elastic.co/guide/en/beats/filebeat/current/filebeat-installation-configuration.html
https://www.elastic.co/guide/en/beats/metricbeat/current/metricbeat-installation-configuration.html

## Конфиги

**/etc/heartbeat/heartbeat.yml**

```
heartbeat.monitors:
- type: http
  enabled: true
  id: otus-hw-monitor
  name: Otus HW Monitor
  urls: ["https://otus.ru", "https://google.com"]
  schedule: '@every 10s'

output.elasticsearch:
  hosts: ["localhost:9200"]
  preset: balanced
  protocol: "https"
  username: "elastic"
  password: "10DM27XoI52iiDxXZa78"
  ssl.verification_mode: none
```

**/etc/metricbeat/metricbeat.yml**

```
output.elasticsearch:
  hosts: ["localhost:9200"]
  preset: balanced
  protocol: "https"
  username: "elastic"
  password: "10DM27XoI52iiDxXZa78"
  ssl.verification_mode: none
```

Подключаем модуль

```
metricbeat modules enable system
```

**/etc/metricbeat/modules.d/system.yml**

```
- module: system
  period: 10s
  metricsets:
    - cpu
    - memory
  process.include_top_n:
    by_cpu: 5      # include top 5 processes by CPU
    by_memory: 5   # include top 5 processes by memory

/etc/filebeat/metricbeat.yml
output.elasticsearch:
  hosts: ["localhost:9200"]
  preset: balanced
  protocol: "https"
  username: "elastic"
  password: "10DM27XoI52iiDxXZa78"
  ssl.verification_mode: none
```

Подключаем модуль

```
filebeat modules enable system
```

**/etc/filebeat/modules.d/system.yml**

```
- module: system
  syslog:
    enabled: true
  auth:
    enabled: true
```

## Индексы

![index](img/indexes.png?raw=true "Index" )

## Kibana

![heartbeat](img/heartbeat.png?raw=true "Heartbeat" )

![metricbeat](img/metricbeat.png?raw=true "Metricbeat" )

![filebeat](img/filebeat.png?raw=true "Filebeat" )
