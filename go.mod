module github.com/riser-platform/riser-server

replace github.com/riser-platform/riser-server/api/v1/model => ./api/v1/model

go 1.13

require (
	github.com/bitnami-labs/sealed-secrets v0.7.0
	github.com/dustin/go-humanize v1.0.0
	github.com/go-ozzo/ozzo-validation/v3 v3.8.1
	github.com/gogo/protobuf v1.2.2-0.20190730201129-28a6bbf47e48 // indirect
	github.com/golang-migrate/migrate/v4 v4.6.2
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.1.1
	github.com/imdario/mergo v0.3.8
	github.com/joho/godotenv v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/labstack/echo/v4 v4.1.6
	github.com/lib/pq v1.1.1
	github.com/onrik/logrus v0.4.0
	github.com/pkg/errors v0.8.1
	github.com/riser-platform/riser-server/api/v1/model v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80 // indirect
	golang.org/x/sys v0.0.0-20190726091711-fc99dfbffb4e // indirect
	google.golang.org/grpc v1.21.0 // indirect
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.0.0-20190627205229-acea843d18eb
	k8s.io/apimachinery v0.0.0-20190629005116-7ae370969693
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v0.3.3 // indirect
	knative.dev/pkg v0.0.0-20191101194912-56c2594e4f11
	sigs.k8s.io/yaml v1.1.0
)
