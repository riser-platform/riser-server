package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/pkg/core"
)

type apiKeyRepository struct {
	db *sql.DB
}

func NewApiKeyRepository(db *sql.DB) core.ApiKeyRepository {
	return &apiKeyRepository{db}
}

func (r *apiKeyRepository) GetByUserId(userId uuid.UUID) ([]core.ApiKey, error) {
	apiKeys := []core.ApiKey{}
	rows, err := r.db.Query(`
	SELECT riser_user_id, key_hash
	FROM apikey
	WHERE riser_user_id = $1
	`, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		apiKey := core.ApiKey{}
		err := rows.Scan(&apiKey.UserId, &apiKey.KeyHash)
		if err != nil {
			return nil, err
		}
		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys, nil
}

func (r *apiKeyRepository) Create(userId uuid.UUID, keyHash []byte) error {
	_, err := r.db.Exec("INSERT INTO apikey (riser_user_id, key_hash) VALUES ($1, $2)", userId, keyHash)
	return err
}
