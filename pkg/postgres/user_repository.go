package postgres

import (
	"database/sql"
	"time"

	"github.com/riser-platform/riser-server/pkg/core"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) core.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetByApiKey(keyHash []byte) (*core.User, error) {
	user := &core.User{}
	err := r.db.QueryRow(`SELECT riser_user.id, username, doc
	FROM riser_user
	INNER JOIN apikey ON (riser_user.id = apikey.riser_user_id)
	WHERE apikey.key_hash = $1`, keyHash).Scan(&user.Id, &user.Username, &user.Doc)
	if err == sql.ErrNoRows {
		return nil, core.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// TODO: Refactor common projection, scanning, and error handling
func (r *userRepository) GetByUsername(username string) (*core.User, error) {
	user := &core.User{}
	err := r.db.QueryRow(`SELECT id, username, doc
	FROM riser_user
	WHERE username = $1`, username).Scan(&user.Id, &user.Username, &user.Doc)
	if err == sql.ErrNoRows {
		return nil, core.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Create(newUser *core.NewUser) error {
	doc := &core.UserDoc{Created: time.Now().UTC()}
	_, err := r.db.Exec("INSERT INTO riser_user (id, username, doc) VALUES ($1, $2, $3) RETURNING id", newUser.Id, newUser.Username, doc)
	return err
}

func (r *userRepository) GetActiveCount() (activeUserCount int, err error) {
	err = r.db.QueryRow("SELECT COUNT(1) FROM riser_user INNER JOIN apikey ON riser_user.id = apikey.riser_user_id").Scan(&activeUserCount)
	return activeUserCount, err
}
