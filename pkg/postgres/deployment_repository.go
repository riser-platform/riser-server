package postgres

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
)

type deploymentRepository struct {
	db *sql.DB
}

func NewDeploymentRepository(db *sql.DB) core.DeploymentRepository {
	return &deploymentRepository{db: db}
}

func (r *deploymentRepository) Create(newDeployment *core.Deployment) error {
	_, err := r.db.Exec(`INSERT INTO deployment (name, stage_name, app_name, riser_generation, doc) VALUES ($1,$2,$3,$4,$5)`,
		newDeployment.Name, newDeployment.StageName, newDeployment.AppName, newDeployment.RiserGeneration, &newDeployment.Doc)
	return err
}

func (r *deploymentRepository) Get(deploymentName, stageName string) (*core.Deployment, error) {
	deployment := &core.Deployment{}
	err := r.db.QueryRow(`
	SELECT name, stage_name, app_name, riser_generation, doc FROM deployment
	WHERE name = $1 AND stage_name = $2
	`, deploymentName, stageName).Scan(&deployment.Name, &deployment.StageName, &deployment.AppName, &deployment.RiserGeneration, &deployment.Doc)
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

func (r *deploymentRepository) FindByApp(appName string) ([]core.Deployment, error) {
	deployments := []core.Deployment{}
	rows, err := r.db.Query(`
	SELECT name, stage_name, app_name, riser_generation, doc
	FROM deployment
	WHERE app_name = $1
	ORDER BY stage_name, name
	`, appName)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		deployment := core.Deployment{}
		err := rows.Scan(&deployment.Name, &deployment.StageName, &deployment.AppName, &deployment.RiserGeneration, &deployment.Doc)
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, deployment)
	}

	return deployments, nil
}

func (r *deploymentRepository) IncrementGeneration(deploymentName, stageName string) (generation int64, err error) {
	err = r.db.QueryRow(`
	UPDATE deployment SET riser_generation = riser_generation + 1
	WHERE name = $1 AND stage_name = $2
	RETURNING riser_generation
	`, deploymentName, stageName).Scan(&generation)
	if err != nil {
		return 0, err
	}

	return generation, nil
}

func (r *deploymentRepository) RollbackGeneration(deploymentName, stageName string, failedGeneration int64) (generation int64, err error) {
	err = r.db.QueryRow(`
	UPDATE deployment SET riser_generation = riser_generation - 1
	WHERE name = $1 AND stage_name = $2 AND riser_generation = $3
	RETURNING riser_generation
	`, deploymentName, stageName, failedGeneration).Scan(&generation)
	if err != nil {
		return 0, err
	}

	return generation, nil
}

func (r *deploymentRepository) UpdateStatus(deploymentName, stageName string, status *core.DeploymentStatus) error {
	result, err := r.db.Exec(`
	  UPDATE deployment
		SET doc = jsonb_set(doc, '{status}', $3)
		WHERE name = $1 AND stage_name = $2
		-- Don't update status from an older observed generation
		AND (doc->'status'->>'observedRiserGeneration')::int <= $4
	`, deploymentName, stageName, status, status.ObservedRiserGeneration)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("Deployment not found")
	}

	return nil
}
