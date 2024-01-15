**ДЗ-01: Prometheus, exporters** 



**Установка/настройка БД**

apt install mariadb-server

apt install mariadb-client

mysql\_secure\_installation

mysql -u root -p

CREATE DATABASE joomladb;

CREATE USER 'joomladbuser'@'localhost' IDENTIFIED BY 'pas';

GRANT ALL ON joomladb.\* TO 'joomladbuser'@'localhost' WITH GRANT OPTION;

Пользователь для exporter

CREATE USER 'exporter'@'localhost' IDENTIFIED BY 'pas' WITH MAX\_USER\_CONNECTIONS 3;

GRANT PROCESS, REPLICATION CLIENT, SELECT ON \*.\* TO 'exporter'@'localhost';

FLUSH PRIVILEGES;

EXIT;

**Установка PHP, PHP-FPM**

apt install php

apt install php-fpm php-mysql

apt install php8.2-zip php8.2-intl php8.2-xml php8.2-curl php8.2-mbstring php8.2-gd

**Установка Joomla**

unzip -d /var/www/joomla /tmp/Joomla\_5.0.2-Stable-Full\_Package.zip

chown -R www-data:www-data /var/www/joomla/

**Установка NGINX, настройка для Joomla** 

apt update

apt install nginx

Конфиг:

/etc/nginx/sites-available/joomla.conf

server {

`    `listen 81;

`    `listen [::]:81;

`    `root /var/www/joomla;

`    `index  index.php index.html index.htm;

`    `server\_name  example.com www.example.com;



`    `client\_max\_body\_size 100M;

`    `autoindex off;



`    `location / {

`        `try\_files $uri $uri/ /index.php?$args;

`    `}



`    `# deny running scripts inside writable directories

`    `location ~\* /(images|cache|media|logs|tmp)/.\*.(php|pl|py|jsp|asp|sh|cgi)$ {

`      `return 403;

`      `error\_page 403 /403\_error.html;

`    `}



`    `location ~ .php$ {

`        `include snippets/fastcgi-php.conf;

`        `fastcgi\_pass unix:/run/php/php8.0-fpm.sock;

`        `fastcgi\_param SCRIPT\_FILENAME $document\_root$fastcgi\_script\_name;

`        `include fastcgi\_params;

`    `}

}

ln -s /etc/nginx/sites-available/joomla.conf /etc/nginx/sites-enabled/

systemctl restart nginx.service

![](Aspose.Words.1660484c-b1c5-4e33-96e5-35877b6eb110.001.png)

**

**Установка/настройка exporters**

apt-get install prometheus-node-exporter

http://localhost:9100/metrics

![](Aspose.Words.1660484c-b1c5-4e33-96e5-35877b6eb110.002.png)


MySql exporter

<https://www.devopsschool.com/blog/install-and-configure-prometheus-mysql-exporter/>

apt-get install prometheus-mysqld-exporter

/etc/.mysqld\_exporter.cnf

[client]

user=exporter

password=pas

chown root:prometheus /etc/.mysqld\_exporter.cnf

Добавляем параметр для конфига при запуске сервиса  /lib/systemd/system/prometheus-mysqld-exporter.service

--config.my-cnf /etc/.mysqld\_exporter.cnf

systemctl daemon-reload

systemctl enable prometheus-mysqld-exporter

systemctl start prometheus-mysqld-exporter

![](Aspose.Words.1660484c-b1c5-4e33-96e5-35877b6eb110.003.png)



Blackbox exporter

<https://kirelos.com/how-to-monitor-website-performance-with-blackbox-exporter-and-grafana/>

Скачиваем, делаем кофиг для сервиса

wget https://github.com/prometheus/blackbox\_exporter/releases/download/v0.24.0/blackbox\_exporter-0.24.0.linux-amd64.tar.gz

tar -xzf blackbox\_exporter-0.24.0.linux-amd64.tar.gz

[Unit]

Description=Blackbox Exporter Service

Wants=network-online.target

After=network-online.target



[Service]

Type=simple

