package state

import "k8s.io/apimachinery/pkg/runtime/schema"

type KubeResource interface {
	GetName() string
	GetNamespace() string
	GetObjectKind() schema.ObjectKind
}
