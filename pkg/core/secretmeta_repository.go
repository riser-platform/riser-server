package core

type SecretMetaRepository interface {
	// Commit commits a secretmeta at the specified revision. This is an acknowledgement that the secret has been committed to the underlying
	// resource (e.g. the git state repo)
	Commit(secretMeta *SecretMeta) error
	// Save saves the secret meta and returns a new revision. It is up to the caller to modify the secretMeta with the new revision.
	// Important: You must call r.Commit to validate that the object has been committed. Uncommitted secrets are not applied to deployments
	Save(secretMeta *SecretMeta) (revision int64, err error)
	FindByStage(appName string, stageName string) ([]SecretMeta, error)
}

type FakeSecretMetaRepository struct {
	CommitFn        func(*SecretMeta) error
	CommitCallCount int
	SaveFn          func(*SecretMeta) (int64, error)
	SaveCallCount   int
	FindByStageFn   func(string, string) ([]SecretMeta, error)
}

func (fake *FakeSecretMetaRepository) Save(secretMeta *SecretMeta) (int64, error) {
	fake.SaveCallCount++
	return fake.SaveFn(secretMeta)
}

func (fake *FakeSecretMetaRepository) FindByStage(appName string, stageName string) ([]SecretMeta, error) {
	return fake.FindByStageFn(appName, stageName)
}

func (fake *FakeSecretMetaRepository) Commit(secretMeta *SecretMeta) error {
	fake.CommitCallCount++
	return fake.CommitFn(secretMeta)
}
