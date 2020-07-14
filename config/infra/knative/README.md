A base [Kustomization](https://kustomize.io/) is provided for generation of Knative Serving manifests that work with Riser. The base manifests are:
- [serving-core](https://github.com/knative/serving/releases)
- [net-istio](https://github.com/knative/net-istio/releases)
- [net-cert-manager](https://github.com/knative/net-certmanager/releases)

To build run `kustomize build base >/path/to/gitops/repo/kustomization.yaml`

> :information_source: _This is not a production ready Knative Serving configuration. See the [KNative install documentation](https://knative.dev/docs/install/) which contains detailed installation guides depending on your Kubernetes environment._

## Configuring KNative
KNative configuration is vast and is dependant largely on your needs. The following is meant to help get you started with a basic demo of Riser and is not intended as being exhaustive.

### Domain Name
You should configure a wildcard domain for each environment and namespace using a pattern like `<environment>.riser.<your-domain>` (e.g. for the `dev` environment `dev.riser.your-domain.org`. To do this, create a ConfigMap like the following example:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-domain
  namespace: knative-serving
data:
  dev.riser.your-domain.org: ""
```



