# Autoscaling with Prometheus Adapter

This sample application provides some basic mechanism to play around with Horizontal Pod Autoscaler,
including Prometheus Adapter which allows to scale based on custom metrics.

## Prerequisites

* Kubernetes Cluster, e.g. minikube
* Docker Registry

### Minikube

The following minikube addons (`minikube addons list`) have been enabled: `ingress`, `ingress-dns`, `metrics-server`.
To make the ingress configuration work the resolve service must be configured:

```
sudo mkdir /etc/systemd/resolved.conf.d
sudo tee /etc/systemd/resolved.conf.d/minikube.conf << EOF
[Resolve]
DNS=$(minikube ip)
Domains=~test
EOF

sudo systemctl restart systemd-resolved.service
```

#### Some commands


## Deploy Prometheus
```
kubectl create namespace prometheus
kubectl apply -f deployments/prometheus -n prometheus
```

Alternative access to a service when ingress is not available:

```
kubectl expose deployment prometheus-deployment --type=NodePort --port=9090
minikube service prometheus-deployment -n prometheus --url
```

## Build and deploy sample application

```
kubectl create namespace autoscaling
kubectl config set-context --current --namespace autoscaling
eval $(minikube docker-env)
go build && docker build -t autoscaling:latest .
kubectl create -f deployments/autoscaling/
# or the reload the image after a build
kubectl rollout restart deployment.apps/autoscaling
```

### Security

The `cpu` endpoint can be protected with basic auth which is simply enabled with two environment variables:

```
export USER=username
export PASSWORD=your secret
```

A client must call the protected endpoint with an `Authorization` header, e.g.:

```
curl -H "Authorization: Basic $(echo -n $USER:$PASSWORD | base64)" http://localhost:3333/cpu
```

## Deploy Prometheus Adapter
```
kubectl create namespace monitoring
kubectl config set-context --current --namespace monitoring
kubectl apply -f deployments/prometheus-adapter
kubectl get apiservices
kubectl get --raw /apis/custom.metrics.k8s.io/v1beta2 | jq
```

## HPA with CPU metric

```
kubectl autoscale deployment autoscaling --cpu-percent=50 --min=1 --max=10
```

> Remove other HPA deployment, e.g. the `autoscaling-hpa` deployed with the sample application!

Execute `curl http://autoscaling.test/cpu` to generate some cpu load. 
> Take care with this command, with each request one cpu is used by 50% for 10 minutes!

### Prometheus-Query

Open the Prometheus UI in  a Browser: http://localhost:9090

```
sum(rate(container_cpu_usage_seconds_total{pod=~"autoscaling.*"}[1m])) 
/ 
sum(kube_pod_container_resource_requests{pod=~"autoscaling.*",resource="cpu"}) * 100
```

## HPA with custom metric

```
kubectl create -f deployments/autoscaling/hpa.yaml
```

Execute `curl http://autoscaling.test/add` to increase the metric and `curl http://autoscaling.test/remove` to decrease the metric.

## References

* https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/
* https://github.com/kubernetes-sigs/prometheus-adapter/blob/master/docs/walkthrough.md
* https://github.com/kubernetes-sigs/prometheus-adapter/blob/master/docs/config-walkthrough.md