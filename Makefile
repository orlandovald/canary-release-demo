CLUSTER_NAME=canary-demo
NAMESPACE=argo-rollouts
ISTIO_NAMESPACE=istio-system

.PHONY: cluster-up
cluster-up:
	k3d cluster create $(CLUSTER_NAME) --config ./k3d/k3d-config.yaml
	kubectl cluster-info --context k3d-$(CLUSTER_NAME)

.PHONY: cluster-down
cluster-down:
	k3d cluster delete $(CLUSTER_NAME)

.PHONY: install-istio
install-istio:
	@echo "\033[33m=== Installing Istio. This might take a while ===\033[0m"
	kubectl get ns $(ISTIO_NAMESPACE) > /dev/null 2>&1 || kubectl create ns $(ISTIO_NAMESPACE)
	helm repo add istio https://istio-release.storage.googleapis.com/charts && helm repo update
	helm upgrade -i istio-base istio/base -n $(ISTIO_NAMESPACE) --wait
	helm upgrade -i istiod istio/istiod -n $(ISTIO_NAMESPACE) --wait
	helm upgrade -i istio-ingressgateway istio/gateway -n $(ISTIO_NAMESPACE) --wait
	kubectl label namespace default istio-injection=enabled

.PHONY: install-prometheus
install-prometheus:
	@echo "\033[33mInstalling Prometheus ===\033[0m"
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update
	helm upgrade -i -n prometheus --create-namespace prometheus prometheus-community/prometheus --set server.global.scrape_interval=10s

.PHONY: install-argo
install-argo:
	@echo "\033[33m=== Installing Argo Rollouts ===\033[0m"
	kubectl get ns $(NAMESPACE) > /dev/null 2>&1 || kubectl create ns $(NAMESPACE)
	kubectl apply -n $(NAMESPACE) -f https://github.com/argoproj/argo-rollouts/releases/download/v1.8.0/install.yaml

.PHONY: build-demo-app
build-demo-app:
	cd demo-app && docker build . -t canary-demo-app:green --build-arg APP_VERSION=green 
	cd demo-app && docker build . -t canary-demo-app:blue --build-arg APP_VERSION=blue 
	cd demo-app && docker build . -t canary-demo-app:yellow --build-arg APP_VERSION=yellow --build-arg ERROR_RATE=30 

.PHONY: cluster-import-image
cluster-import-image:
	k3d images import canary-demo-app:green --cluster $(CLUSTER_NAME)
	k3d images import canary-demo-app:blue --cluster $(CLUSTER_NAME)
	k3d images import canary-demo-app:yellow --cluster $(CLUSTER_NAME)
