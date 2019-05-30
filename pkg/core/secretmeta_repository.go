package core

type SecretMetaRepository interface {
	Save(secretMeta *SecretMeta) error
	FindByStage(appName string, stageName string) ([]SecretMeta, error)
}

type FakeSecretMetaRepository struct {
	SaveFn        func(*SecretMeta) error
	SaveCallCount int
	FindByStageFn func(string, string) ([]SecretMeta, error)
}

func (fake *FakeSecretMetaRepository) Save(secretMeta *SecretMeta) error {
	fake.SaveCallCount++
	return fake.SaveFn(secretMeta)
}

func (fake *FakeSecretMetaRepository) FindByStage(appName string, stageName string) ([]SecretMeta, error) {
	return fake.FindByStageFn(appName, stageName)
}
