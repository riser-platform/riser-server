A base [Kustomization](https://kustomize.io/) is provided for generation of Knative Serving manifests that work with Riser. The base manifests are:
- [operator](https://github.com/knative/operator/releases)
- [net-certmanager](https://github.com/knative/net-certmanager/releases)

To build run `kustomize build base >/path/to/gitops/repo/kustomization.yaml`

> :information_source: _This is not a production ready Knative Serving configuration. See the [Knative install documentation](https://knative.dev/docs/install/) which contains detailed installation guides depending on your Kubernetes environment._

## Configuring Knative
Knative configuration is vast and is dependant largely on your needs. The following is meant to help get you started with a basic demo of Riser and is not intended as being exhaustive.

### Domain Name
You should configure a wildcard domain for each environment and namespace using a pattern like `<environment>.riser.<your-domain>` (e.g. for the `dev` environment `dev.riser.your-domain.org`. To do this, add your domain to the configuration found in `knative.serving.yaml` file e.g.

```yaml
apiVersion: operator.knative.dev/v1alpha1
kind: KnativeServing
metadata:
  name: knative-serving
  namespace: knative-serving
spec:
  version: 1.0.0
  config:
    # ---v--- example domain configuration
    domain:
      dev.riser.your-domain.org: ""
    # ---^--- example domain configuration
```



