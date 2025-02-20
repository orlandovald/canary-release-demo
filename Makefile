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
	@echo "Installing Istio! This might take a while..."
	kubectl get ns $(ISTIO_NAMESPACE) > /dev/null 2>&1 || kubectl create ns $(ISTIO_NAMESPACE)
	helm repo add istio https://istio-release.storage.googleapis.com/charts && helm repo update
	helm upgrade -i istio-base istio/base -n $(ISTIO_NAMESPACE) --wait
	helm upgrade -i istiod istio/istiod -n $(ISTIO_NAMESPACE) --wait
	helm upgrade -i istio-ingressgateway istio/gateway -n $(ISTIO_NAMESPACE) --wait
	kubectl label namespace default istio-injection=enabled

.PHONY: install-argo
install-argo:
	kubectl get ns $(NAMESPACE) > /dev/null 2>&1 || kubectl create ns $(NAMESPACE)
	kubectl apply -n $(NAMESPACE) -f https://github.com/argoproj/argo-rollouts/releases/download/v1.8.0/install.yaml

.PHONY: build-demo-app
build-demo-app:
	cd demo-app && docker buildx build -t canary-demo:latest .

.PHONY: cluster-import-image
cluster-import-image:
	k3d images import canary-demo:latest --cluster $(CLUSTER_NAME)
