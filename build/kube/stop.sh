alias k=kubectl

#clusterrole
k delete clusterrole fluentd

#clusterrolebinding
k delete clusterrolebinding fluentd

#namespace
k delete ns stories
