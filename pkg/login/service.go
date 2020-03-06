package login

import (
	"crypto/sha256"
	"strings"

	"github.com/google/uuid"

	"github.com/pkg/errors"

	"github.com/riser-platform/riser-server/pkg/core"
)

const RootUsername = "root"
const ApiKeyMinCharacterLength = 32

// ErrRootUserExists is returned when the root user exists with an active API key
var ErrRootUserExists = errors.New("The root user already exists")

// ErrInvalidLogin is returned when the login is not valid
var ErrInvalidLogin = errors.New("Invalid login")

// ErrNoBoostrapNoUsers is returned when no bootstrap is specified and no active users are provisioned in the DB
var ErrNoBootstrapNoUsers = errors.New("You must specify RISER_BOOTSTRAP_APIKEY is required when there are no users. Use \"riser ops generate-apikey\" to generate the key.")

type Service interface {
	LoginWithApiKey(apiKeyPlainText string) (*core.User, error)
	BootstrapRootUser(apiKeyPlainText string) error
}

type service struct {
	users   core.UserRepository
	apikeys core.ApiKeyRepository
}

func NewService(users core.UserRepository, apikeys core.ApiKeyRepository) Service {
	return &service{users, apikeys}
}

func (s *service) LoginWithApiKey(apiKeyPlainText string) (*core.User, error) {
	hash := hashApiKey([]byte(strings.TrimSpace(apiKeyPlainText)))
	user, err := s.users.GetByApiKey(hash)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, ErrInvalidLogin
		}
		return nil, err
	}

	return user, nil
}

// BootstrapRootUser is an idempotent function that will create the root user with the specified API key if needed.
// An error is returned if the root user is already created.
// Passing an empty value for the API key results in a NOOP unless no logins are specified, in which case an operator friendly error is returned.
// If the root user exists with an API key, the request is ignored and ErrRootUserExists is returned
func (s *service) BootstrapRootUser(apiKeyPlainText string) error {
	if apiKeyPlainText == "" {
		activeUserCount, err := s.users.GetActiveCount()
		if err != nil {
			return err
		}

		if activeUserCount == 0 {
			return ErrNoBootstrapNoUsers
		}

		// We have active users so it's okay to ignore this bootstrap request
		return nil
	}

	if len(apiKeyPlainText) < ApiKeyMinCharacterLength {
		return errors.Errorf("API Key must be a minimum of %d characters. It is highly recommended to use `riser ops generate-apikey` to generate the key", ApiKeyMinCharacterLength)
	}

	rootUser, err := s.users.GetByUsername(RootUsername)
	var rootUserId uuid.UUID
	if err == nil {
		rootUserId = rootUser.Id
		// In the event that that the root user exists without an apikey
		apikeys, err := s.apikeys.GetByUserId(rootUserId)
		if err != nil {
			return errors.Wrap(err, "Unable to retrieve root API keys")
		}

		if len(apikeys) > 0 {
			return ErrRootUserExists
		}
	} else {
		if err != core.ErrNotFound {
			return errors.Wrap(err, "Unable to retrieve root user")
		}
		rootUserId = uuid.New()
		err = s.users.Create(&core.NewUser{Id: rootUserId, Username: RootUsername})
		if err != nil {
			return errors.Wrap(err, "Unable to create root user")
		}
	}

	err = s.apikeys.Create(rootUserId, hashApiKey([]byte(apiKeyPlainText)))
	if err != nil {
		return errors.Wrap(err, "Error creating root API key")
	}

	return nil
}

/*
Important! Changing this algorithm could be a breaking change:
	- Update the DB to store the algorithm used for existing keys and/or provide a way to rehash the keys on next login
	- Update the `riser ops generate-apikey` command to use the updated algorithm
*/
func hashApiKey(in []byte) []byte {
	arr := sha256.Sum256(in)
	// Convert to a slice because this is how we will always work with it.
	return arr[:]
}
