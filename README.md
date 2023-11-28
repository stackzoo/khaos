# KHAOS
[![releaser](https://github.com/stackzoo/khaos/actions/workflows/release.yaml/badge.svg)](https://github.com/stackzoo/khaos/actions/workflows/release.yaml)  [![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)  [![Latest Release](https://img.shields.io/github/release/stackzoo/khaos.svg)](https://github.com/stackzoo/khaos/releases/latest)  

<br/>
<img src="docs/images/klogo.png" alt="logo" width="230" height="230">  
<br/>
A lightweight kubernetes operator to test cluster and application resilience via chaos engineering 💣 ☸️  

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


## Local Testing and Debugging
First of all clone the repository:  
```console
git clone https://github.com/stackzoo/khaos && cd khaos
```  

The repo contains a `Makefile` with all that you need.  
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

You can spin up a local dev cluster with [KinD](https://kind.sigs.k8s.io/) with the following command:  
```console
make cluster-up
```   

Install and list the operator CRDs with the following command:  
```console
make install && kubectl get crds

NAME                                       CREATED AT
commandinjections.khaos.stackzoo.io        2023-11-28T12:55:25Z
containerresourcechaos.khaos.stackzoo.io   2023-11-28T12:55:25Z
nodedestroyers.khaos.stackzoo.io           2023-11-28T12:55:25Z
poddestroyers.khaos.stackzoo.io            2023-11-28T12:55:25Z
secretdestroyers.khaos.stackzoo.io         2023-11-28T12:55:25Z
```  

In order to run the operator on your cluster (current context - i.e. whatever cluster `kubectl cluster-info` shows).) run:  
```console
make run
```  


In order to debug this project locally, I strongly suggest using [vscode](https://code.visualstudio.com/).  

In vscode you need to create a `.vscode/launch.json` file similar to the following:  
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

## Examples


## Examples

<details>
  <summary>Delete Pods</summary>

Create a new namespace called `prod` and apply an example deployment:  

```console
kubectl create namespace prod && kubectl apply -f examples/test-deployment.yaml
```  

Wait for all the pods to be up and running and then apply the `PodDestroyer` manifest:  


```console
kubectl apply -f examples/pod-destroyer.yaml
```

Now you can observe 2 things:  
1. the pods in prod namespace are being Terminating (and recreated by the replicaset):  
```console
NAME                                READY   STATUS              RESTARTS   AGE
nginx-deployment-7bf8c77b5b-5fvrc   1/1     Running             0          6s
nginx-deployment-7bf8c77b5b-5qcx4   1/1     Running             0          6s
nginx-deployment-7bf8c77b5b-6kmbd   0/1     ContainerCreating   0          6s
nginx-deployment-7bf8c77b5b-75bg6   1/1     Running             0          6s
nginx-deployment-7bf8c77b5b-bcbk5   1/1     Running             0          6s
nginx-deployment-7bf8c77b5b-f5wkh   1/1     Running             0          6s
nginx-deployment-7bf8c77b5b-gfdzl   1/1     Running             0          6s
nginx-deployment-7bf8c77b5b-gmhr2   1/1     Running             0          6s
nginx-deployment-7bf8c77b5b-gsprh   1/1     Terminating         0          6s
nginx-deployment-7bf8c77b5b-hvsff   1/1     Running             0          6s
nginx-deployment-7bf8c77b5b-v4j9v   0/1     ContainerCreating   0          6s
nginx-deployment-7bf8c77b5b-zxxv7   0/1     Terminating         0          6s
nginx-deployment-7bf8c77b5b-6kmbd   1/1     Running             0          6s
nginx-deployment-7bf8c77b5b-zxxv7   0/1     Terminating         0          6s
nginx-deployment-7bf8c77b5b-zxxv7   0/1     Terminating         0          6s
nginx-deployment-7bf8c77b5b-zxxv7   0/1     Terminating         0          6s
nginx-deployment-7bf8c77b5b-v4j9v   1/1     Running             0          7s
nginx-deployment-7bf8c77b5b-gsprh   0/1     Terminating         0          32s
nginx-deployment-7bf8c77b5b-gsprh   0/1     Terminating         0          33s
nginx-deployment-7bf8c77b5b-gsprh   0/1     Terminating         0          33s
nginx-deployment-7bf8c77b5b-gsprh   0/1     Terminating         0          33s
```  
2. Our operator shows the reconciliation logic's logs:  
```console   
2023-11-28T14:07:18+01:00       INFO    Reconciling PodDestroyer: default/nginx-destroyer       {"controller": "poddestroyer", "controllerGroup": "khaos.stackzoo.io", "controllerKind": "PodDestroyer", "PodDestroyer": {"name":"nginx-destroyer","namespace":"default"}, "namespace": "default", "name": "nginx-destroyer", "reconcileID": "1e16a7d2-825a-4b46-b4e5-ac1228bc1c36"}
2023-11-28T14:07:18+01:00       INFO    Selector: {map[app:nginx] []}   {"controller": "poddestroyer", "controllerGroup": "khaos.stackzoo.io", "controllerKind": "PodDestroyer", "PodDestroyer": {"name":"nginx-destroyer","namespace":"default"}, "namespace": "default", "name": "nginx-destroyer", "reconcileID": "1e16a7d2-825a-4b46-b4e5-ac1228bc1c36"}
2023-11-28T14:07:18+01:00       INFO    MaxPods: 9      {"controller": "poddestroyer", "controllerGroup": "khaos.stackzoo.io", "controllerKind": "PodDestroyer", "PodDestroyer": {"name":"nginx-destroyer","namespace":"default"}, "namespace": "default", "name": "nginx-destroyer", "reconcileID": "1e16a7d2-825a-4b46-b4e5-ac1228bc1c36"}
2023-11-28T14:07:18+01:00       INFO    Namespace: prod {"controller": "poddestroyer", "controllerGroup": "khaos.stackzoo.io", "controllerKind": "PodDestroyer", "PodDestroyer": {"name":"nginx-destroyer","namespace":"default"}, "namespace": "default", "name": "nginx-destroyer", "reconcileID": "1e16a7d2-825a-4b46-b4e5-ac1228bc1c36"}
```  

Now we can inspect the status of our PodDestroyer custom resource:  
```console 
kubectl get poddestroyer

NAME              AGE
nginx-destroyer   4m51s
```  

```console
kubectl get poddestroyer nginx-destroyer -o yaml
```  
This will retrieve our resource in `yaml` format:  
```yaml
apiVersion: khaos.stackzoo.io/v1alpha1
kind: PodDestroyer
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"khaos.stackzoo.io/v1alpha1","kind":"PodDestroyer","metadata":{"annotations":{},"name":"nginx-destroyer","namespace":"default"},"spec":{"maxPods":9,"namespace":"prod","selector":{"matchLabels":{"app":"nginx"}}}}
  creationTimestamp: "2023-11-28T13:07:18Z"
  generation: 1
  name: nginx-destroyer
  namespace: default
  resourceVersion: "2009"
  uid: fbba6287-6f70-406b-821e-9000f097afc5
spec:
  maxPods: 9
  namespace: prod
  selector:
    matchLabels:
      app: nginx
status:
  numPodsDestroyed: 9
```  

The `status` spec tells you how many pods have been successfully destroyed.  


</details>  



<details>
  <summary>Delete Nodes</summary>

First, retrieve nodes info for your cluster:  
```console
kubectl get nodes

NAME                                  STATUS   ROLES           AGE   VERSION
test-operator-cluster-control-plane   Ready    control-plane   24m   v1.27.3
test-operator-cluster-worker          Ready    <none>          24m   v1.27.3
test-operator-cluster-worker2         Ready    <none>          24m   v1.27.3
test-operator-cluster-worker3         Ready    <none>          24m   v1.27.3

```  

Now apply the following `NodeDestroyer` manifest:  

```yaml
apiVersion: khaos.stackzoo.io/v1alpha1
kind: NodeDestroyer
metadata:
  name: example-node-destroyer
spec:
  nodeNames:
    - test-operator-cluster-worker
    - test-operator-cluster-worker3
```

```console
kubectl apply -f examples/node-destroyer.yaml
```

Now, once again, retrieve the node list from the kuber-apiserver:  
```console
kubectl get nodes

NAME                                  STATUS   ROLES           AGE   VERSION
test-operator-cluster-control-plane   Ready    control-plane   25m   v1.27.3
test-operator-cluster-worker2         Ready    <none>          25m   v1.27.3

```  

As you can see the operator succesfully removed the specified nodes.  


</details>  



<details>
  <summary>Delete Secrets</summary>

First create a new kubernetes secret (empty secret is fine):  

```console
kubectl -n prod create secret generic test-secret

secret/test-secret created
```  

Now apply the following `SecretDestroyer` manifest:  

```yaml
apiVersion: khaos.stackzoo.io/v1alpha1
kind: SecretDestroyer
metadata:
  name: example-secret-destroyer
spec:
  namespace: prod
  secretNames:
    - test-secret
```

```console
kubectl apply -f examples/secret-destroyer.yaml
```  

Try to list all the secrets in the `prod` namespace:  
```console
kubectl -n prod get secrets

No resources found in prod namespace.
```  

The specified secret was successfully removed.  



</details>  


<details>
  <summary>Apply New Container Resource Limits</summary>  

Apply the following `ContainerResourceChaos` manifest:  

```yaml
apiVersion: khaos.stackzoo.io/v1alpha1
kind: ContainerResourceChaos
metadata:
  name: example-container-resource-chaos
  namespace: prod
spec:
  namespace: prod
  DeploymentName: nginx-deployment
  containerName: nginx
  maxCPU: "666m"
  maxRAM: "512Mi"

```  

```console
kubectl apply -f examples/container-resource-chaos.yaml
```  

Now retrieve one of the pod in the prod namespace in `yaml` format and take a look at the resources:  
```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: "2023-11-28T13:43:37Z"
  generateName: nginx-deployment-c54b8b4b4-
  labels:
    app: nginx
    pod-template-hash: c54b8b4b4
  name: nginx-deployment-c54b8b4b4-jvw4k
  namespace: prod
  ownerReferences:
  - apiVersion: apps/v1
    blockOwnerDeletion: true
    controller: true
    kind: ReplicaSet
    name: nginx-deployment-c54b8b4b4
    uid: a73e8483-a51b-4f43-806d-38b8976ee61d
  resourceVersion: "6128"
  uid: 6be9fe17-f6b8-418b-96a1-bdf70da8eb95
spec:
  containers:
  - image: nginx:latest
    imagePullPolicy: Always
    name: nginx
    resources: # modified
      limits:
        cpu: 666m
        memory: 512Mi
      requests:
        cpu: 666m
        memory: 512Mi
```   


</details>  


<br/>  


## Operator Installation
This repo contains a [github action](https://github.com/stackzoo/khaos/blob/main/.github/workflows/release.yaml) that publish  the operator *oci image*  to *github registry* when new releases tag are pushed to the main branch.  
In order to install the operator as a pod in the cluster you can leverage one of the *make* targets:  
```console
make deploy IMG=ghcr.io/stackzoo/khaos:0.0.3
```  

This command will install all the required *CRDs* and *RBAC manifests* and then start the operator as a pod:  
```console
kubectl get pods -n khaos-system

NAME                                       READY   STATUS             RESTARTS   AGE
khaos-controller-manager-8887957bf-5b8g9   1/1     Running               0       107s
```  

> [!NOTE]  
> If you encounter RBAC errors, you may need to grant yourself cluster-admin privileges or be logged in as admin.


