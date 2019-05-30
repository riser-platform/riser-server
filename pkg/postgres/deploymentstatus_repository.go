package postgres

import (
	"database/sql"

	"github.com/riser-platform/riser-server/pkg/core"
)

type deploymentStatusRepository struct {
	db *sql.DB
}

func NewDeploymentStatusRepository(db *sql.DB) core.DeploymentStatusRepository {
	return &deploymentStatusRepository{db: db}
}

func (r *deploymentStatusRepository) FindByApp(appName string) ([]core.DeploymentStatus, error) {
	statuses := []core.DeploymentStatus{}
	rows, err := r.db.Query(`
	SELECT app_name, deployment_name, stage_name, doc
	FROM deployment_status
	WHERE app_name = $1
	ORDER BY stage_name, deployment_name
	`, appName)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		status := core.DeploymentStatus{}
		err := rows.Scan(&status.AppName, &status.DeploymentName, &status.StageName, &status.Doc)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (r *deploymentStatusRepository) Save(status *core.DeploymentStatus) error {
	_, err := r.db.Exec(`
	  INSERT INTO deployment_status(app_name, deployment_name, stage_name, doc) VALUES($1,$2,$3,$4)
		ON CONFLICT (app_name, deployment_name, stage_name) DO
		UPDATE SET
			doc=$4;

	`, status.AppName, status.DeploymentName, status.StageName, status.Doc)

	return err
}
