package postgres

import (
	"database/sql"

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
		app_id := app_id FROM app WHERE name = $1 AND namespace = $2
		INSERT INTO secret_meta (app_id, stage_name, name, revision) VALUES (app_id,$3,$4,0)
		ON CONFLICT(app_id, stage_name, name) DO
		UPDATE SET
			revision=secret_meta.revision + 1
		RETURNING secret_meta.revision
		`, secretMeta.App.Name, secretMeta.App.Namespace, secretMeta.StageName, secretMeta.Name)

	var revision int64
	err := row.Scan(&revision)
	return revision, err
}

func (r *secretMetaRepository) Commit(secretMeta *core.SecretMeta) error {

	result, err := r.db.Exec(`
		app_id := app_id FROM app WHERE name = $1 AND namespace = $2
		UPDATE secret_meta
		SET committed_revision = revision
		WHERE app_id = app_id AND stage_name = $3 AND name = $4 AND revision = $5
	`, secretMeta.App.Name, secretMeta.App.Namespace, secretMeta.StageName, secretMeta.Name, secretMeta.Revision)

	if err != nil && !ResultHasRows(result) {
		return core.ErrConflictNewerVersion
	}

	return err
}

func (r *secretMetaRepository) ListByAppInStage(appName *core.NamespacedName, stageName string) ([]core.SecretMeta, error) {
	secretMetas := []core.SecretMeta{}
	rows, err := r.db.Query(`
	SELECT
		app.name,
		app.namespace,
		secret_meta.stage_name,
		secret_meta.name,
		secret_meta.committed_revision
	FROM secret_meta
	INNER JOIN app ON app.id = secret_meta.app_id
	WHERE app.name = $1 AND app.namespace = $2 AND stage_name = $2
	ORDER BY name
	`, appName.Name, appName.Namespace, stageName)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		secretMeta := core.SecretMeta{}
		err := rows.Scan(&secretMeta.App.Name, &secretMeta.App.Namespace, &secretMeta.StageName, &secretMeta.Name, &secretMeta.Revision)
		if err != nil {
			return nil, err
		}
		secretMetas = append(secretMetas, secretMeta)
	}

	return secretMetas, nil
}
