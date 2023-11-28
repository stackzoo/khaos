# Khaos
[![releaser](https://github.com/stackzoo/khaos/actions/workflows/release.yaml/badge.svg)](https://github.com/stackzoo/khaos/actions/workflows/release.yaml)  [![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)  
<br/>
<img src="docs/images/klogo.png" alt="logo" width="230" height="230">  
<br/>
A lightweight kubernetes operator to test cluster resilience via chaos engineering 💥🌪 ☸️  

## Abstract
**Khaos** is a streamlined Kubernetes [operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) made with [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) and designed for executing [Chaos Engineering](https://en.wikipedia.org/wiki/Chaos_engineering) activities.  
Through the implementation of custom controllers and resources, Khaos facilitates the configuration and automation  
of operations such as the targeted deletion of pods within a specified namespace, the removal of nodes from the cluster, the deletion of secrets and more.  

## Supported features
- [X] Delete specified pods in specified namespace
- [x] Delete specified cluster nodes
- [X] Delete specified secrets in specified namespace  
- [X] Inject resource constraints in the specified containers  of the specified deployment of the specified namespace
- [X] Inject the specified command inside the pods of the specified deployment in the specified namespace (experimental).  


## Local Testing
First of all clone the repository:  
```console
git clone https://github.com/stackzoo/khaos && cd khaos
```  

The repo contain a Makefile with all that you need.  
Inspect the make targets with the following command:  
```console
make help

Usage:
  make <target>

General
  help             Display this help.

Development
  manifests        Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
  generate         Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
  fmt              Run go fmt against code.
  vet              Run go vet against code.
  cluster-up       Create a kind cluster named "test-operator-cluster" with a master and 3 worker nodes.
  cluster-down     Delete the kind cluster named "test-operator-cluster".
  test             Run tests.
  lint             Run golangci-lint linter & yamllint
  lint-fix         Run golangci-lint linter and perform fixes

Build
  build            Build manager binary.
  run              Run a controller from your host.
  docker-build     Build docker image with the manager.
  docker-push      Push docker image with the manager.
  docker-buildx    Build and push docker image for the manager for cross-platform support

Deployment
  install          Install CRDs into the K8s cluster specified in ~/.kube/config.
  uninstall        Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
  deploy           Deploy controller to the K8s cluster specified in ~/.kube/config.
  undeploy         Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.

Build Dependencies
  kustomize        Download kustomize locally if necessary. If wrong version is installed, it will be removed before downloading.
  controller-gen   Download controller-gen locally if necessary. If wrong version is installed, it will be overwritten.
  envtest          Download envtest-setup locally if necessary.
```   


## Instructions
```console
go mod init github.com/stackzoo/khaos && 
kubebuilder init --domain stackzoo.io --repo stackzoo.io/khaos &&
kubebuilder create api --group khaos --version v1alpha1 --kind PodDestroyer
```   

Install:
```console
minikube start --driver=docker --profile operator-cluster --memory 8192 --cpus 4 && 
make install && kubectl get crds

NAME                              CREATED AT
poddestroyers.khaos.stackzoo.io   2023-11-27T16:29:37Z
```  

## Local Debug
  
```console
make run
&& kubectl apply -f examples
```  

In order to debug this project locally, I strongly suggest using [vscode](https://code.visualstudio.com/).  

In vscode you need to create a`.vscode/launch.json` similar to the following:  
```json
{
    "version": "0.2.0",
    "configurations": [
      {
        "name": "Debug Khaos Operator",
        "type": "go",
        "request": "launch",
        "mode": "auto",
        "program": "${workspaceFolder}/cmd/main.go",
        "args": []
      }
    ]
  }
```   




