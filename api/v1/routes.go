package v1

import (
	"database/sql"

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

func RegisterRoutes(e *echo.Echo, repo git.GitRepoProvider, db *sql.DB) {
	v1 := e.Group("/api/v1")

	// TODO: Refactor dependency management
	appRepository := postgres.NewAppRepository(db)
	appService := app.NewService(appRepository)
	secretMetaRepository := postgres.NewSecretMetaRepository(db)
	stageRepository := postgres.NewStageRepository(db)
	stageService := stage.NewService(stageRepository)
	secretService := secret.NewService(secretMetaRepository, stageRepository)
	deploymentService := deployment.NewService(secretService, stageRepository)
	deploymentStatusRepository := postgres.NewDeploymentStatusRepository(db)
	deploymentStatusService := deploymentstatus.NewService(deploymentStatusRepository, stageService)
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

	// NOTE: May just limit to /apps/:namespace in the future
	v1.GET("/apps", func(c echo.Context) error {
		return ListApps(c, appRepository)
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

	v1.PUT("/secrets", func(c echo.Context) error {
		return PutSecret(c, repo, secretService, stageService)
	})
	v1.GET("/secrets/:appName/:stageName", func(c echo.Context) error {
		return GetSecrets(c, secretService, stageService)
	})

	v1.POST("/status", func(c echo.Context) error {
		return PostStatus(c, deploymentStatusService)
	})
	v1.PUT("/status", func(c echo.Context) error {
		return PostStatus(c, deploymentStatusService)
	})
	v1.GET("/status", func(c echo.Context) error {
		return GetStatus(c, deploymentStatusService)
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
