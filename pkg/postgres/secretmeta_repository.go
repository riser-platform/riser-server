package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/pkg/core"
)

type secretMetaRepository struct {
	db *sql.DB
}

func NewSecretMetaRepository(db *sql.DB) core.SecretMetaRepository {
	return &secretMetaRepository{db}
}

func (r *secretMetaRepository) Save(secretMeta *core.SecretMeta) (int64, error) {
	row := r.db.QueryRow(`
		INSERT INTO secret_meta (app_id, stage_name, name, revision) VALUES ($1,$2,$3,0)
		ON CONFLICT(app_id, stage_name, name) DO
		UPDATE SET
			revision=secret_meta.revision + 1
		RETURNING secret_meta.revision
		`, secretMeta.AppId, secretMeta.StageName, secretMeta.Name)

	var revision int64
	err := row.Scan(&revision)
	return revision, err
}

func (r *secretMetaRepository) Commit(secretMeta *core.SecretMeta) error {
	result, err := r.db.Exec(`
	UPDATE secret_meta
		SET committed_revision = revision
		WHERE app_id = $1 AND stage_name = $2 AND name = $3 AND revision = $4
	`, secretMeta.AppId, secretMeta.StageName, secretMeta.Name, secretMeta.Revision)

	if err != nil && !ResultHasRows(result) {
		return core.ErrConflictNewerVersion
	}

	return err
}

func (r *secretMetaRepository) FindByStage(appId uuid.UUID, stageName string) ([]core.SecretMeta, error) {
	secretMetas := []core.SecretMeta{}
	rows, err := r.db.Query(`
	SELECT app_id, stage_name, name, committed_revision
	FROM secret_meta
	WHERE app_id = $1 AND stage_name = $2
	ORDER BY name
	`, appId, stageName)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		secretMeta := core.SecretMeta{}
		err := rows.Scan(&secretMeta.AppId, &secretMeta.StageName, &secretMeta.Name, &secretMeta.Revision)
		if err != nil {
			return nil, err
		}
		secretMetas = append(secretMetas, secretMeta)
	}

	return secretMetas, nil
}
