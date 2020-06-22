You may apply `istio_operator.yaml` or follow [these instructions](https://istio.io/docs/setup/install/standalone-operator/) to install the Istio operator. You may optionally apply the included `profile-riser-dev.yaml` profile which has been optimized to reduce the footprint for Riser development.

> :warning: *This profile is not a production installation of istio. This should only be used for local development or demo purposes.*

It is recommended that you enable mTLS in STRICT mode (see `peerauthentication.yaml`). It is possible that a future version of Riser will
require mTLS for certain features. See also: https://istio.io/latest/docs/concepts/security/#dependency-on-mutual-tls










