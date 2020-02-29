package postgres

import (
	"database/sql"

	"github.com/riser-platform/riser-server/pkg/core"
)

type deploymentReservationRepository struct {
	db *sql.DB
}

func NewDeploymentReservationRepository(db *sql.DB) core.DeploymentReservationRepository {
	return &deploymentReservationRepository{db: db}
}

func (r *deploymentReservationRepository) Create(reservation *core.DeploymentReservation) error {
	_, err := r.db.Exec(`INSERT INTO deployment (id, app_id, name, namespace)ReservationS ($1,$2,$3,$4)`,
		reservation.Id, reservation.AppId, reservation.Name, reservation.Namespace)
	return err
}

func (r *deploymentReservationRepository) GetByName(name *core.NamespacedName) (*core.DeploymentReservation, error) {
	reservation := &core.DeploymentReservation{}
	err := r.db.QueryRow(`SELECT id, app_id, name, namespace FROM deployment_reservation WHERE name = $1 AND namespace = $2`,
		name.Name, name.Namespace).Scan(reservation.Id, reservation.AppId, reservation.Name, reservation.Namespace)
	if err == sql.ErrNoRows {
		return nil, core.ErrNotFound
	}

	return reservation, err
}
