package core

type ApiKeyRepository interface {
	GetByUserId(userId int) ([]ApiKey, error)
	Create(userId int, keyHash []byte) (int, error)
}

type FakeApiKeyRepository struct {
	GetByUserIdFn   func(int) ([]ApiKey, error)
	CreateFn        func(int, []byte) (int, error)
	CreateCallCount int
}

func (r *FakeApiKeyRepository) GetByUserId(userId int) ([]ApiKey, error) {
	return r.GetByUserIdFn(userId)
}

func (r *FakeApiKeyRepository) Create(userId int, keyHash []byte) (int, error) {
	r.CreateCallCount++
	return r.CreateFn(userId, keyHash)
}
