The files in the [apply](./apply) folder provide a basic install of the KNative Serving components from [v0.11.0](https://github.com/knative/serving/releases/tag/v0.11.1). See the [KNative installation](https://knative.dev/docs/install/) documentation which contains detailed guides depending on your Kubernetes environment. Note that Riser only requires that the Serving components are installed. Due to a bug in `kubectl` you may have to run `kubectl apply` again after a few moments after the CRDs are finished registering.

## Configuring KNative
Like Kubernetes, KNative configuration is vast (although much more constrained) and is dependant largely on your needs. The following is meant to help get you started with a basic demo of Riser and is not intended as being exhaustive.

### Domain Name
You should configure a wildcard domain for each stage and namespace using a pattern like `<stage>.riser.<your-domain>` (e.g. for the `dev` stage `dev.riser.your-domain.org`. To do this, create a ConfigMap like the following example:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-domain
  namespace: knative-serving
data:
  dev.riser.your-domain.org: ""
```

If hosting the riser server you can add the riser domain as another line e.g. `api.riser.your-domain.org`.

### TLS
It is recommended that TLS is enabled and enforced for all endpoints, even in non production environments. To configure this, add the following ConfigMap:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-network
  namespace: knative-serving
data:
  autoTLS: Enabled
  httpProtocol: Redirected
```

Your app will be accessible via https and http requests will 301 redirect to the https endpoint.

### CertManager

Managing TLS certificates is a complicated subject. You will have to assess your own certificate management strategy along with your security requirements to determine the best approach. The easiest way to manage certificates is with the installed CertManager. The following is an example configuration that uses LetsEncrypt. It assumes that you have created a LetsEncrypt ClusterIssuer.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-certmanager
  namespace: knative-serving
  labels:
    networking.knative.dev/certificate-provider: cert-manager
data:
  issuerRef: |
    kind: ClusterIssuer
    name: letsencrypt
  solverConfig: |
    dns01:
      provider: default
```

[Read more about configuring KNative and CertManager](https://knative.dev/docs/serving/using-auto-tls/)

