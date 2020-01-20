package postgres

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

const MigrationsUrl = "file://migrations"

func Migrate(postgresConn string, logger *logrus.Logger) error {
	m, err := migrate.New(MigrationsUrl, postgresConn)
	if err != nil {
		return err
	}
	m.Log = &logrusMigrateAdapter{logger}
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("No new postgres migrations found")
		}
	} else {
		return err
	}

	return nil
}

type logrusMigrateAdapter struct {
	logger *logrus.Logger
}

func (a *logrusMigrateAdapter) Printf(format string, v ...interface{}) {
	// TODO: No logrus categories? Maybe switch loggers...
	a.logger.WithField("category", "postgresmigration").Printf(format, v...)
}

func (a *logrusMigrateAdapter) Verbose() bool {
	return a.logger.IsLevelEnabled(logrus.DebugLevel)
}
