package core

import (
	"fmt"
	"strings"
)

const DefaultNamespace = "apps"

type Namespace struct {
	Name string
}

type NamespacedName struct {
	Name      string
	Namespace string
}

func (v *NamespacedName) String() string {
	return fmt.Sprintf("%s.%s", v.Name, v.Namespace)
}

func ParseNamespacedName(namespacedName string) *NamespacedName {
	parts := strings.Split(namespacedName, ".")
	if len(parts) == 2 {
		return &NamespacedName{parts[0], parts[1]}
	}
	return &NamespacedName{Name: namespacedName}
}

func NewNamespacedName(name, namespace string) *NamespacedName {
	if namespace == "" {
		namespace = DefaultNamespace
	}
	return &NamespacedName{Name: name, Namespace: namespace}
}
