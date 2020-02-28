package namespace

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/git"
	"github.com/riser-platform/riser-server/pkg/state"
	"github.com/riser-platform/riser-server/pkg/state/resources"
)

// DefaultName is the default namespace
const DefaultName = "apps"

type Service interface {
	EnsureDefaultNamespace(committer state.Committer) error
	// EnsureNamespaceInStage ensures that a namespace has been committed to a stage. Returns an error if the namespace has not been created
	EnsureNamespaceInStage(namespaceName string, stageName string, committer state.Committer) error
	Create(namespaceName string, committer state.Committer) error
}

type service struct {
	namespaces core.NamespaceRepository
	stages     core.StageRepository
}

func NewService(namespaces core.NamespaceRepository, stages core.StageRepository) Service {
	return &service{namespaces, stages}
}

func (s *service) EnsureDefaultNamespace(committer state.Committer) error {
	_, err := s.namespaces.Get(DefaultName)
	if err != nil {
		if err == core.ErrNotFound {
			return s.Create(DefaultName, committer)
		}
		return err
	}

	return nil
}

func (s *service) EnsureNamespaceInStage(namespaceName string, stageName string, committer state.Committer) error {
	_, err := s.namespaces.Get(namespaceName)
	if err != nil {
		if err == core.ErrNotFound {
			return core.NewValidationErrorMessage(fmt.Sprintf("the namespace %q does not exist", namespaceName))
		}
		return err
	}
	err = commitNamespace(namespaceName, stageName, committer)
	if err == git.ErrNoChanges {
		return nil
	}
	return err
}

func (s *service) Create(namespaceName string, committer state.Committer) error {
	// TODO: Validate name against ban list
	err := s.namespaces.Create(&core.Namespace{Name: namespaceName})
	if err != nil {
		return errors.Wrap(err, "error creating namespace")
	}

	stages, err := s.stages.List()
	if err != nil {
		return err
	}
	for _, stage := range stages {
		err = commitNamespace(namespaceName, stage.Name, committer)
		if err != nil && err != git.ErrNoChanges {
			return err
		}
	}

	return nil
}

func commitNamespace(namespaceName string, stageName string, committer state.Committer) error {
	nsResource, err := resources.CreateNamespace(namespaceName, stageName)
	if err != nil {
		return err
	}
	resourceFiles, err := state.RenderGeneric(stageName, nsResource)
	if err != nil {
		return err
	}

	return committer.Commit(fmt.Sprintf("Updating namespace %q in stage %q", namespaceName, stageName), resourceFiles)
}
