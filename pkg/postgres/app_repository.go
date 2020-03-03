package postgres

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/riser-platform/riser-server/pkg/core"
)

type appRepository struct {
	db *sql.DB
}

func NewAppRepository(db *sql.DB) core.AppRepository {
	return &appRepository{db: db}
}

func (r *appRepository) Get(id uuid.UUID) (*core.App, error) {
	app := &core.App{}
	err := r.db.QueryRow("SELECT id, name FROM app WHERE id = $1", id).Scan(&app.Id, &app.Name)
	if err == sql.ErrNoRows {
		return nil, core.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (r *appRepository) GetByName(name *core.NamespacedName) (*core.App, error) {
	app := &core.App{}
	err := r.db.QueryRow("SELECT id, name FROM app WHERE name = $1 and namespace = $2", name.Name, name.Namespace).Scan(&app.Id, &app.Name)
	if err == sql.ErrNoRows {
		return nil, core.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (r *appRepository) Create(app *core.App) error {
	_, err := r.db.Exec("INSERT INTO app (id, name, namespace) VALUES ($1,$2,$3)", app.Id, app.Name, app.Namespace)
	return err
}

func (r *appRepository) ListApps() ([]core.App, error) {
	apps := []core.App{}
	rows, err := r.db.Query("SELECT id, name FROM app")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		app := core.App{}
		err := rows.Scan(&app.Id, &app.Name)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}
