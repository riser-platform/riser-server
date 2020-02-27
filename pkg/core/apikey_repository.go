package core

import "github.com/google/uuid"

type ApiKeyRepository interface {
	GetByUserId(userId uuid.UUID) ([]ApiKey, error)
	Create(userId uuid.UUID, keyHash []byte) error
}

type FakeApiKeyRepository struct {
	GetByUserIdFn   func(uuid.UUID) ([]ApiKey, error)
	CreateFn        func(uuid.UUID, []byte) error
	CreateCallCount int
}

func (r *FakeApiKeyRepository) GetByUserId(userId uuid.UUID) ([]ApiKey, error) {
	return r.GetByUserIdFn(userId)
}

func (r *FakeApiKeyRepository) Create(userId uuid.UUID, keyHash []byte) error {
	r.CreateCallCount++
	return r.CreateFn(userId, keyHash)
}
