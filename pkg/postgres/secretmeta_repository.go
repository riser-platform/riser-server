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
		INSERT INTO secret_meta (app_id, environment_name, name, revision)
		SELECT app.id, $3, $4, 1
		FROM app
		WHERE
			app.name = $1
			AND app.namespace = $2
		ON CONFLICT(app_id, environment_name, name) DO
		UPDATE SET
			revision=secret_meta.revision + 1
		RETURNING secret_meta.revision
		`, secretMeta.App.Name, secretMeta.App.Namespace, secretMeta.EnvironmentName, secretMeta.Name)

	var revision int64
	err := row.Scan(&revision)
	return revision, err
}

func (r *secretMetaRepository) Commit(secretMeta *core.SecretMeta) error {
	result, err := r.db.Exec(`
		UPDATE secret_meta
		SET committed_revision = revision
		FROM app
		WHERE
			secret_meta.app_id = app.id
			AND app.name = $1
			AND app.namespace = $2
			AND secret_meta.environment_name = $3
			AND secret_meta.name = $4
			AND secret_meta.revision = $5
	`, secretMeta.App.Name, secretMeta.App.Namespace, secretMeta.EnvironmentName, secretMeta.Name, secretMeta.Revision)

	if err != nil && !resultHasRows(result) {
		return core.ErrConflictNewerVersion
	}

	return err
}

func (r *secretMetaRepository) ListByAppInEnvironment(appName *core.NamespacedName, envName string) ([]core.SecretMeta, error) {
	secretMetas := []core.SecretMeta{}
	rows, err := r.db.Query(`
	SELECT
		app.name,
		app.namespace,
		secret_meta.environment_name,
		secret_meta.name,
		secret_meta.committed_revision
	FROM secret_meta
	INNER JOIN app ON app.id = secret_meta.app_id
	WHERE
		app.name = $1
		AND app.namespace = $2
		AND secret_meta.environment_name = $3
	ORDER BY secret_meta.name
	`, appName.Name, appName.Namespace, envName)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		secretMeta := core.SecretMeta{App: &core.NamespacedName{}}
		err := rows.Scan(&secretMeta.App.Name, &secretMeta.App.Namespace, &secretMeta.EnvironmentName, &secretMeta.Name, &secretMeta.Revision)
		if err != nil {
			return nil, err
		}
		secretMetas = append(secretMetas, secretMeta)
	}

	return secretMetas, nil
}
