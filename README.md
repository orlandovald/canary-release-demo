# Canary Release Demo

This is a companion repo for my talk about Canary Releases with Argo Rollouts

## Pre-requisites

- Make
- [Docker](https://docs.docker.com/get-docker/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/)
- [K3d](https://k3d.io/stable/#releases)
- [Helm](https://helm.sh/docs/intro/install/)
- [Argo Rollouts Kubectl plugin](https://argoproj.github.io/argo-rollouts/installation/) (optional)

## Set up

First step is to run a local cluster with k3d and install the required
components. We can do that using the below Make targets.

### Create and prepare cluster

```sh
make cluster-up # creates a local K3d cluster
make install-argo
make install-istio # required for traffic routing
make install-prometheus # required to run analyses
```

### Set up demo app

```sh
make build-demo-app # build green, blue and yellow images
make cluster-import-image # load the images into the local cluster
```

### Clean up

Run the below command to destroy all created resources

```sh
make cluster-down
```

## Running the demo

To run the demo, apply to the cluster any of the below Kubernetes manifess
in the `default` namespace

- 101_canary-simple.yaml | Simple rollout
- 203_canary-traffic.yaml | Rollout with traffic routing. Requires filed 201 and 202
- 302_canary-analysis.yaml | Rollout with traffic routing and canary analyses. Required 201, 202 and 301
