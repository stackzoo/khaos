# KHAOS
[![releaser](https://github.com/stackzoo/khaos/actions/workflows/release.yaml/badge.svg)](https://github.com/stackzoo/khaos/actions/workflows/release.yaml)  [![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)  [![Latest Release](https://img.shields.io/github/v/release/stackzoo/khaos?logo=github)](https://github.com/stackzoo/khaos/releases/latest)  ![GitHub last commit (branch)](https://img.shields.io/github/last-commit/stackzoo/khaos/main)


<br/>
<img src="docs/images/klogo.png" alt="logo" width="230" height="230">  
<br/>
A lightweight kubernetes operator to test cluster and application resilience via chaos engineering 💣 ☸️  

## Abstract
**Khaos** is a straightforward Kubernetes [operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) made with [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) and designed for executing [Chaos Engineering](https://en.wikipedia.org/wiki/Chaos_engineering) activities.  
Through the implementation of custom controllers and resources, Khaos facilitates the configuration and automation of operations such as the targeted deletion of pods within a specified namespace, the removal of nodes from the cluster, the deletion of secrets and more.  
Khaos is an **unopinionated** operator, in the sense that it provides simple and *atomic primitives* that engineers can use as building blocks in order to compose their preferred chaos strategy.  
Currently, Khaos does not implement *cronjobs*; any scheduling of Khaos Custom Resources is delegated to external logic outside the cluster, possibly through a GitOps approach.  

> [!WARNING]  
> This operator will introduce faults and unpredicatbility in your infrastructure, use with caution  

## Supported features
- [X] Delete pods
- [x] Delete cluster nodes
- [X] Delete secrets
- [X] Delete configmaps
- [X] Cordon nodes
- [X] Inject resource constraints in pods
- [X] Add o remove labels in pods
- [X] Flood api server with calls
- [X] Consume resources in a namespace
- [X] Create custom kubernetes events
- [X] Exec commands inside pods (**experimental**).  



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

You can spin up a local dev cluster with [KinD](https://kind.sigs.k8s.io/) via the following command:  
```console
make cluster-up
```   

Install and list all the available operator's CRDs with the following command:  
```console
make install && kubectl get crds

NAME                                       CREATED AT
apiserveroverloads.khaos.stackzoo.io          2023-12-06T13:20:49Z
commandinjections.khaos.stackzoo.io           2023-12-06T13:20:49Z
configmapdestroyers.khaos.stackzoo.io         2023-12-06T13:20:49Z
consumenamespaceresources.khaos.stackzoo.io   2023-12-06T13:20:49Z
containerresourcechaos.khaos.stackzoo.io      2023-12-06T13:20:49Z
cordonnodes.khaos.stackzoo.io                 2023-12-06T13:20:49Z
eventsentropies.khaos.stackzoo.io             2023-12-06T13:20:49Z
nodedestroyers.khaos.stackzoo.io              2023-12-06T13:20:49Z
poddestroyers.khaos.stackzoo.io               2023-12-06T13:20:49Z
podlabelchaos.khaos.stackzoo.io               2023-12-06T13:20:49Z
secretdestroyers.khaos.stackzoo.io            2023-12-06T13:20:49Z
```  

In order to run the operator on your cluster (current context - i.e. whatever cluster `kubectl cluster-info` shows) run:  
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



## Some Examples

In order to test the following examples, you can use the local *KinD* cluster (see the `Local Testing and Debugging` section).  
Once you have the cluster up and running, procede to create a new namespace called `prod` and apply an example deployment:  

```console
kubectl create namespace prod && kubectl apply -f examples/test-deployment.yaml
```  

Now you can procede with the examples!  

<details>
  <summary>DELETE PODS</summary>

Wait for all the pods in the `prod` namespace to be up and running and then apply the `PodDestroyer` manifest:  

```yaml
apiVersion: khaos.stackzoo.io/v1alpha1
kind: PodDestroyer
metadata:
  name: nginx-destroyer
spec:
  selector:
    matchLabels:
      app: nginx
  maxPods: 9
  namespace: prod
```  



```console
kubectl apply -f examples/pod-destroyer.yaml
```

Now you can observe 2 things:  
1. the pods in prod namespace are being Terminated (and recreated by the replicaset):  
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
1. Our operator shows the reconciliation logic's logs:  
```console   
2023-11-28T14:07:18+01:00       INFO    Reconciling PodDestroyer: default/nginx-destroyer       {"controller": "poddestroyer", "controllerGroup": "khaos.stackzoo.io", "controllerKind": "PodDestroyer", "PodDestroyer": {"name":"nginx-destroyer","namespace":"default"}, "namespace": "default", "name": "nginx-destroyer", "reconcileID": "1e16a7d2-825a-4b46-b4e5-ac1228bc1c36"}
2023-11-28T14:07:18+01:00       INFO    Selector: {map[app:nginx] []}   {"controller": "poddestroyer", "controllerGroup": "khaos.stackzoo.io", "controllerKind": "PodDestroyer", "PodDestroyer": {"name":"nginx-destroyer","namespace":"default"}, "namespace": "default", "name": "nginx-destroyer", "reconcileID": "1e16a7d2-825a-4b46-b4e5-ac1228bc1c36"}
2023-11-28T14:07:18+01:00       INFO    MaxPods: 9      {"controller": "poddestroyer", "controllerGroup": "khaos.stackzoo.io", "controllerKind": "PodDestroyer", "PodDestroyer": {"name":"nginx-destroyer","namespace":"default"}, "namespace": "default", "name": "nginx-destroyer", "reconcileID": "1e16a7d2-825a-4b46-b4e5-ac1228bc1c36"}
2023-11-28T14:07:18+01:00       INFO    Namespace: prod {"controller": "poddestroyer", "controllerGroup": "khaos.stackzoo.io", "controllerKind": "PodDestroyer", "PodDestroyer": {"name":"nginx-destroyer","namespace":"default"}, "namespace": "default", "name": "nginx-destroyer", "reconcileID": "1e16a7d2-825a-4b46-b4e5-ac1228bc1c36"}
```  

Now we can inspect the status of our PodDestroyer object:  
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
  <summary>DELETE NODES</summary>

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
  <summary>DELETE SECRETS</summary>

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
  <summary>DELETE CONFIGMAPS</summary>

First create a new kubernetes configmap:  

```console
kubectl create configmap test-configmap --namespace=prod --from-literal=message=ready && kubectl -n prod get configmap

configmap/test-configmap created

NAME               DATA   AGE
kube-root-ca.crt   1      2m24s
test-configmap     1      1s

```  

Now apply the following `ConfigMapDestroyer` manifest:  

```yaml
apiVersion: khaos.stackzoo.io/v1alpha1
kind: ConfigMapDestroyer
metadata:
  name: example-configmap-destroyer
spec:
  namespace: prod
  configMapNames:
    - test-configmap
```

```console
kubectl apply -f examples/config-map-destroyer.yaml
```  

Try to list all the configmaps in the `prod` namespace:  
```console
kubectl -n prod get configmap

NAME               DATA   AGE
kube-root-ca.crt   1      9m26s
```  

The specified configmap was successfully removed.  

</details>  


<details>
  <summary>APPLY NEW CONTAINER RESOURCE LIMITS</summary>  

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




<details>
  <summary>MODIFY POD LABELS</summary>  

Apply the following `PodLabelChaos` manifest:  

```yaml
apiVersion: khaos.stackzoo.io/v1alpha1
kind: PodLabelChaos
metadata:
  name: podlabelchaos-test
spec:
  deploymentName: nginx-deployment
  namespace: prod
  labels:
    chaos: "true"
  addLabels: true

```  

```console
kubectl apply -f examples/pod-label-chaos.yaml
```  

Now retrieve one of the pod in the prod namespace in `yaml` format and take a look at the labels:  
```yaml

apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: "2023-11-28T15:27:22Z"
  generateName: nginx-deployment-6bb89bf6cd-
  labels:
    app: nginx
    chaos: "true"
    pod-template-hash: 6bb89bf6cd
  name: nginx-deployment-6bb89bf6cd-52j42
  namespace: prod

```   


</details>  



<details>
  <summary>CONSUME NAMESPACE RESOURCES</summary>  
This feature of the operator will spin up a busybox deployment with the specified replicas in the specified namespace.  
All the busybox's pod will execute the following command:  

```console
while true; do echo 'Doing extensive tasks'; sleep 1; done
```  


First of all we need to install the **metrics server** on our cluster:  
```console
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml  \
&& kubectl patch -n kube-system deployment metrics-server --type=json -p '[{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--kubelet-insecure-tls"}]'
```   
Wait for the metric server pod to be up and running and check cluster (nodes) resources:  
```console
kubectl top nodes

NAME                                  CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%
test-operator-cluster-control-plane   221m         2%     697Mi           4%
test-operator-cluster-worker          31m          0%     230Mi           1%
test-operator-cluster-worker2         29m          0%     253Mi           1%
test-operator-cluster-worker3         42m          0%     242Mi           1%
```  


Now apply the following `ConsumeNamespaceResources` manifest:  

```yaml
apiVersion: khaos.stackzoo.io/v1alpha1
kind: ConsumeNamespaceResources
metadata:
  name: example-consume-resources
spec:
  targetNamespace: prod
  numPods: 200

```  

```console
kubectl apply -f examples/consume-namespace-resources.yaml  
```  

Let's inspect the deployment in the `prod` namespace:  
```console
kubectl -n prod get deployment

NAME                 READY   UP-TO-DATE   AVAILABLE     AGE
busybox-deployment   200/200   80           80          44s
nginx-deployment     10/10     10           10          10m
```  

Let's now review the nodes usagge:  
```console
kubectl top nodes

NAME                                  CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%   
test-operator-cluster-control-plane   845m         10%    904Mi           5%        
test-operator-cluster-worker          1790m        22%    938Mi           5%        
test-operator-cluster-worker2         1494m        18%    1039Mi          6%        
test-operator-cluster-worker3         1673m        20%    1045Mi          6%
```  

As we can see, our deployment in the *prod* namespace is consuming resources!  
Now try deleting the `ConsumeNamespaceResources` object:  
```console
kubectl delete -f examples/consume-namespace-resources.yaml

consumenamespaceresources.khaos.stackzoo.io "example-consume-resources" deleted
```   

Check the operator's logs:  

```console
2023-11-30T15:45:40+01:00       INFO    Object deleted, finalizing resources    {"controller": "consumenamespaceresources", "controllerGroup": "khaos.stackzoo.io", "controllerKind": "ConsumeNamespaceResources", "ConsumeNamespaceResources": {"name":"example-consume-resources","namespace":"default"}, "namespace": "default", "name": "example-consume-resources", "reconcileID": "b35fdd79-5308-4080-a718-027e2d9d7d13"}
```  

The resource's controller contains a finalizer and it is deleting our busybox deployment in the *prod* namespace!  
Check the deployments in the *prod* namespace:  
```console
kubectl -n prod get deployment

NAME               READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   10/10   10           10          21m
```  
Cool, our deployment has been successfully deleted.  



</details>  



<details>
  <summary>CREATE CUSTOM KUBERNETES EVENTS</summary>  

Apply the following `EventsEntropy` manifest:  

```yaml
apiVersion: khaos.stackzoo.io/v1alpha1
kind: EventsEntropy
metadata:
  name: example-eventsentropy
spec:
  events:
    - "Custom event 1 with some gibberish - dfsdfsdffdgt egeg4e 😊"
    - "Custom event 2 - with some gibberish dfsdfsdffdgt 676565 🥴"
    - "Custom event 3 - with some gibberish 8/ihfwgf sufdh  🤪"

```  

```console
kubectl apply -f examples/events-entropy.yaml
```  

Now retrieve kubernetes events via kubectl:  
```console
kubectl get events | grep gibberish

<unknown>               Custom event 1 with some gibberish - dfsdfsdffdgt egeg4e 😊
<unknown>               Custom event 3 - with some gibberish 8/ihfwgf sufdh  🤪
<unknown>               Custom event 2 - with some gibberish dfsdfsdffdgt 676565 🥴

```   


</details>  



<details>
  <summary>CORDON NODES</summary>  

Apply the following `CordonNodes` manifest:  

```yaml
apiVersion: khaos.stackzoo.io/v1alpha1
kind: CordonNode
metadata:
  name: example-cordon-node
spec:
  nodesToCordon:
    - test-operator-cluster-worker
    - test-operator-cluster-worker2
    - test-operator-cluster-worker3

```  

```console
kubectl apply -f examples/cordon-nodes.yaml
```  

Now check the status of the resource:  

```console
kubectl describe cordonnodes.khaos.stackzoo.io example-cordon-node | grep "Nodes Cordoned"

Nodes Cordoned:  3
```   


Now run a busybox pod:  
```console
kubectl apply -f examples/test-node-cordon-pod.yaml

pod/busybox-pod created
```   

Let's check that pod:  
```console
kubectl -n default describe pod busybox-pod | grep Warning

Warning  FailedScheduling  63s   default-scheduler  0/4 nodes are available: 1 node(s) had untolerated taint {node-role.kubernetes.io/control-plane: }, 3 node(s) were unschedulable. preemption: 0/4 nodes are available: 4 Preemption is not helpful for scheduling..
```  

</details>  


<br/>  


## Operator Installation
This repo contains a [github action](https://github.com/stackzoo/khaos/blob/main/.github/workflows/release.yaml) that publish  the operator *oci image*  to *github registry* when new releases tag are pushed to the main branch.  
In order to install the operator as a pod in the cluster you can leverage one of the *make* targets:  
```console
make deploy IMG=ghcr.io/stackzoo/khaos:0.0.8
```  

This command will install all the required *CRDs* and *RBAC manifests* and then start the operator as a pod:  
```console
kubectl get pods -n khaos-system

NAME                                       READY   STATUS             RESTARTS   AGE
khaos-controller-manager-8887957bf-5b8g9   1/1     Running               0       107s
```  

> [!NOTE]  
> If you encounter RBAC errors, you may need to grant yourself cluster-admin privileges or be logged in as admin.  
  



## Useful References

- [kubuilder docs](https://book.kubebuilder.io/)
- [programming kubernetes book](https://www.oreilly.com/library/view/programming-kubernetes/9781492047094/)
- [kubernetes programming with go book](https://link.springer.com/book/10.1007/978-1-4842-9026-2)
- [chaos engineering book](https://www.oreilly.com/library/view/chaos-engineering/9781492043850/)
