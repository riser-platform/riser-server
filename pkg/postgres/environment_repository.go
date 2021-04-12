package postgres

import (
	"database/sql"

	"github.com/riser-platform/riser-server/pkg/core"
)

type environmentRepository struct {
	db *sql.DB
}

func NewEnvironmentRepository(db *sql.DB) core.EnvironmentRepository {
	return &environmentRepository{db}
}

func (r *environmentRepository) Get(name string) (*core.Environment, error) {
	environment := &core.Environment{}
	err := r.db.QueryRow("SELECT name, doc FROM environment WHERE name = $1", name).Scan(&environment.Name, &environment.Doc)
	if err == sql.ErrNoRows {
		return nil, core.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return environment, nil
}

func (r *environmentRepository) List() ([]core.Environment, error) {
	environments := []core.Environment{}
	rows, err := r.db.Query("SELECT name, doc FROM environment")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		environment := core.Environment{}
		err := rows.Scan(&environment.Name, &environment.Doc)
		if err != nil {
			return nil, err
		}
		environments = append(environments, environment)
	}

	return environments, nil
}

func (r *environmentRepository) Save(environment *core.Environment) error {
	_, err := r.db.Exec(`
		INSERT INTO environment(name, doc) VALUES($1,$2)
		ON CONFLICT (name) DO
		UPDATE SET
			doc = $2;`, environment.Name, &environment.Doc)

	return err
}
