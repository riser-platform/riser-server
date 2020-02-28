package postgres

import (
	"database/sql"

	"github.com/riser-platform/riser-server/pkg/core"
)

type namespaceRepository struct {
	db *sql.DB
}

func NewNamespaceRepository(db *sql.DB) core.NamespaceRepository {
	return &namespaceRepository{db: db}
}

func (r *namespaceRepository) Create(namespace *core.Namespace) error {
	_, err := r.db.Exec("INSERT INTO namespace (name) VALUES ($1)", namespace.Name)
	return err
}

func (r *namespaceRepository) Get(namespaceName string) (*core.Namespace, error) {
	ns := &core.Namespace{}
	// Effectively used just to make sure that the namespace exists. This will do in the future as we add more fields.
	err := r.db.QueryRow("SELECT name FROM namespace WHERE name = $1", namespaceName).Scan(&ns.Name)
	if err == sql.ErrNoRows {
		return nil, core.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return ns, nil
}
