package secret

import (
	"fmt"
	"time"

	"github.com/riser-platform/riser-server/pkg/state/resources"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/state"
)

type Service interface {
	SealAndSave(plaintextSecret string, secretMeta *core.SecretMeta, namespace string, comitter state.Committer) error
	FindByStage(appName, stageName string) ([]core.SecretMeta, error)
	FindNamesByStage(appName, stageName string) ([]string, error)
}

type service struct {
	secretMetas core.SecretMetaRepository
	stages      core.StageRepository
}

func NewService(secretMetas core.SecretMetaRepository, stages core.StageRepository) Service {
	return &service{secretMetas: secretMetas, stages: stages}
}

func (s *service) FindNamesByStage(appName, stageName string) ([]string, error) {
	secretMetas, err := s.FindByStage(appName, stageName)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, secretMeta := range secretMetas {
		names = append(names, secretMeta.SecretName)
	}

	return names, nil
}

func (s *service) FindByStage(appName, stageName string) ([]core.SecretMeta, error) {
	secretMetas, err := s.secretMetas.FindByStage(appName, stageName)
	if err != nil {
		return nil, errors.Wrap(err, "Error retrieving secret metadata")
	}

	return secretMetas, nil
}

func (s *service) SealAndSave(plaintextSecret string, secretMeta *core.SecretMeta, namespace string, comitter state.Committer) error {
	stage, err := s.stages.Get(secretMeta.StageName)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error retrieving stage %q", secretMeta.StageName))
	}

	if len(stage.Doc.Config.SealedSecretCert) == 0 {
		return errors.New(fmt.Sprintf("No certificate configured in stage %q", secretMeta.StageName))
	}

	sealedSecret, err := resources.CreateSealedSecret(plaintextSecret, secretMeta, namespace, stage.Doc.Config.SealedSecretCert)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error creating sealed secret %q in stage %q", secretMeta.SecretName, secretMeta.StageName))
	}

	resourceFiles, err := state.RenderSealedSecret(secretMeta.AppName, secretMeta.StageName, sealedSecret)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error rendering sealed secret resource %q in stage %q", secretMeta.SecretName, secretMeta.StageName))
	}

	secretMeta.Doc.LastUpdated = time.Now().UTC()
	err = s.secretMetas.Save(secretMeta)
	if err != nil {
		return errors.Wrap(err, "Error saving secret metadata")
	}

	err = comitter.Commit(fmt.Sprintf("Updating secret %q in stage %q", sealedSecret.Name, secretMeta.StageName), resourceFiles)
	if err != nil {
		return errors.Wrap(err, "Error commiting sealed secret resources")
	}
	return nil
}
