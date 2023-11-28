# Khaos
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)  
<br/>
<img src="docs/images/klogo.png" alt="logo" width="230" height="230">  
<br/>
A lightweight kubernetes operator to test cluster resilience via chaos engineering 💥🌪 ☸️  

## Abstract
**Khaos** is a streamlined Kubernetes [operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) made with [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) and designed for executing [Chaos Engineering](https://en.wikipedia.org/wiki/Chaos_engineering) activities.  
Through the implementation of custom controllers and resources, Khaos facilitates the configuration and automation  
of operations such as the targeted deletion of pods within a specified namespace, the removal of nodes from the cluster, the deletion of secrets and more.  

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




