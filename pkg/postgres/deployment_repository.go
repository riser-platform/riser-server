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
	INSERT INTO deployment (id, deployment_reservation_id, environment_name, riser_revision, doc)
	VALUES ($1,$2,$3,$4,$5)`,
		deployment.Id, deployment.ReservationId, deployment.EnvironmentName, deployment.RiserRevision, &deployment.Doc)
	return err
}

func (r *deploymentRepository) Delete(name *core.NamespacedName, envName string) error {
	_, err := r.db.Exec(`
	UPDATE deployment SET deleted_at=now()
	FROM deployment_reservation
	WHERE
	 deployment.deployment_reservation_id = deployment_reservation.id
	 AND deployment_reservation.name = $1
	 AND deployment_reservation.namespace = $2
	 AND deployment.environment_name = $3
	`, name.Name, name.Namespace, envName)

	return noRowsErrorHandler(err)
}

// GetByName returns a deployment by its name whether or not it's been deleted.
func (r *deploymentRepository) GetByName(name *core.NamespacedName, envName string) (*core.Deployment, error) {
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
		deployment.environment_name,
		deployment.riser_revision,
		deployment.doc
	FROM deployment
	INNER JOIN deployment_reservation ON deployment.deployment_reservation_id=deployment_reservation.id
	WHERE deployment_reservation.name=$1 AND deployment_reservation.namespace=$2 AND deployment.environment_name=$3
	`, name.Name, name.Namespace, envName).Scan(
		&deployment.DeploymentReservation.Id,
		&deployment.AppId,
		&deployment.Name,
		&deployment.Namespace,
		&deployment.DeploymentRecord.Id,
		&deployment.DeletedAt,
		&deployment.ReservationId,
		&deployment.EnvironmentName,
		&deployment.RiserRevision,
		&deployment.Doc)

	return deployment, noRowsErrorHandler(err)
}

// GetByReservation returns a deployment by its reservation ID whether or not it has been deleted.
func (r *deploymentRepository) GetByReservation(reservationId uuid.UUID, envName string) (*core.Deployment, error) {
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
		deployment.environment_name,
		deployment.riser_revision,
		deployment.doc
	FROM deployment
	INNER JOIN deployment_reservation ON deployment.deployment_reservation_id = deployment_reservation.id
	WHERE deployment_reservation.id = $1 AND deployment.environment_name = $2
	`, reservationId, envName).Scan(
		&deployment.DeploymentReservation.Id,
		&deployment.AppId,
		&deployment.Name,
		&deployment.Namespace,
		&deployment.DeploymentRecord.Id,
		&deployment.DeletedAt,
		&deployment.ReservationId,
		&deployment.EnvironmentName,
		&deployment.RiserRevision,
		&deployment.Doc)
	return deployment, noRowsErrorHandler(err)
}

// FindByApp returns all active deployments in all environments by a given app ID
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
		deployment.environment_name,
		deployment.riser_revision,
		deployment.doc
	FROM deployment
	INNER JOIN deployment_reservation ON deployment.deployment_reservation_id = deployment_reservation.id
	WHERE deployment_reservation.app_id = $1 AND deployment.deleted_at IS NULL
	ORDER BY deployment.environment_name, deployment_reservation.name
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
			&deployment.EnvironmentName,
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
func (r *deploymentRepository) IncrementRevision(name *core.NamespacedName, envName string) (revision int64, err error) {
	err = r.db.QueryRow(`
	UPDATE deployment SET riser_revision = riser_revision + 1, deleted_at = NULL
	FROM deployment_reservation
	WHERE
	deployment.deployment_reservation_id = deployment_reservation.id
	AND deployment_reservation.name = $1
	AND deployment_reservation.namespace = $2
	AND environment_name = $3
	RETURNING riser_revision
	`, name.Name, name.Namespace, envName).Scan(&revision)
	if err != nil {
		return 0, err
	}

	return revision, nil
}

func (r *deploymentRepository) RollbackRevision(name *core.NamespacedName, envName string, failedRevision int64) (revision int64, err error) {
	err = r.db.QueryRow(`
	UPDATE deployment SET riser_revision = riser_revision - 1
	FROM deployment_reservation
	WHERE
	deployment.deployment_reservation_id = deployment_reservation.id
	AND deployment_reservation.name = $1
	AND deployment_reservation.namespace = $2
	AND environment_name = $3
	AND riser_revision = $4
	RETURNING riser_revision
	`, name.Name, name.Namespace, envName, failedRevision).Scan(&revision)
	if err != nil {
		return 0, err
	}

	return revision, nil
}

func (r *deploymentRepository) UpdateStatus(name *core.NamespacedName, envName string, status *core.DeploymentStatus) error {
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
		AND deployment.environment_name = $3
		-- Don't update status from an older observed revision
		AND (
			(deployment.doc->'status'->>'observedRiserRevision')::int <= $5
			OR deployment.doc->'status' IS NULL
			OR deployment.doc->'status'->>'observedRiserRevision' IS NULL
		)
	`, name.Name, name.Namespace, envName, status, status.ObservedRiserRevision)

	if err != nil {
		return err
	}
	return r.handleUpdateStatusResult(result)
}

func (deploymentRepository) handleUpdateStatusResult(r sql.Result) error {
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}

	/*
		Assume that no updates = a version conflict.
		While this assumption is not strictly true, this is the most common reason that the update predicate fails to match
	*/
	if rows == 0 {
		return core.ErrConflictNewerVersion
	}

	return nil
}

func (r *deploymentRepository) UpdateTraffic(name *core.NamespacedName, envName string, riserRevision int64, traffic core.TrafficConfig) error {
	result, err := r.db.Exec(`
		UPDATE deployment
		SET doc = jsonb_set(doc, '{traffic}', $5)
		FROM deployment_reservation
		WHERE
		deployment.deployment_reservation_id = deployment_reservation.id
		AND deployment_reservation.name = $1
		AND deployment_reservation.namespace = $2
		AND deployment.environment_name = $3
		AND riser_revision = $4
		AND deleted_at IS NULL
	`, name.Name, name.Namespace, envName, riserRevision, traffic)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("Deployment not found or has been updated by another process")
	}

	return nil
}
