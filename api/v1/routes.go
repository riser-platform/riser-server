package v1

import (
	"database/sql"

	"github.com/riser-platform/riser-server/pkg/deploymentreservation"
	"github.com/riser-platform/riser-server/pkg/environment"

	"github.com/riser-platform/riser-server/pkg/namespace"

	"github.com/riser-platform/riser-server/pkg/rollout"

	"github.com/labstack/echo/v4/middleware"

	"github.com/riser-platform/riser-server/pkg/app"
	"github.com/riser-platform/riser-server/pkg/deployment"
	"github.com/riser-platform/riser-server/pkg/deploymentstatus"
	"github.com/riser-platform/riser-server/pkg/git"
	"github.com/riser-platform/riser-server/pkg/login"
	"github.com/riser-platform/riser-server/pkg/postgres"
	"github.com/riser-platform/riser-server/pkg/secret"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, repo git.Repo, db *sql.DB) {
	v1 := e.Group("/api/v1")

	// TODO: Refactor dependency management
	environmentRepository := postgres.NewEnvironmentRepository(db)
	environmentService := environment.NewService(environmentRepository)
	namespaceRepository := postgres.NewNamespaceRepository(db)
	namespaceService := namespace.NewService(namespaceRepository, environmentRepository)
	deploymentReservationRepository := postgres.NewDeploymentReservationRepository(db)
	appRepository := postgres.NewAppRepository(db)
	appService := app.NewService(appRepository, namespaceService)
	secretMetaRepository := postgres.NewSecretMetaRepository(db)
	secretService := secret.NewService(secretMetaRepository, environmentRepository)
	deploymentReservationService := deploymentreservation.NewService(deploymentReservationRepository)
	deploymentRepository := postgres.NewDeploymentRepository(db)
	deploymentService := deployment.NewService(appRepository, namespaceService, secretMetaRepository, environmentRepository, deploymentRepository, deploymentReservationService)
	deploymentStatusService := deploymentstatus.NewService(deploymentRepository, environmentService)
	rolloutService := rollout.NewService(appRepository, deploymentRepository)
	userRepository := postgres.NewUserRepository(db)
	apiKeyRepository := postgres.NewApiKeyRepository(db)
	loginService := login.NewService(userRepository, apiKeyRepository)

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		// We will probably use the "Bearer" scheme for OIDC
		AuthScheme: "Apikey",
		Validator: func(apikey string, c echo.Context) (bool, error) {
			return loginWithApiKey(c, loginService, apikey)
		},
	}))

	v1.GET("/apps", func(c echo.Context) error {
		return ListApps(c, appRepository)
	})

	v1.GET("/apps/:namespace/:appName", func(c echo.Context) error {
		return GetApp(c, appRepository)
	})

	v1.GET("/apps/:namespace/:appName/status", func(c echo.Context) error {
		return GetAppStatus(c, appService, deploymentStatusService)
	})

	v1.POST("/apps", func(c echo.Context) error {
		return PostApp(c, appService)
	})

	v1.POST("/deployments", func(c echo.Context) error {
		return PostDeployment(c, repo, appService, deploymentService, environmentService)
	})
	v1.PUT("/deployments", func(c echo.Context) error {
		return PostDeployment(c, repo, appService, deploymentService, environmentService)
	})

	v1.DELETE("/deployments/:envName/:namespace/:deploymentName", func(c echo.Context) error {
		return DeleteDeployment(c, repo, deploymentService)
	})

	v1.PUT("/deployments/:envName/:namespace/:deploymentName/status", func(c echo.Context) error {
		return PutDeploymentStatus(c, deploymentRepository)
	})

	v1.PUT("/rollout/:envName/:namespace/:deploymentName", func(c echo.Context) error {
		return PutRollout(c, rolloutService, environmentService, repo)
	})

	v1.PUT("/secrets", func(c echo.Context) error {
		return PutSecret(c, repo, secretService, environmentService)
	})

	v1.GET("/secrets/:envName/:namespace/:appName", func(c echo.Context) error {
		return GetSecrets(c, secretMetaRepository, environmentService)
	})

	v1.GET("/namespaces", func(c echo.Context) error {
		return GetNamespaces(c, namespaceRepository)
	})

	v1.POST("/namespaces", func(c echo.Context) error {
		return PostNamespace(c, namespaceService, repo)
	})

	v1.PUT("/environments/:envName/config", func(c echo.Context) error {
		return PutEnvironmentConfig(c, environmentService)
	})

	v1.POST("/environments/:envName/ping", func(c echo.Context) error {
		return PostEnvironmentPing(c, environmentService)
	})

	v1.GET("/environments", func(c echo.Context) error {
		return ListEnvironments(c, environmentRepository)
	})

	v1.POST("/validate/appconfig", func(c echo.Context) error {
		return PostValidateAppConfig(c, appService, environmentService)
	})
}
