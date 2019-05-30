package main

import (
	"database/sql"

	"github.com/riser-platform/riser-server/pkg/login"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/postgres"

	"os"

	apiv1 "github.com/riser-platform/riser-server/api/v1"
	"github.com/riser-platform/riser-server/pkg/git"

	"github.com/joho/godotenv"

	"github.com/kelseyhightower/envconfig"

	"github.com/labstack/echo/v4"
	echolog "github.com/onrik/logrus/echo"
	"github.com/sirupsen/logrus"
)

// All env vars are prefixed with RISER_
const envPrefix = "RISER"

// DotEnv file typically used For local development
const dotEnvFile = ".env"

var logger = logrus.StandardLogger()

func main() {
	logger.SetFormatter(&logrus.JSONFormatter{})

	err := loadDotEnv()
	exitIfError(err, "Error loading .env file")

	var rc core.RuntimeConfig
	err = envconfig.Process(envPrefix, &rc)
	exitIfError(err, "Error loading environment variables")

	if rc.DeveloperMode {
		logger.SetFormatter(&logrus.TextFormatter{})
		logger.Info("Developer mode active")
	}

	postgresConn, err := postgres.AddAuthToConnString(rc.PostgresUrl, rc.PostgresUsername, rc.PostgresPassword)
	exitIfError(err, "Error creating postgres connection url")

	postgresDb, err := postgres.NewDB(postgresConn)
	exitIfError(err, "Unable to initialize postgres")

	if rc.PostgresMigrateOnStartup {
		logger.Info("Applying Postgres migrations")
		err = postgres.Migrate(postgresConn, logger)
		exitIfError(err, "Error performing Postgres migrations")
	}

	bootstrapApiKey(postgresDb, &rc)

	repoSettings := git.RepoSettings{
		URL: rc.GitUrl,
		// Note: The underlying git code has multiple bugs related to non-master branches.
		Branch:      "master",
		Username:    rc.GitUsername,
		Password:    rc.GitPassword,
		LocalGitDir: rc.GitDir,
	}

	repo, err := git.NewGitRepoProvider(repoSettings)
	exitIfError(err, "Error initializing git")

	e := echo.New()
	e.HideBanner = true

	e.Logger = echolog.NewLogger(logger, "")
	e.Use(echolog.Middleware(echolog.DefaultConfig))
	e.HTTPErrorHandler = apiv1.ApiErrorHandler

	apiv1.RegisterRoutes(e, repo, postgresDb)
	err = e.Start(rc.BindAddress)
	exitIfError(err, "Error starting server")
}

func bootstrapApiKey(db *sql.DB, rc *core.RuntimeConfig) {
	loginService := login.NewService(postgres.NewUserRepository(db), postgres.NewApiKeyRepository(db))
	err := loginService.BootstrapRootUser(rc.BootstrapApikey)
	if err != nil {
		if err == login.ErrRootUserExists {
			logger.Info("Ignoring environment variable RISER_BOOTSTRAP_APIKEY: root user already exists.")
		} else {
			exitIfError(err, "Unable to bootstrap API KEY")
		}
	}
}

func loadDotEnv() error {
	_, err := os.Stat(dotEnvFile)
	if !os.IsNotExist(err) {
		return godotenv.Load(dotEnvFile)
	}

	return nil
}

func exitIfError(err error, message string) {
	if err != nil {
		logger.Fatalf("%s: %s", message, err)
	}
}
