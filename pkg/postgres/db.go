package postgres

import (
	"database/sql"
	"net/url"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"

	// Required for pq lib dynamic driver loading
	_ "github.com/lib/pq"
)

func NewDB(postgresConn string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", postgresConn)
	if err != nil {
		return nil, errors.Wrap(err, "error opening connection to posgres")
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "error pinging postgres")
	}

	return db, nil
}

func AddAuthToConnString(postgresConn string, username string, password string) (string, error) {
	postgresUrl, err := url.Parse(postgresConn)
	if err != nil {
		return "", err
	}
	postgresUrl.User = url.UserPassword(username, password)
	return postgresUrl.String(), nil
}

func resultHasRows(r sql.Result) bool {
	if r == nil {
		return false
	}
	rows, err := r.RowsAffected()
	return err == nil && rows > 0
}

func noRowsErrorHandler(err error) error {
	if err == sql.ErrNoRows {
		return core.ErrNotFound
	}

	return errors.WithStack(err)
}
