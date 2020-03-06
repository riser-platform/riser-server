package core

type SecretMeta struct {
	Name      string
	App       *NamespacedName
	StageName string
	Revision  int64
}
