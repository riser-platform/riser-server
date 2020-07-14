module github.com/riser-platform/riser-server/pkg/sdk

replace github.com/riser-platform/riser-server/api/v1/model => ../../api/v1/model

go 1.14

require (
	github.com/google/uuid v1.1.1
	github.com/pkg/errors v0.9.1
	github.com/riser-platform/riser-server/api/v1/model v0.0.16
	github.com/stretchr/testify v1.4.0
	k8s.io/apimachinery v0.17.3 // indirect
)
