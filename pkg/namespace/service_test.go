package namespace

import (
	"errors"
	"testing"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Create(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		CreateFn: func(namespace *core.Namespace) error {
			assert.Equal(t, "myns", namespace.Name)
			return nil
		},
	}

	environments := &core.FakeEnvironmentRepository{
		ListFn: func() ([]core.Environment, error) {
			return []core.Environment{
				{Name: "myenv1"},
				{Name: "myenv2"},
			}, nil
		},
	}

	svc := &service{namespaces, environments}

	err := svc.Create("myns")

	assert.NoError(t, err)
	assert.Equal(t, 1, namespaces.CreateCallCount)
}

func Test_Create_WhenNamespaceCreateErr(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		CreateFn: func(namespace *core.Namespace) error {
			return errors.New("test")
		},
	}

	svc := &service{namespaces: namespaces}

	err := svc.Create("myns")

	assert.Equal(t, "error creating namespace: test", err.Error())
}

func Test_EnsureDefaultNamespace_ReturnsErr(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		GetFn: func(string) (*core.Namespace, error) {
			return nil, errors.New("test")
		},
	}

	svc := &service{namespaces: namespaces}

	err := svc.EnsureDefaultNamespace()

	assert.Equal(t, "test", err.Error())
}

func Test_EnsureDefaultNamespace_WhenExists_Noop(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		GetFn: func(string) (*core.Namespace, error) {
			return &core.Namespace{}, nil
		},
	}

	svc := &service{namespaces: namespaces}

	err := svc.EnsureDefaultNamespace()

	assert.NoError(t, err)
}

func Test_ValidateDeployable_NamespaceExists(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		GetFn: func(namespaceArg string) (*core.Namespace, error) {
			assert.Equal(t, "myns", namespaceArg)
			return &core.Namespace{Name: namespaceArg}, nil
		},
	}

	svc := &service{namespaces: namespaces}

	err := svc.ValidateDeployable("myns")

	assert.NoError(t, err)
	assert.Equal(t, 1, namespaces.GetCallCount)
}

func Test_ValidateDeployable_NamespaceMissing(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		GetFn: func(namespaceArg string) (*core.Namespace, error) {
			return nil, core.ErrNotFound
		},
		ListFn: func() ([]core.Namespace, error) {
			return []core.Namespace{
				{Name: "ns1"},
				{Name: "ns2"},
			}, nil
		},
	}

	svc := &service{namespaces: namespaces}

	err := svc.ValidateDeployable("myns")

	require.IsType(t, &core.ValidationError{}, err, err.Error())
	assert.Equal(t, `Invalid namespace "myns". Must be one of: ns1, ns2`, err.Error())
}

func Test_ValidateDeployable_NamespaceMissing_ListError(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		GetFn: func(namespaceArg string) (*core.Namespace, error) {
			return nil, core.ErrNotFound
		},
		ListFn: func() ([]core.Namespace, error) {
			return nil, errors.New("test")
		},
	}

	svc := &service{namespaces: namespaces}

	err := svc.ValidateDeployable("myns")

	assert.Equal(t, "test", err.Error())
}

func Test_ValidateDeployable_GetError(t *testing.T) {
	namespaces := &core.FakeNamespaceRepository{
		GetFn: func(namespaceArg string) (*core.Namespace, error) {
			return nil, errors.New("test")
		},
	}

	svc := &service{namespaces: namespaces}

	err := svc.ValidateDeployable("myns")

	assert.Equal(t, "test", err.Error())
}
