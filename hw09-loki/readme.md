# ДЗ-09: Vector

Перенаправил Vector из конфигурации 8-го ДЗ на Loki.
 
## Установка

https://grafana.com/docs/loki/latest/setup/install/

## Конфиги

**loki** – стандартный конфиг для хранения логов на файловой системе.

```
auth_enabled: false
 
server:
  http_listen_port: 3100
  grpc_listen_port: 9096
 
common:
  instance_addr: 127.0.0.1
  path_prefix: /tmp/loki
  storage:
    filesystem:
      chunks_directory: /tmp/loki/chunks
      rules_directory: /tmp/loki/rules
  replication_factor: 1
  ring:
    kvstore:
      store: inmemory
 
query_range:
  results_cache:
    cache:
      embedded_cache:
        enabled: true
        max_size_mb: 100
 
schema_config:
  configs:
    - from: 2020-10-24
      store: tsdb
      object_store: filesystem
      schema: v13
      index:
        prefix: index_
        period: 24h
 
analytics:
  reporting_enabled: false
```

**vector.yaml**

```
 sources:
  sshd_http_source:
    type: "socket"
    address: "127.0.0.1:10514"
    mode: "udp"
    decoding:
      codec: "json"
  sshd_log_source:
    type: "file"
    read_from: "beginning"
    include:
      - "/var/log/sshd.log"
 
 
transforms:
  sshd_http_transform:
    type: "remap"
    inputs: ["sshd_http_source"]
    source: '.programname=upcase!(.programname) + "-from-vector-http"'
  sshd_log_transform:
    type: "remap"
    inputs: ["sshd_log_source"]
    source: |
     . |= parse_syslog!(.message)
     .@timestamp = .timestamp
     .programname = "SSHD-from-vector-log-file"
 
sinks:
  print:
    type: "console"
    inputs: ["sshd_http_transform", "sshd_log_transform"]
    encoding:
      codec: "json"
  loki:
    type: "loki"
    inputs: ["sshd_*_transform"]
    endpoint: http://localhost:3100
    encoding:
      codec: json
    labels:
      source: vector
      programname: "{{ programname }}"
      severity: "{{ severity }}"
    out_of_order_action: drop
    path: /loki/api/v1/push
```

## Получение логов

Проверка, что метки появились:

![labels](img/labels.png?raw=true "Labels" )

Получение логов:

![losg](img/logs.png?raw=true "logs" )
