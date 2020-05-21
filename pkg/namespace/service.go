package namespace

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/git"
	"github.com/riser-platform/riser-server/pkg/state"
	"github.com/riser-platform/riser-server/pkg/state/resources"
)

type Service interface {
	// ValidateDeployable validates that a namespace is deployable. Returns a ValidationError if it is not.
	ValidateDeployable(namespaceName string) error
	// EnsureDefaultNamespace ensures that the default namespace has been provisioned. Designed to be used only at server startup.
	EnsureDefaultNamespace(committer state.Committer) error
	// EnsureNamespaceInEnvironment ensures that a namespace has been committed to a environment. Returns an error if the namespace has not been created
	EnsureNamespaceInEnvironment(namespaceName string, envName string, committer state.Committer) error
	Create(namespaceName string, committer state.Committer) error
}

type service struct {
	namespaces   core.NamespaceRepository
	environments core.EnvironmentRepository
}

func NewService(namespaces core.NamespaceRepository, environments core.EnvironmentRepository) Service {
	return &service{namespaces, environments}
}

func (s *service) EnsureDefaultNamespace(committer state.Committer) error {
	_, err := s.namespaces.Get(core.DefaultNamespace)
	if err != nil {
		if err == core.ErrNotFound {
			return s.Create(core.DefaultNamespace, committer)
		}
		return err
	}

	return nil
}

func (s *service) EnsureNamespaceInEnvironment(namespaceName string, envName string, committer state.Committer) error {
	_, err := s.namespaces.Get(namespaceName)
	if err != nil {
		if err == core.ErrNotFound {
			return core.NewValidationErrorMessage(fmt.Sprintf("the namespace %q does not exist", namespaceName))
		}
		return err
	}
	err = commitNamespace(namespaceName, envName, committer)
	if err == git.ErrNoChanges {
		return nil
	}
	return err
}

func (s *service) Create(namespaceName string, committer state.Committer) error {
	err := s.namespaces.Create(&core.Namespace{Name: namespaceName})
	if err != nil {
		return errors.Wrap(err, "error creating namespace")
	}

	environments, err := s.environments.List()
	if err != nil {
		return err
	}
	for _, environment := range environments {
		err = commitNamespace(namespaceName, environment.Name, committer)
		if err != nil && err != git.ErrNoChanges {
			return err
		}
	}

	return nil
}

func (s *service) ValidateDeployable(namespaceName string) error {
	_, err := s.namespaces.Get(namespaceName)

	if err == core.ErrNotFound {
		namespaces, nsListErr := s.namespaces.List()
		if nsListErr == nil {
			validNamespaceNames := toNameList(namespaces)
			return core.NewValidationErrorMessage(fmt.Sprintf("Invalid namespace %q. Must be one of: %s", namespaceName, strings.Join(validNamespaceNames, ", ")))
		} else {
			return nsListErr
		}
	}

	return err
}

func toNameList(namespaces []core.Namespace) []string {
	names := []string{}
	for _, namespace := range namespaces {
		names = append(names, namespace.Name)
	}
	return names
}

func commitNamespace(namespaceName string, envName string, committer state.Committer) error {
	nsResource, err := resources.CreateNamespace(namespaceName, envName)
	if err != nil {
		return err
	}
	resourceFiles, err := state.RenderGeneric(envName, nsResource)
	if err != nil {
		return err
	}

	return committer.Commit(fmt.Sprintf("Updating namespace %q in environment %q", namespaceName, envName), resourceFiles)
}
