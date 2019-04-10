# ownership-controller
Kubernetes custom controller that controls all owned resources to be successfully created

## Preface
This application is designed to familiarize myself with the CRD (custom resource definition) and CRD controller which are Kubernetes core features that allow make Kubernetes API extensions very easy.
`client-go` library was used in main as it has all necessary tools for communication with Kubernetes cluster

## Prerequisites
Ownership controller is designed to work inside of cluster so you should have a configured one. You can get acquainted how to configure your own Kubernetes cluster on Ubuntu 18 nice and easy using [this article](https://linuxconfig.org/how-to-install-kubernetes-on-ubuntu-18-04-bionic-beaver-linux).
Also you should have `go` installed on your workstation and `$GOPATH` must be set.

## Building
Get the repo with:

`go get github.com/kacejot/ownership-controller`

Go to the project directory:

`cd $GOPATH/src/github.com/kacejot/ownership-controller`

Run codegen sript to generate helper types for controller work:

`./scripts/update_codegen.sh`

Build and push the docker image with it to Docker Hub:

`docker build --no-cache -t "<user>/<tag>" . && docker push <user>/<tag>`

## Running controller
Let the project root dir be `$PROJECT_ROOT`. `$PROJECT_ROOT/yaml/crd.yaml` is a file that defines custom resource type called `Owner`. Simply register it in your cluster:

`kubectl create -f $PROJECT_ROOT/yaml/crd.yaml`

In `$PROJECT_ROOT/yaml` you have `deploy-controller.yaml`. This file is responsible for controller creation inside of the cluster. It has several resources: `ServiceAccount`, `ClusterRoleBinding` and `Deployment`. `ClusterRoleBinding` binds cluster-admin role to our `ServiceAccount`. That is necessary to allow controller to list `Owner` resource in cluster scope. `Deployment` is resource that pulls down image with controller and runs it in the pod. Write your image in `image` field and launch the controller inside the cluster:

`kubectl create -f $PROJECT_ROOT/yaml/deploy-controller.yaml`

Now controller works and CRD is registered on the cluster. The last thing that left is to test its work. Create resource with `Owner` type like described in `owner-example.yaml` and `error-owner-example.yaml`

You also can check controller work with logs and pod status. Check that pod is running: 

`kubectl get pods -n kube-system`

Check controller logs:

`kubectl logs -n kube-system controller-******`

## Owner usage
You can now use `Owner` to be sure that all resources that owner owns will be deleted if someone fails on its creation.
Owner structure:`Owner` contains array of owned resource. Each resource must be represented with 3 fields: namespace, name and resource. Resource field is a plural form for your resource kind definded in its CRD. All the standard resource kinds has plural form by default.

Example of usage:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: sample-service-account-1
  namespace: kube-system
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: sample-service-account-2
  namespace: kube-system
---
apiVersion: myproject.com/v1alpha1
kind: Owner
metadata:
  name: accounts-owner
spec:
  ownedResources:
  - resource: serviceaccounts
    name: sample-service-account-1
    namespace: kube-system
  - resource: serviceaccounts
    name: sample-service-account-2
    namespace: kube-system
```
