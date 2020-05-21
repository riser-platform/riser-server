package core

type SecretMeta struct {
	Name            string
	App             *NamespacedName
	EnvironmentName string
	Revision        int64
}
