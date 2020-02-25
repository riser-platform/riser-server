package namespace

import (
	"fmt"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/state"
	"github.com/riser-platform/riser-server/pkg/state/resources"
)

// Default is the default namespace
const Default = "apps"

type Service interface {
	Save(namespace *core.Namespace, committer state.Committer) error
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) Save(namespace *core.Namespace, committer state.Committer) error {
	nsResource, err := resources.CreateNamespace(namespace)
	if err != nil {
		return err
	}
	resourceFiles, err := state.RenderGeneric(namespace.Stage, nsResource)
	if err != nil {
		return err
	}

	return committer.Commit(fmt.Sprintf("Updating namespace %q in stage %q", namespace.Name, namespace.Stage), resourceFiles)
}
