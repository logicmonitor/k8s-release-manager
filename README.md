> **Note:** Release Manager is a community driven project. LogicMonitor support
will not assist in any issues related to Release Manager.

## Release Manager is a tool for importing and exporting Helm release state

-  **Export Helm release state:**
Release Manager can contact the configured cluster, collect
all metadata for each deployed release, and write that metadata to the
configured backend. This metadata can later be consumed by Release Manager
import to re-install the saved releases to a different cluster.

-  **Continuously export Helm release state:**
Release Manager can also be run in daemon mode to continuously update the
stored state to reflect ongoing changes to the cluster.
When running in daemon mode, it is HIGHLY recommended when running to use the
official [Release Manager Helm chart](#installing-via-helm-chart).

-  **Import Helm release state:**
Release Manager can retrieve stored state from the configured
backend and install all exported releases to the current Kubernetes cluster.

-  **Configurable storage backends:**
Release Manager can interact with multiple release state storage backends.
Currently supported backends are:
    - S3 (**Note:** S3 object versioning is strongly recommended for backend buckets)
    - Local (**Note:** the local backend is intended for non-production use only)

## Release Manager Overview
Release Manager provides functionality for exporting and importing the
state of Helm releases currently deployed to a Kubernetes cluster. The state
of the installed releases is saved to a configurable backend for easy
restoration of previously-deployed releases or for simplified re-deployment of
those releases to a new Kubernetes cluster.

Release Manager operations can be run locally or within the Kubernetes cluster.
The application also supports a daemon mode that will periodically update the
saved state.

To export releases, Release Manager queries the target Kubernetes cluster to 
collect metadata for all releases currently deployed in the source cluster 
and writes this metadata to the configured backend data store. 
If the Release Manager is deployed in daemon mode via its own Helm chart, 
it will also store metadata about itself. This metadata is used to prevent 
import operations from creating a new Release Manager with the same 
configuration as the previous managed, causing both instances to write 
conflicting state to the backend.

To import releases, Release Manager retrieves the state stored in the backend, 
connects to the target Kubernetes cluster, 
and deploys the saved releases to the cluster..

Release Manager will use --kubeconfig/--kubecontext, $KUBECONFIG, or
~/.kube/config to establish a connection to the Kubernetes cluster. If none of
these configuraitons are set, an in-cluster connection will be attempted. All
actions will be performed against the current cluster and a given command will
only perform actions against a single cluster, i.e. 'export' will
export releases from the configured cluster while 'import' will deploy releases
to the configured cluster and 'clear' requires no custer connection whatsoever.

## [Command usage](docs/releasemanager.md)

## Installing via Helm Chart
Installing releasemanager daemon via Helm chart
```
helm repo add logicmonitor https://logicmonitor.github.io/k8s-helm-charts
helm install logicmonitor/releasemanager \
  --set path=$BACKEND_STORAGE_PATH \
  --name releasemanager-$CURRENT_CLUSTER
```

For detailed information about the Helm chart, see the [README](https://github.com/logicmonitor/k8s-helm-charts/blob/master/releasemanager/README.md)

## Viewing installed releases
When running in daemon mode, the Release Manager exposes an endpoint to view
the list of releases currently stored in the backend. This endpoint is
available at /releases.

Note that if the Release Manager is running in-cluster, you'll
need to expose its service via ingress using --set ingress.hosts={...} when
installing the Helm chart.

## Retrieving metrics
When running in daemon mode, the Release Manager exposes an endpoint providing
monitoring metrics. This endpoint is available at /debug/vars
at /debug/vars.

Note that if the Release Manager is running in-cluster, you'll
need to expose its service via ingress using --set ingress.hosts={...} when
installing the Helm chart.

## Use case examples
Release Manager was created with the goal of solving two common cluster
management problems. These use cases are outlined below along with a general
instructions for using Release Manager to solve the problem.

### Redeploying applications after a blue/green switch
Using blue/green deployments for Kubernetes clusters is a great upgrade
strategy; however, when a cluster is running dozens or hundreds of application
deployments supporting several different environments, it can become burdensome
to redploy all of these applications every time a cluster is upgraded. To solve
this problem, Release Manager makes it easy to take a snapshot of applications
deployed to the current cluster and redeploy those applications to the new
cluster.

#### 1. Install Release Manager locally

```shell
go get github.com/logicmonitor/k8s-release-manager/cmd/releasemanager
```

#### 2. Export the releases currently installed in the source cluster

```shell
releasemanager export local \
  --kubeconfig $SOURCE_CLUSTER_KUBECONFIG \
  --path $LOCAL_RELEASE_STATE_PATH
```

#### 3. Deploy the saved releases to the destination cluster

```shell
releasemanager import local \
  --kubeconfig $DESTINATION_CLUSTER_KUBECONFIG \
  --path $LOCAL_RELEASE_STATE_PATH
```

### Deploying applications in a disaster recovery scenario
Some disaster scenarios can result in the need to create a brand new Kubernetes
cluster for failover. When a cluster is running dozens or hundreds of
application deployments, it can become burdensome and incur huge time costs to
redploy all of these applications manually. To solve this problem, Release
Manager can operate as a long-running process inside your cluster to
periodically export snapshots of deployed applications to an external storage
backend. In the event of a disaster recovery operation requiring a cluster
failover, Release Manager can retrieve the stored state and quickly deploy all
of the applications from the failed cluster to the new cluster.

#### 1. Add the LogicMonitor Helm repository

```shell
helm repo add logicmonitor https://logicmonitor.github.io/k8s-helm-charts
```

#### 2. Install Release Manager locally

```shell
go get github.com/logicmonitor/k8s-release-manager/cmd/releasemanager
```

#### 3. Deploy Release Manager in daemon mode when provisioning a production cluster

```shell
helm install logicmonitor/releasemanager \
  --set path=$PROD_CLUSTER_BACKEND_PATH \
  --set s3.bucket=$RELEASE_MANAGER_STATE_BUCKET \
  --set s3.region=$BUCKET_REGION \
  --name releasemanager-$PROD_CLUSTER_NAME
```

#### 4. Provision a failover cluster during a disaster scenario

#### 5. Deploy saved releases to the failover cluster.

**NOTE!** Be sure to --new-path. This prevents Release Manager in failover cluster from overwriting the production state

```shell
releasemanager import s3 \
  --kubeconfig $FAILOVER_CLUSTER_KUBECONFIG \
  --path $PROD_CLUSTER_BACKEND_PATH \
  --new-path $FAILOVER_CLUSTER_BACKEND_PATH \
  --bucket $RELEASE_MANAGER_STATE_BUCKET \
  --region $BUCKET_REGION
```

## License
[![license](https://img.shields.io/github/license/logicmonitor/k8s-argus.svg?style=flat-square)](https://github.com/logicmonitor/k8s-argus/blob/master/LICENSE)
