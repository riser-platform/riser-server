package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
)

type deploymentRepository struct {
	db *sql.DB
}

func NewDeploymentRepository(db *sql.DB) core.DeploymentRepository {
	return &deploymentRepository{db: db}
}

func (r *deploymentRepository) Create(deployment *core.DeploymentRecord) error {
	_, err := r.db.Exec(`
	INSERT INTO deployment (id, deployment_reservation_id, stage_name, riser_revision, doc)
	VALUES ($1,$2,$3,$4,$5)`,
		deployment.Id, deployment.ReservationId, deployment.StageName, deployment.RiserRevision, &deployment.Doc)
	return err
}

func (r *deploymentRepository) Delete(name *core.NamespacedName, stageName string) error {
	_, err := r.db.Exec(`
	UPDATE deployment SET deleted_at=now()
	INNER JOIN deployment_reservation ON deployment.deployment_reservation_id = deployment_reservation.id
	WHERE deployment_reservation.name = $1 deployment_reservation.namespace = $2 AND deployment.stage_name = $2
	`, name.Name, name.Namespace, stageName)

	return noRowsErrorHandler(err)
}

// GetByName returns a deployment by its name whether or not it's been deleted.
func (r *deploymentRepository) GetByName(name *core.NamespacedName, stageName string) (*core.Deployment, error) {
	deployment := &core.Deployment{}
	err := r.db.QueryRow(`
	SELECT
		deployment_reservation.id,
		deployment_reservation.app_id,
		deployment_reservation.name,
		deployment_reservation.namespace,
		deployment.id,
		deployment.deleted_at,
		deployment.deployment_reservation_id,
		deployment.stage_name,
		deployment.riser_revision,
		deployment.doc
	FROM deployment
	INNER JOIN deployment_reservation ON deployment.deployment_reservation_id=deployment_reservation.id
	WHERE deployment_reservation.name=$1 AND deployment_reservation.namespace=$2 AND deployment.stage_name=$3
	`, name.Name, name.Namespace, stageName).Scan(
		&deployment.DeploymentReservation.Id,
		&deployment.AppId,
		&deployment.Name,
		&deployment.Namespace,
		&deployment.DeploymentRecord.Id,
		&deployment.DeletedAt,
		&deployment.ReservationId,
		&deployment.StageName,
		&deployment.RiserRevision,
		&deployment.Doc)

	return deployment, noRowsErrorHandler(err)
}

// GetByReservation returns a deployment by its reservation ID whether or not it has been deleted.
func (r *deploymentRepository) GetByReservation(reservationId uuid.UUID, stageName string) (*core.Deployment, error) {
	deployment := &core.Deployment{}

	err := r.db.QueryRow(`
	SELECT
		deployment_reservation.id,
		deployment_reservation.app_id,
		deployment_reservation.name,
		deployment_reservation.namespace,
		deployment.id,
		deployment.deleted_at,
		deployment.deployment_reservation_id,
		deployment.stage_name,
		deployment.riser_revision,
		deployment.doc
	FROM deployment
	INNER JOIN deployment_reservation ON deployment.deployment_reservation_id = deployment_reservation.id
	WHERE deployment_reservation.id = $1 AND deployment.stage_name = $2
	`, reservationId, stageName).Scan(
		&deployment.DeploymentReservation.Id,
		&deployment.AppId,
		&deployment.Name,
		&deployment.Namespace,
		&deployment.DeploymentRecord.Id,
		&deployment.DeletedAt,
		&deployment.ReservationId,
		&deployment.StageName,
		&deployment.RiserRevision,
		&deployment.Doc)
	return deployment, noRowsErrorHandler(err)
}

