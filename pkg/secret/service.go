package secret

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/riser-platform/riser-server/pkg/state/resources"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/state"
)

type Service interface {
	SealAndSave(plaintextSecret string, secretMeta *core.SecretMeta, committer state.Committer) error
}

type service struct {
	secretMetas  core.SecretMetaRepository
	environments core.EnvironmentRepository
	rand         io.Reader
}

func NewService(secretMetas core.SecretMetaRepository, environments core.EnvironmentRepository) Service {
	return &service{secretMetas, environments, rand.Reader}
}

func (s *service) SealAndSave(plaintextSecret string, secretMeta *core.SecretMeta, committer state.Committer) error {
	sealedSecretCert, err := s.getSealedSecretCert(plaintextSecret, secretMeta.EnvironmentName)
	if err != nil {
		return err
	}

	return s.sealAndSave(plaintextSecret, sealedSecretCert, secretMeta, committer)
}

func (s *service) sealAndSave(plaintextSecret string, sealedSecretCert []byte, secretMeta *core.SecretMeta, committer state.Committer) error {
	revision, err := s.secretMetas.Save(secretMeta)
	if err != nil {
		return errors.Wrap(err, "Error saving secret metadata")
	}

	secretMeta.Revision = revision

	sealedSecret, err := resources.CreateSealedSecret(plaintextSecret, secretMeta, sealedSecretCert, s.rand)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error creating sealed secret %q in environment %q", secretMeta.Name, secretMeta.EnvironmentName))
	}

	resourceFiles, err := state.RenderSealedSecret(secretMeta.App.Name, secretMeta.EnvironmentName, sealedSecret)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error rendering sealed secret resource %q in environment %q", secretMeta.Name, secretMeta.EnvironmentName))
	}

	err = committer.Commit(fmt.Sprintf("Updating secret %q in environment %q", sealedSecret.Name, secretMeta.EnvironmentName), resourceFiles)
	if err != nil {
		return errors.Wrap(err, "Error committing sealed secret resources")
	}

	err = s.secretMetas.Commit(secretMeta)
	if err != nil {
		// Let the client handle this error specifically
		if err == core.ErrConflictNewerVersion {
			return err
		}
		return errors.Wrap(err, "Error committing sealed secret metadata")
	}
	return nil
}

func (s *service) getSealedSecretCert(plaintextSecret, envName string) ([]byte, error) {
	environment, err := s.environments.Get(envName)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error retrieving environment %q", envName))
	}

	if len(environment.Doc.Config.SealedSecretCert) == 0 {
		return nil, errors.New(fmt.Sprintf("No certificate configured in environment %q", envName))
	}

	return environment.Doc.Config.SealedSecretCert, nil
}
