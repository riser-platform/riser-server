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

func (r *secretMetaRepository) Save(secretMeta *core.SecretMeta) error {
	_, err := r.db.Exec(`
		INSERT INTO secret_meta (app_name, stage_name, secret_name, doc) VALUES ($1,$2,$3,$4)
		ON CONFLICT(app_name, stage_name, secret_name) DO
		UPDATE SET
			doc=$4;
		`, secretMeta.AppName, secretMeta.StageName, secretMeta.SecretName, &secretMeta.Doc)

	return err
}

func (r *secretMetaRepository) FindByStage(appName string, stageName string) ([]core.SecretMeta, error) {
	secretMetas := []core.SecretMeta{}
	rows, err := r.db.Query(`
	SELECT app_name, stage_name, secret_name, doc
	FROM secret_meta
	WHERE app_name = $1 AND stage_name = $2
	ORDER BY secret_name
	`, appName, stageName)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		secretMeta := core.SecretMeta{}
		err := rows.Scan(&secretMeta.AppName, &secretMeta.StageName, &secretMeta.SecretName, &secretMeta.Doc)
		if err != nil {
			return nil, err
		}
		secretMetas = append(secretMetas, secretMeta)
	}

	return secretMetas, nil
}