// FindByApp returns all active deployments in all stages by a given app ID
func (r *deploymentRepository) FindByApp(appId uuid.UUID) ([]core.Deployment, error) {
	deployments := []core.Deployment{}
	rows, err := r.db.Query(`
	SELECT
		deployment_reservation.id,
		deployment_reservation.app_id,
		deployment_reservation.name,
		deployment_reservation.namespace,
		deployment.id,
		deployment.deleted_at,
		deployment.deployment_reservation_id,
		deployment.stage_name,
		deployment.riser_revision,
		deployment.doc
	FROM deployment
	INNER JOIN deployment_reservation ON deployment.deployment_reservation_id = deployment_reservation.id
	WHERE deployment_reservation.app_id = $1 AND deployment.deleted_at IS NULL
	ORDER BY deployment.stage_name, deployment_reservation.name
	`, appId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		deployment := core.Deployment{}
		err := rows.Scan(
			&deployment.DeploymentReservation.Id,
			&deployment.AppId,
			&deployment.Name,
			&deployment.Namespace,
			&deployment.DeploymentRecord.Id,
			&deployment.DeletedAt,
			&deployment.ReservationId,
			&deployment.StageName,
			&deployment.RiserRevision,
			&deployment.Doc)
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, deployment)
	}

	return deployments, nil
}

// IncrementeRevision increments the revision of a deployment. If the deployment was previously soft deleted, it will mark
// the deployment as no longer being deleted
func (r *deploymentRepository) IncrementRevision(name *core.NamespacedName, stageName string) (revision int64, err error) {
	err = r.db.QueryRow(`
	UPDATE deployment SET riser_revision = riser_revision + 1, deleted_at = NULL
	FROM deployment_reservation
	WHERE
	deployment.deployment_reservation_id = deployment_reservation.id
	AND deployment_reservation.name = $1
	AND deployment_reservation.namespace = $2
	AND stage_name = $3
	RETURNING riser_revision
	`, name.Name, name.Namespace, stageName).Scan(&revision)
	if err != nil {
		return 0, err
	}

	return revision, nil
}

func (r *deploymentRepository) RollbackRevision(name *core.NamespacedName, stageName string, failedRevision int64) (revision int64, err error) {
	err = r.db.QueryRow(`
	UPDATE deployment SET riser_revision = riser_revision - 1
	FROM deployment_reservation
	WHERE
	deployment.deployment_reservation_id = deployment_reservation.id
	AND deployment_reservation.name = $1
	AND deployment_reservation.namespace = $2
	AND stage_name = $3
	AND riser_revision = $4
	RETURNING riser_revision
	`, name.Name, name.Namespace, stageName, failedRevision).Scan(&revision)
	if err != nil {
		return 0, err
	}

	return revision, nil
}

func (r *deploymentRepository) UpdateStatus(name *core.NamespacedName, stageName string, status *core.DeploymentStatus) error {
	result, err := r.db.Exec(`
	  UPDATE deployment
		SET doc = jsonb_set(doc, '{status}', $4),
		-- If we receive a status update, we "undelete" the deployment
		deleted_at = null
		FROM deployment_reservation
		WHERE
		deployment.deployment_reservation_id = deployment_reservation.id
		AND deployment_reservation.name = $1
		AND deployment_reservation.namespace = $2
		AND deployment.stage_name = $3
		-- Don't update status from an older observed revision
		AND (
			(deployment.doc->'status'->>'observedRiserRevision')::int <= $5
			OR deployment.doc->'status' IS NULL
			OR deployment.doc->'status'->>'observedRiserRevision' IS NULL
		)
	`, name.Name, name.Namespace, stageName, status, status.ObservedRiserRevision)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("Deployment not found or status is outdated")
	}

	return nil
}

func (r *deploymentRepository) UpdateTraffic(name *core.NamespacedName, stageName string, riserRevision int64, traffic core.TrafficConfig) error {
	result, err := r.db.Exec(`
		UPDATE deployment
		SET doc = jsonb_set(doc, '{traffic}', $5)
		FROM deployment_reservation
		WHERE
		deployment.deployment_reservation_id = deployment_reservation.id
		AND deployment_reservation.name = $1
		AND deployment_reservation.namespace = $2
		AND deployment.stage_name = $3
		AND riser_revision = $4
		AND deleted_at IS NULL
	`, name.Name, name.Namespace, stageName, riserRevision, traffic)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("Deployment not found or has been updated by another process")
	}

	return nil
}
