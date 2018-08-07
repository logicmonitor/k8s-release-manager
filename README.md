> **Note:** Release Manager is a community driven project. LogicMonitor support
will not assist in any issues related to Release Manager.

## Release Manager is a tool for importing and exporting Helm release state

-  **Export Helm release state:**
Release Manager can contact Tiller in the configured cluster, collect
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
    - S3
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

To export releases, Release Manager queries Tiller to collect metadata for all
releases currently deployed in the source cluster and writes this metadata to
the configured backend data store. If the Release Manager is deployed in
daemon mode via its own Helm chart, it will also store metadata about itself.
This metadata is used to prevent import operations from creating a new Release
Manager with the same configuration as the previous managed, causing both
instances to write conflicting state to the backend.

To import releases, Release Manager retrieves the state stored in the backend,
connects to Tiller in the target Kubernetes cluster, and deploys the saved
releases to the cluster.

Release Manager will use --kubeconfig/--kubecontext, $KUBECONFIG, or
~/.kube/config to establish a connection to the Kubernetes cluster. If none of
these configuraitons are set, an in-cluster connection will be attempted. All
actions will be performed against the current cluster and a given command will
only perform actions against a single cluster, i.e. 'export' will
export releases from the configured cluster while 'import' will deploy releases
to the configured cluster and 'clear' requires no custer connection whatsoever.

## Installing via Helm Chart
Installing releasemanager daemon via Helm chart
```
helm repo add logicmonitor https://logicmonitor.github.io/k8s-helm-charts
helm install logicmonitor/releasemanager \
  --set path=$BACKEND_STORAGE_PATH \
  --name releasemanager-$CURRENT_CLUSTER
```

## Use cases
### Redeploying applications after a blue/green switch
### Deploying applications in a disaster recovery scenario

## License
[![license](https://img.shields.io/github/license/logicmonitor/k8s-argus.svg?style=flat-square)](https://github.com/logicmonitor/k8s-argus/blob/master/LICENSE)
