# Autoscaling with Prometheus Adapter

## Prerequisites

* Kubernetes Cluster, e.g. minikube
* Docker Registry


## Build and Deploy

```
kubectl create namespace autoscaling
kubectl config set-context --current --namespace autoscaling
kubectl create -f deployments/autoscaling/
```
Prometheus-Query

```
sum(rate(container_cpu_usage_seconds_total{pod=~"sample-app.*"}[1m])) 
/ 
sum(kube_pod_container_resource_requests{pod=~"sample-app.*",resource="cpu"}) * 100
```
https://github.com/kubernetes-sigs/prometheus-adapter/blob/master/docs/walkthrough.md
https://github.com/kubernetes-sigs/prometheus-adapter/blob/master/docs/config-walkthrough.md

```
kubectl get --raw /apis/custom.metrics.k8s.io/v1beta2
k expose deployment prometheus-deployment --type=NodePort --port=9090
minikube service prometheus-deployment -n prometheus --url
k get apiservices
```