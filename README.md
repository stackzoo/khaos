# Khaos
A lightweight kubernetes operator to test cluster resilience via chaos engineering 💥🌪 ☸️

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

## Testing
  
```console
make run
&& kubectl apply -f examples
```  

Debug in VSCODE (`.vscode/launch.json`):
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




