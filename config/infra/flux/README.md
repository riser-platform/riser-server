[Flux](https://github.com/fluxcd/flux) is used to automatically make changes to the kube cluster based on the resources committed in git. A basic installation is provided here as an example to help get you up and running quickly. You may wish to review the [helm configuration docs](https://github.com/fluxcd/flux/tree/master/chart/flux#configuration) or [Flux documentation](https://docs.fluxcd.io/) for more advanced installation options.

# Installation

## Create the Namespace

```
kubectl create namespace flux
```

## Configure Flux
Create a secret with your git URL and your git path. The git path
should reflect the riser stage name so that only resources for that stage are applied to this cluster.

```
kubectl create secret generic flux-git --namespace flux --from-literal=GIT_URL=<GIT_URL>
--from-literal=GIT_PATH=state/<STAGE_NAME>/kube-resources
```
> :warning: Do not include a leading "slash" in the `GIT_PATH`

### Example Secret
If your cluster is serving the riser stage named `dev`:

```
kubectl create secret generic flux-git --namespace flux --from-literal=GIT_URL=https://myoathtoken@github.com/myorg/riser-state
--from-literal=GIT_PATH=state/dev/kube-resources
```


> :information_source: If the `riser-state` is private be sure to include the username/password or auth token with read access e.g. For github  `https://<YOUR-OAUTH-TOKEN>@github.com/...`.

> :information_source: If you wish to use an SSH key, you will need to customize the accompanied helm chart or use the `fluxctl` installer. Please refer to the [helm configuration docs](https://github.com/fluxcd/flux/tree/master/chart/flux#configuration) for more information


## Install Flux

Finally, switch to the `flux` namespace and apply the resources in `apply.yaml`:

```
kubectl apply -f apply.yaml --namespace flux
```