package secret

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/riser-platform/riser-server/pkg/state/resources"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/state"
)

type Service interface {
	SealAndSave(plaintextSecret string, secretMeta *core.SecretMeta, namespace string, committer state.Committer) error
	FindByStage(appId uuid.UUID, stageName string) ([]core.SecretMeta, error)
}

type service struct {
	apps        core.AppRepository
	secretMetas core.SecretMetaRepository
	stages      core.StageRepository
}

func NewService(apps core.AppRepository, secretMetas core.SecretMetaRepository, stages core.StageRepository) Service {
	return &service{apps, secretMetas, stages}
}

func (s *service) FindByStage(appId uuid.UUID, stageName string) ([]core.SecretMeta, error) {
	secretMetas, err := s.secretMetas.FindByStage(appId, stageName)
	if err != nil {
		return nil, errors.Wrap(err, "Error retrieving secret metadata")
	}

	return secretMetas, nil
}

// TODO: Move namespace parameter into secretMeta (probably as part of namespace feature)
func (s *service) SealAndSave(plaintextSecret string, secretMeta *core.SecretMeta, namespace string, committer state.Committer) error {
	sealedSecretCert, err := s.getSealedSecretCert(plaintextSecret, secretMeta.StageName)
	if err != nil {
		return err
	}

	return s.sealAndSave(plaintextSecret, namespace, sealedSecretCert, secretMeta, committer)
}

func (s *service) sealAndSave(plaintextSecret, namespace string, sealedSecretCert []byte, secretMeta *core.SecretMeta, committer state.Committer) error {
	app, err := s.apps.Get(secretMeta.AppId)
	if err != nil {
		return errors.Wrap(err, "Error getting app")
	}

	revision, err := s.secretMetas.Save(secretMeta)
	if err != nil {
		return errors.Wrap(err, "Error saving secret metadata")
	}

	secretMeta.Revision = revision

	sealedSecret, err := resources.CreateSealedSecret(plaintextSecret, app.Name, secretMeta, namespace, sealedSecretCert)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error creating sealed secret %q in stage %q", secretMeta.Name, secretMeta.StageName))
	}

	resourceFiles, err := state.RenderSealedSecret(app.Name, secretMeta.StageName, sealedSecret)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error rendering sealed secret resource %q in stage %q", secretMeta.Name, secretMeta.StageName))
	}

	err = committer.Commit(fmt.Sprintf("Updating secret %q in stage %q", sealedSecret.Name, secretMeta.StageName), resourceFiles)
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

func (s *service) getSealedSecretCert(plaintextSecret, stageName string) ([]byte, error) {
	stage, err := s.stages.Get(stageName)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error retrieving stage %q", stageName))
	}

	if len(stage.Doc.Config.SealedSecretCert) == 0 {
		return nil, errors.New(fmt.Sprintf("No certificate configured in stage %q", stageName))
	}

	return stage.Doc.Config.SealedSecretCert, nil
}
