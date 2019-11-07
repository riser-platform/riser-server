package postgres

import (
	"database/sql"

	"github.com/riser-platform/riser-server/pkg/core"
)

type appRepository struct {
	db *sql.DB
}

func NewAppRepository(db *sql.DB) core.AppRepository {
	return &appRepository{db: db}
}

func (r *appRepository) Get(name string) (*core.App, error) {
	app := &core.App{}
	err := r.db.QueryRow("SELECT name, hashid FROM app WHERE name = $1", name).Scan(&app.Name, &app.Hashid)
	if err == sql.ErrNoRows {
		return nil, core.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (r *appRepository) Create(app *core.App) error {
	_, err := r.db.Exec("INSERT INTO app (name, hashid) VALUES ($1,$2)", app.Name, app.Hashid)
	return err
}

func (r *appRepository) ListApps() ([]core.App, error) {
	apps := []core.App{}
	rows, err := r.db.Query("SELECT name, hashid FROM app")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		app := core.App{}
		err := rows.Scan(&app.Name, &app.Hashid)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}
