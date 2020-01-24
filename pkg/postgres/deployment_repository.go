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
	_, err := r.db.Exec(`INSERT INTO deployment (name, stage_name, app_name, riser_revision, doc) VALUES ($1,$2,$3,$4,$5)`,
		newDeployment.Name, newDeployment.StageName, newDeployment.AppName, newDeployment.RiserRevision, &newDeployment.Doc)
	return err
}

func (r *deploymentRepository) Get(deploymentName, stageName string) (*core.Deployment, error) {
	deployment := &core.Deployment{}
	err := r.db.QueryRow(`
	SELECT name, stage_name, app_name, riser_revision, doc FROM deployment
	WHERE name = $1 AND stage_name = $2
	`, deploymentName, stageName).Scan(&deployment.Name, &deployment.StageName, &deployment.AppName, &deployment.RiserRevision, &deployment.Doc)
	if err != nil {
		if err == sql.ErrNoRows {
			err = core.ErrNotFound
		}
		return nil, err
	}
	return deployment, nil
}

func (r *deploymentRepository) FindByApp(appName string) ([]core.Deployment, error) {
	deployments := []core.Deployment{}
	rows, err := r.db.Query(`
	SELECT name, stage_name, app_name, riser_revision, doc
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
		err := rows.Scan(&deployment.Name, &deployment.StageName, &deployment.AppName, &deployment.RiserRevision, &deployment.Doc)
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, deployment)
	}

	return deployments, nil
}

func (r *deploymentRepository) IncrementRevision(deploymentName, stageName string) (revision int64, err error) {
	err = r.db.QueryRow(`
	UPDATE deployment SET riser_revision = riser_revision + 1
	WHERE name = $1 AND stage_name = $2
	RETURNING riser_revision
	`, deploymentName, stageName).Scan(&revision)
	if err != nil {
		return 0, err
	}

	return revision, nil
}

func (r *deploymentRepository) RollbackRevision(deploymentName, stageName string, failedRevision int64) (revision int64, err error) {
	err = r.db.QueryRow(`
	UPDATE deployment SET riser_revision = riser_revision - 1
	WHERE name = $1 AND stage_name = $2 AND riser_revision = $3
	RETURNING riser_revision
	`, deploymentName, stageName, failedRevision).Scan(&revision)
	if err != nil {
		return 0, err
	}

	return revision, nil
}

func (r *deploymentRepository) UpdateStatus(deploymentName, stageName string, status *core.DeploymentStatus) error {
	result, err := r.db.Exec(`
	  UPDATE deployment
		SET doc = jsonb_set(doc, '{status}', $3)
		WHERE name = $1 AND stage_name = $2
		-- Don't update status from an older observed revision
		AND ((doc->'status'->>'observedRiserRevision')::int <= $4 OR doc->'status' IS NULL OR doc->'status'->>'observedRiserRevision' IS NULL)
	`, deploymentName, stageName, status, status.ObservedRiserRevision)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("Deployment not found or status is outdated")
	}

	return nil
}

func (r *deploymentRepository) UpdateTraffic(deploymentName, stageName string, riserRevision int64, traffic core.TrafficConfig) error {
	result, err := r.db.Exec(`
		UPDATE DEPLOYMENT
		SET doc = jsonb_set(doc, '{traffic}', $4)
		WHERE Name = $1 AND stage_name = $2 AND riser_revision = $3
	`, deploymentName, stageName, riserRevision, traffic)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("Deployment not found or has been updated by another process")
	}

	return nil
}
