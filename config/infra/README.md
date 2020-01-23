> :information_source: For a simple single cluster stage demo, try `riser demo install`. This README is for more advanced exploration of the riser platform.

This folder contains demo infrastructure for components required by Riser. These are configured for demonstration purposes only and have not been rigorously tested for stability or security. It is recommended that you configure and install these dependencies using the recommended approach by each respective dependency. This documentation assumes that the reader already has prior experience installing Kubernetes.

### Riser Server
The Riser Server spans all stages and only needs to be installed on one stage. Please review the [server README](../server/README.md) for installing the riser server.

## Create Kubernetes Cluster
While Riser is supported on theoretically any kubernetes cluster, the demo has been tested with GKE. You may wish review the example script in `gke/create.sh`. Once the cluster is created you may continue.

## Install and configure kube-applier
See the [README](kubeapplier/README.md).

## Creating a new stage (one per cluster)
Riser requires a git repo to manage all of its state, referred to as the `riser-state` repo. It is recommended that you use it to manage Riser required infrastructure as well. Note that
at the time of writing that GitHub is the only officially supported git host. Others are planned to be be supported in the future. You may share
a single repo between multiple Riser stages.

Review each required component's README in this folder. With the exception of kube-applier, the final yaml for each component should be placed in your `riser-state` git repo in the
 `/stages/<stageName>/kube-resources/infra` folder. Kube-applier must be installed manually after the cluster is ready. Once you push your changes, `kube-applier` should begin installing the remaining components. This may take a few minutes depending on cluster capacity and internet speed.


# Install and configure KNative
See the [README](knative/README.md).

### Update DNS to point to the istio ingress gateway external IP
Create an A record for your domain (e.g. `*.apps.<stage>.riser.<your-domain>`) using the IP from the command below.

```
kubectl get service istio-ingressgateway -n istio-system -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
```

If hosting the riser server do the same for the riser api e.g. `api.riser.<your-domain>`.

>Note: It may take a few minutes for the gateway to get an IP from the load balancer. A GKE cluster created with the provided script will automatically create the load balancer. Other cloud providers may need additional configuration in order to create a load balancer that points to the istio ingress.

### Optional: Install riser-server secrets
If you're installing the riser server in this cluster, you will have to apply any required secrets as those should not be committed to the `riser-state` repo. See the [server README](../server/README.md)
for more information.

### Optional: Apply cert-manager DNS secrets
For TLS, this demo uses cert-manager with LetsEncrypt using DNS-01 challenges. In order to do this you must be using Google Cloud DNS to manage your domain and cert-manager needs permissions to manage the zone records. The `gke/gcloud-dns-admin-sa.sh` script creates a service account required and generates a secret for kubernetes. After generating apply the secret:

```
kubectl create secret generic cert-manager-credentials -n istio-system --from-file=./gcp-dns-admin.json
```


