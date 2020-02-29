package v1

import (
	"database/sql"

	"github.com/riser-platform/riser-server/pkg/deploymentreservation"

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
	"github.com/riser-platform/riser-server/pkg/stage"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, repo git.Repo, db *sql.DB) {
	v1 := e.Group("/api/v1")

	// TODO: Refactor dependency management
	appRepository := postgres.NewAppRepository(db)
	appService := app.NewService(appRepository)
	secretMetaRepository := postgres.NewSecretMetaRepository(db)
	stageRepository := postgres.NewStageRepository(db)
	stageService := stage.NewService(stageRepository)
	secretService := secret.NewService(appRepository, secretMetaRepository, stageRepository)
	namespaceRepository := postgres.NewNamespaceRepository(db)
	namespaceService := namespace.NewService(namespaceRepository, stageRepository)
	deploymentReservationRepository := postgres.NewDeploymentReservationRepository(db)
	deploymentReservationService := deploymentreservation.NewService(deploymentReservationRepository)
	deploymentRepository := postgres.NewDeploymentRepository(db)
	deploymentService := deployment.NewService(appRepository, namespaceService, secretService, stageRepository, deploymentRepository, deploymentReservationService)
	deploymentStatusService := deploymentstatus.NewService(deploymentRepository, stageService)
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

	// TODO: /apps/:namespace/:appName
	v1.GET("/apps/:appIdOrName", func(c echo.Context) error {
		return GetApp(c, appService)
	})

	// TODO: /apps/:namespace/:appName
	v1.GET("/apps/:appIdOrName/status", func(c echo.Context) error {
		return GetAppStatus(c, appService, deploymentStatusService)
	})

	v1.POST("/apps", func(c echo.Context) error {
		return PostApp(c, appService)
	})

	v1.POST("/deployments", func(c echo.Context) error {
		return PostDeployment(c, repo, appService, deploymentService, stageService)
	})
	v1.PUT("/deployments", func(c echo.Context) error {
		return PostDeployment(c, repo, appService, deploymentService, stageService)
	})

	// TODO(sdk)
	v1.DELETE("/deployments/:stageName/:namespace/:deploymentName", func(c echo.Context) error {
		return DeleteDeployment(c, repo, deploymentService)
	})

	// TODO(sdk)
	v1.PUT("/deployments/:stageName/:namespace/:deploymentName/status", func(c echo.Context) error {
		return PutDeploymentStatus(c, deploymentStatusService)
	})

	// TODO(sdk)
	v1.PUT("/rollout/:stageName/:namespace/:deploymentName", func(c echo.Context) error {
		return PutRollout(c, rolloutService, stageService, repo)
	})

	v1.PUT("/secrets", func(c echo.Context) error {
		return PutSecret(c, repo, secretService, stageService)
	})

	// TODO(sdk)
	v1.GET("/secrets/:appIdOrName/:stageName", func(c echo.Context) error {
		return GetSecrets(c, secretService, stageService)
	})

	v1.GET("/namespaces", func(c echo.Context) error {
		return GetNamespaces(c, namespaceRepository)
	})

	v1.POST("/namespaces", func(c echo.Context) error {
		return PostNamespace(c, namespaceService, repo)
	})

	v1.PUT("/stages/:stageName/config", func(c echo.Context) error {
		return PutStageConfig(c, stageService)
	})

	v1.POST("/stages/:stageName/ping", func(c echo.Context) error {
		return PostStagePing(c, stageService)
	})

	v1.GET("/stages", func(c echo.Context) error {
		return ListStages(c, stageRepository)
	})

	v1.POST("/validate/appconfig", PostValidateAppConfig)
}
