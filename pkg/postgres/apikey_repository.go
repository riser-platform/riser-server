package postgres

import (
	"database/sql"

	"github.com/riser-platform/riser-server/pkg/core"
)

type apiKeyRepository struct {
	db *sql.DB
}

func NewApiKeyRepository(db *sql.DB) core.ApiKeyRepository {
	return &apiKeyRepository{db}
}

func (r *apiKeyRepository) GetByUserId(userId int) ([]core.ApiKey, error) {
	apiKeys := []core.ApiKey{}
	rows, err := r.db.Query(`
	SELECT id, riser_user_id, key_hash
	FROM apikey
	WHERE riser_user_id = $1
	`, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		apiKey := core.ApiKey{}
		err := rows.Scan(&apiKey.Id, &apiKey.UserId, &apiKey.KeyHash)
		if err != nil {
			return nil, err
		}
		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys, nil
}

func (r *apiKeyRepository) Create(userId int, keyHash []byte) (id int, err error) {
	err = r.db.QueryRow("INSERT INTO apikey (riser_user_id, key_hash) VALUES ($1, $2) RETURNING id", userId, keyHash).Scan(&id)
	return id, err
}