User=root

Group=root

ExecStart=/etc/blackbox\_eporter/blackbox\_exporter --config.file==/etc/blackbox\_eporter/blackbox.yml



[Install]

WantedBy=multi-user.target

systemctl daemon-reload

systemctl enable prometheus-blackbox-exporter

systemctl start prometheus-blackbox-exporter


![](Aspose.Words.1660484c-b1c5-4e33-96e5-35877b6eb110.004.png)

**Один порт для метрик и аутентификация**

Для использования одного порта добавил в NGIX настройки для перенаправления запросов

location /metrics/db/ {

`  `proxy\_pass http://localhost:9104/metrics;

}

location /metrics/node/ {

`  `proxy\_pass http://localhost:9100/metrics;

}

location /metrics/blackbox/ {

`  `proxy\_pass http://localhost:9115/;

}

И добавил базовую атентифкацию

htpasswd -c /etc/nginx/conf.d/.htpasswd prom

server{

...

`    `auth\_basic "Restricted Access!";

`    `auth\_basic\_user\_file /etc/nginx/conf.d/.htpasswd;

...

}


![](Aspose.Words.1660484c-b1c5-4e33-96e5-35877b6eb110.005.png)

#### **Настройка Prometheus**
<https://computingforgeeks.com/how-to-install-prometheus-and-node-exporter-on-debian/>

groupadd --system prometheus

useradd -s /sbin/nologin --system -g prometheus prometheus

mkdir /var/lib/prometheus

for i in rules rules.d files\_sd; do sudo mkdir -p /etc/prometheus/${i}; done

mkdir -p /tmp/prometheus &amp;&amp; cd /tmp/prometheus

curl -s https://api.github.com/repos/prometheus/prometheus/releases/latest|grep browser\_download\_url|grep linux-amd64|cut -d '"' -f 4|wget -qi -

tar xvf prometheus\*.tar.gz

cd prometheus\*/

mv prometheus promtool /usr/local/bin/

mv prometheus.yml  /etc/prometheus/prometheus.yml

mv consoles/ console\_libraries/ /etc/prometheus/

cd ~/

rm -rf /tmp/prometheus

Конфиг для сервиса

tee /etc/systemd/system/prometheus.service<<EOF

[Unit]

Description=Prometheus

Documentation=https://prometheus.io/docs/introduction/overview/

Wants=network-online.target

After=network-online.target



[Service]

Type=simple

User=prometheus

Group=prometheus

ExecReload=/bin/kill -HUP $MAINPID

ExecStart=/usr/local/bin/prometheus \

`  `--config.file=/etc/prometheus/prometheus.yml \

`  `--storage.tsdb.path=/var/lib/prometheus \

`  `--web.console.templates=/etc/prometheus/consoles \

`  `--web.console.libraries=/etc/prometheus/console\_libraries \

`  `--web.listen-address=0.0.0.0:9090 \

`  `--web.external-url=



SyslogIdentifier=prometheus

Restart=always



[Install]

WantedBy=multi-user.target

EOF

for i in rules rules.d files\_sd; do sudo chown -R prometheus:prometheus /etc/prometheus/${i}; done

for i in rules rules.d files\_sd; do sudo chmod -R 775 /etc/prometheus/${i}; done

chown -R prometheus:prometheus /var/lib/prometheus/

Reload systemd daemon and start the service.

sudo systemctl daemon-reload

sudo systemctl start prometheus

sudo systemctl enable prometheus

![](Aspose.Words.1660484c-b1c5-4e33-96e5-35877b6eb110.006.png)

Конфиг 



global:

`  `scrape\_interval: 15s

`  `scrape\_timeout: 10s

`  `scrape\_protocols:

`  `- OpenMetricsText1.0.0

`  `- OpenMetricsText0.0.1

`  `- PrometheusText0.0.4

`  `evaluation\_interval: 15s

alerting:

`  `alertmanagers:

`  `- follow\_redirects: true

`    `enable\_http2: true

`    `scheme: http

`    `timeout: 10s

