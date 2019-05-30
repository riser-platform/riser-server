package postgres

import (
	"database/sql"

	"github.com/riser-platform/riser-server/pkg/core"
)

type stageRepository struct {
	db *sql.DB
}

func NewStageRepository(db *sql.DB) core.StageRepository {
	return &stageRepository{db}
}

func (r *stageRepository) Get(name string) (*core.Stage, error) {
	stage := &core.Stage{}
	err := r.db.QueryRow("SELECT name, doc FROM stage WHERE name = $1", name).Scan(&stage.Name, &stage.Doc)
	if err == sql.ErrNoRows {
		return nil, core.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return stage, nil
}

func (r *stageRepository) List() ([]core.Stage, error) {
	stages := []core.Stage{}
	rows, err := r.db.Query("SELECT name, doc FROM stage")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		stage := core.Stage{}
		err := rows.Scan(&stage.Name, &stage.Doc)
		if err != nil {
			return nil, err
		}
		stages = append(stages, stage)
	}

	return stages, nil
}

func (r *stageRepository) Save(stage *core.Stage) error {
	_, err := r.db.Exec(`
		INSERT INTO stage(name, doc) VALUES($1,$2)
		ON CONFLICT (name) DO
		UPDATE SET
			doc = $2;`, stage.Name, &stage.Doc)

	return err
}
