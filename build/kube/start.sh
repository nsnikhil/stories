alias k=kubectl

#namespace
k create -f common/namespace.yaml

#influxdb
k create -f monitoring/influxdb.yaml

#telegraf
k create -f monitoring/telegraf.yaml

#grafana
k create -f monitoring/grafana.yaml

#elasticsearch
k create -f monitoring/elasticsearch.yaml

#kibana
k create -f monitoring/kibana.yaml

#fluentd
k create -f monitoring/fluentd.yaml

#app
k create -f app/app.yaml