`    `api\_version: v2

`    `static\_configs:

`    `- targets: []

scrape\_configs:

\- job\_name: prometheus

`  `honor\_timestamps: true

`  `track\_timestamps\_staleness: false

`  `scrape\_interval: 15s

`  `scrape\_timeout: 10s

`  `scrape\_protocols:

`  `- OpenMetricsText1.0.0

`  `- OpenMetricsText0.0.1

`  `- PrometheusText0.0.4

`  `metrics\_path: /metrics

`  `scheme: http

`  `enable\_compression: true

`  `follow\_redirects: true

`  `enable\_http2: true

`  `static\_configs:

`  `- targets:

`    `- localhost:9090

\- job\_name: db

`  `honor\_timestamps: true

`  `track\_timestamps\_staleness: false

`  `scrape\_interval: 15s

`  `scrape\_timeout: 10s

`  `scrape\_protocols:

`  `- OpenMetricsText1.0.0

`  `- OpenMetricsText0.0.1

`  `- PrometheusText0.0.4

`  `metrics\_path: /metrics/db/

`  `scheme: http

`  `enable\_compression: true

`  `basic\_auth:

`    `username: prom

`    `password: <secret>

`  `follow\_redirects: true

`  `enable\_http2: true

`  `static\_configs:

`  `- targets:

`    `- 192.168.0.125:82

\- job\_name: node

`  `honor\_timestamps: true

`  `track\_timestamps\_staleness: false

`  `scrape\_interval: 15s

`  `scrape\_timeout: 10s

`  `scrape\_protocols:

`  `- OpenMetricsText1.0.0

`  `- OpenMetricsText0.0.1

`  `- PrometheusText0.0.4

`  `metrics\_path: /metrics/node/

`  `scheme: http

`  `enable\_compression: true

`  `basic\_auth:

`    `username: prom

`    `password: <secret>

`  `follow\_redirects: true

`  `enable\_http2: true

`  `static\_configs:

`  `- targets:

`    `- 192.168.0.125:82

\- job\_name: blackbox

`  `honor\_timestamps: true

`  `track\_timestamps\_staleness: false

`  `params:

`    `module:

`    `- http\_2xx

`  `scrape\_interval: 15s

`  `scrape\_timeout: 10s

`  `scrape\_protocols:

`  `- OpenMetricsText1.0.0

`  `- OpenMetricsText0.0.1

`  `- PrometheusText0.0.4

`  `metrics\_path: /metrics/blackbox/probe

`  `scheme: http

`  `enable\_compression: true

`  `basic\_auth:

`    `username: prom

`    `password: <secret>

`  `follow\_redirects: true

`  `enable\_http2: true

`  `relabel\_configs:

`  `- source\_labels: [\_\_address\_\_]

`    `separator: ;

`    `regex: (.\*)

`    `target\_label: \_\_param\_target

`    `replacement: $1

`    `action: replace

`  `- source\_labels: [\_\_param\_target]

`    `separator: ;

`    `regex: (.\*)

`    `target\_label: instance

`    `replacement: $1

`    `action: replace

`  `- separator: ;

`    `regex: (.\*)

`    `target\_label: \_\_address\_\_

`    `replacement: 192.168.0.125:82

`    `action: replace

`  `static\_configs:

`  `- targets:

`    `- http://192.168.0.125:81

\- job\_name: blackbox\_exporter

`  `honor\_timestamps: true

`  `track\_timestamps\_staleness: false

`  `scrape\_interval: 15s

`  `scrape\_timeout: 10s

`  `scrape\_protocols:

`  `- OpenMetricsText1.0.0

`  `- OpenMetricsText0.0.1

`  `- PrometheusText0.0.4

`  `metrics\_path: /metrics/blackbox/metrics

`  `scheme: http

`  `enable\_compression: true

`  `basic\_auth:

`    `username: prom

`    `password: <secret>

`  `follow\_redirects: true

`  `enable\_http2: true

`  `static\_configs:

`  `- targets:

`    `- 192.168.0.125:82



