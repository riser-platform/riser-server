module github.com/riser-platform/riser-server/pkg/sdk

replace github.com/riser-platform/riser-server/api/v1/model => ../../api/v1/model

go 1.15

require (
	github.com/google/uuid v1.1.2
	github.com/pkg/errors v0.9.1
	github.com/riser-platform/riser-server/api/v1/model v0.0.21
	github.com/stretchr/testify v1.6.1
)
