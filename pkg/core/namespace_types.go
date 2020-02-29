package core

import (
	"fmt"
	"strings"
)

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

// TODO: Should we allow non-namespaced names and use default?
func ParseNamespacedName(namespacedName string) *NamespacedName {
	parts := strings.Split(namespacedName, ".")
	if len(parts) == 2 {
		return &NamespacedName{parts[0], parts[1]}
	}
	return &NamespacedName{Name: namespacedName}
}

func NewNamespacedName(name, namespace string) *NamespacedName {
	return &NamespacedName{Name: name, Namespace: namespace}
}
