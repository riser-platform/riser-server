package login

import (
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/assert"
)

const testValidKey = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func Test_LoginWithApiKey(t *testing.T) {
	plainText := "aabbccdd"
	var hash []byte
	user := &core.User{Username: "test"}
	userRepository := &core.FakeUserRepository{
		GetByApiKeyFn: func(hashArg []byte) (*core.User, error) {
			hash = hashArg
			return user, nil
		},
	}
	service := service{users: userRepository}

	result, err := service.LoginWithApiKey(plainText)

	assert.Equal(t, user, result)
	assert.NoError(t, err)
	assert.Equal(t, hashApiKey([]byte(plainText)), hash)
}

func Test_LoginWithApiKey_Trims(t *testing.T) {
	plainText := " aabbccdd "
	var hash []byte
	user := &core.User{Username: "test"}
	userRepository := &core.FakeUserRepository{
		GetByApiKeyFn: func(hashArg []byte) (*core.User, error) {
			hash = hashArg
			return user, nil
		},
	}
	service := service{users: userRepository}

	result, err := service.LoginWithApiKey(plainText)

	assert.Equal(t, user, result)
	assert.NoError(t, err)
	assert.Equal(t, hashApiKey([]byte("aabbccdd")), hash)
}

func Test_LoginWithApiKey_NotFound_ReturnsError(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetByApiKeyFn: func([]byte) (*core.User, error) {
			return nil, core.ErrNotFound
		},
	}
	service := service{users: userRepository}

	result, err := service.LoginWithApiKey("nope")

	assert.Nil(t, result)
	assert.Equal(t, ErrInvalidLogin, err)
}

func Test_LoginWithApiKey_Error_ReturnsError(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetByApiKeyFn: func([]byte) (*core.User, error) {
			return nil, errors.New("test")
		},
	}
	service := service{users: userRepository}

	result, err := service.LoginWithApiKey("nope")

	assert.Nil(t, result)
	assert.Equal(t, "test", err.Error())
}

func Test_BootstrapRootUser(t *testing.T) {
	var rootUserId uuid.UUID
	userRepository := &core.FakeUserRepository{
		GetByUsernameFn: func(username string) (*core.User, error) {
			assert.Equal(t, RootUsername, username)
			return nil, core.ErrNotFound
		},
		CreateFn: func(newUser *core.NewUser) error {
			assert.NotEqual(t, uuid.Nil, newUser.Id)
			assert.Equal(t, RootUsername, newUser.Username)
			rootUserId = newUser.Id
			return nil
		},
	}
	apikeyRepository := &core.FakeApiKeyRepository{
		CreateFn: func(userId uuid.UUID, keyHash []byte) error {
			assert.Equal(t, rootUserId, userId)
			assert.Equal(t, hashApiKey([]byte(testValidKey)), keyHash)
			return nil
		},
	}

	service := service{userRepository, apikeyRepository}

	err := service.BootstrapRootUser(testValidKey)

	assert.NoError(t, err)
	assert.Equal(t, 1, userRepository.CreateCallCount)
	assert.Equal(t, 1, apikeyRepository.CreateCallCount)
}

func Test_BootstrapRootUser_UnableToCreateApiKey_ReturnsError(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetByUsernameFn: func(username string) (*core.User, error) {
			return nil, core.ErrNotFound
		},
		CreateFn: func(newUser *core.NewUser) error {
			return nil
		},
	}
	apikeyRepository := &core.FakeApiKeyRepository{
		CreateFn: func(uuid.UUID, []byte) error {
			return errors.New("test")
		},
	}

	service := service{userRepository, apikeyRepository}

	err := service.BootstrapRootUser(testValidKey)

	assert.Equal(t, "Error creating root API key: test", err.Error())
}

func Test_BootstrapRootUser_UserWithKeyExists_ReturnsError(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetByUsernameFn: func(string) (*core.User, error) {
			return &core.User{Id: uuid.New()}, nil
		},
	}
	apikeyRepository := &core.FakeApiKeyRepository{
		GetByUserIdFn: func(uuid.UUID) ([]core.ApiKey, error) {
			return []core.ApiKey{core.ApiKey{}}, nil
		},
	}

	service := service{userRepository, apikeyRepository}

	err := service.BootstrapRootUser(testValidKey)

	assert.Equal(t, ErrRootUserExists, err)
}

func Test_BootstrapRootUser_UnableToCreateUser_ReturnsError(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetByUsernameFn: func(string) (*core.User, error) {
			return nil, core.ErrNotFound
		},
		CreateFn: func(newUser *core.NewUser) error {
			return errors.New("test")
		},
	}
	apikeyRepository := &core.FakeApiKeyRepository{}

	service := service{userRepository, apikeyRepository}

	err := service.BootstrapRootUser(testValidKey)

	assert.Equal(t, "Unable to create root user: test", err.Error())
}

func Test_BootstrapRootUser_UnableToQueryKeys_ReturnsError(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetByUsernameFn: func(string) (*core.User, error) {
			return &core.User{Id: uuid.New()}, nil
		},
	}
	apikeyRepository := &core.FakeApiKeyRepository{
		GetByUserIdFn: func(uuid.UUID) ([]core.ApiKey, error) {
			return nil, errors.New("test")
		},
	}

	service := service{userRepository, apikeyRepository}

	err := service.BootstrapRootUser(testValidKey)

	assert.Equal(t, "Unable to retrieve root API keys: test", err.Error())
}

func Test_BootstrapRootUser_UnableToQueryUser_ReturnsError(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetByUsernameFn: func(string) (*core.User, error) {
			return nil, errors.New("test")
		},
	}

	service := service{users: userRepository}

	err := service.BootstrapRootUser(testValidKey)

	assert.Equal(t, "Unable to retrieve root user: test", err.Error())
}

func Test_BootstrapRootUser_ShortKey_ReturnsError(t *testing.T) {
	service := service{}

	err := service.BootstrapRootUser("oops")

	assert.Equal(t, "API Key must be a minimum of 32 characters. It is highly recommended to use `riser ops generate-apikey` to generate the key", err.Error())
}

func Test_BootstrapRootUser_EmptyKeyWithUsers_NOOP(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetActiveCountFn: func() (int, error) {
			return 1, nil
		},
	}
	service := service{users: userRepository}

	err := service.BootstrapRootUser("")

	assert.NoError(t, err)
}

func Test_BootstrapRootUser_EmptyKeyNoUsers_NOOP(t *testing.T) {
	userRepository := &core.FakeUserRepository{
		GetActiveCountFn: func() (int, error) {
			return 0, nil
		},
	}
	service := service{users: userRepository}

	err := service.BootstrapRootUser("")

	assert.Equal(t, "You must specify RISER_BOOTSTRAP_APIKEY is required when there are no users. Use \"riser ops generate-apikey\" to generate the key.", err.Error())
}
