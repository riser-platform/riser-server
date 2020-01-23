[Kube Applier](https://github.com/box/kube-applier) is used to automatically make changes to the kube cluster based on what is in git. This works in tandem with the [Git Sync](https://github.com/kubernetes/git-sync) sidecar to synchronize git changes to a kubernetes cluster into each service. these resources must be applied directly to the kubernetes cluster and not committed to the `riser-state` repo.


## Configuration
You must configure a configMap and secret for the kube applier and git sync services. The key values map to environment variables within each service.

```
kubectl create configmap kube-applier	--namespace kube-applier --from-literal=REPO_PATH=/git-repo/<REPO_NAME>/stages/<STAGE_NAME>/kube-resources
kubectl create secret generic kube-applier --namespace kube-applier --from-literal=GIT_SYNC_REPO=<STATE_REPO_URL>
```

Example for stage `dev` in repo `my-state-repo`

```
kubectl create configmap kube-applier	--namespace kube-applier --from-literal=REPO_PATH=/git-repo/my-state-repo/stages/dev/kube-resources
kubectl create secret kube-applier --namespace kube-applier --from-literal=GIT_SYNC_REPO=https://github.com/my-org/my-state-repo
```

> :information_source: If the `riser-state` is private be sure to include the username/password or auth token with read access e.g. For github  `https://oauthtoken:<YOUR-OAUTH-TOKEN>@github.com/...`. It is possible to use a [SSH key instead](https://github.com/kubernetes/git-sync/blob/master/docs/ssh.md).
