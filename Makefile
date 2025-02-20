CLUSTER_NAME=canary-demo
NAMESPACE=argo-rollouts

.PHONY: cluster-up
cluster-up:
	k3d cluster create $(CLUSTER_NAME) --config ./k3d/k3d-config.yaml
	kubectl cluster-info --context k3d-$(CLUSTER_NAME)

.PHONY: cluster-down
cluster-down:
	k3d cluster delete $(CLUSTER_NAME)

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
