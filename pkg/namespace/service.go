package namespace

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
)

type Service interface {
	// ValidateDeployable validates that a namespace is deployable. Returns a ValidationError if it is not.
	ValidateDeployable(namespaceName string) error
	// EnsureDefaultNamespace ensures that the default namespace has been provisioned. Designed to be used only at server startup.
	EnsureDefaultNamespace() error
	Create(namespaceName string) error
}

type service struct {
	namespaces   core.NamespaceRepository
	environments core.EnvironmentRepository
}

func NewService(namespaces core.NamespaceRepository, environments core.EnvironmentRepository) Service {
	return &service{namespaces, environments}
}

func (s *service) EnsureDefaultNamespace() error {
	_, err := s.namespaces.Get(core.DefaultNamespace)
	if err != nil {
		if err == core.ErrNotFound {
			return s.Create(core.DefaultNamespace)
		}
		return err
	}

	return nil
}

func (s *service) Create(namespaceName string) error {
	err := s.namespaces.Create(&core.Namespace{Name: namespaceName})
	if err != nil {
		return errors.Wrap(err, "error creating namespace")
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
