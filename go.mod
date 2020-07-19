module github.com/riser-platform/riser-server

replace github.com/riser-platform/riser-server/api/v1/model => ./api/v1/model

go 1.14

require (
	github.com/bitnami-labs/sealed-secrets v0.12.4
	github.com/dustin/go-humanize v1.0.0
	github.com/go-ozzo/ozzo-validation/v3 v3.8.1
	github.com/golang-migrate/migrate/v4 v4.11.0
	github.com/google/uuid v1.1.1
	github.com/imdario/mergo v0.3.10
	github.com/joho/godotenv v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/labstack/echo/v4 v4.1.16
	github.com/lib/pq v1.7.0
	github.com/onrik/logrus v0.7.0
	github.com/pkg/errors v0.9.1
	github.com/riser-platform/riser-server/api/v1/model v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.1
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.17.9
	k8s.io/apimachinery v0.17.9
	k8s.io/client-go v0.17.9
	knative.dev/pkg v0.0.0-20191101194912-56c2594e4f11
	sigs.k8s.io/yaml v1.2.0
)
