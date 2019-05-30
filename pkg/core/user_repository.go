package core

type UserRepository interface {
	GetByApiKey(keyHash []byte) (*User, error)
	GetByUsername(username string) (*User, error)
	Create(newUser *NewUser) (int, error)
	GetActiveCount() (int, error)
}

type FakeUserRepository struct {
	GetByApiKeyFn    func(keyHash []byte) (*User, error)
	GetByUsernameFn  func(username string) (*User, error)
	CreateFn         func(newUser *NewUser) (int, error)
	CreateCallCount  int
	GetActiveCountFn func() (int, error)
}

func (r *FakeUserRepository) GetByApiKey(keyHash []byte) (*User, error) {
	return r.GetByApiKeyFn(keyHash)
}

func (r *FakeUserRepository) GetByUsername(username string) (*User, error) {
	return r.GetByUsernameFn(username)
}

func (r *FakeUserRepository) Create(newUser *NewUser) (int, error) {
	r.CreateCallCount++
	return r.CreateFn(newUser)
}

func (r *FakeUserRepository) GetActiveCount() (int, error) {
	return r.GetActiveCountFn()
}
