package core

type ResourceFile struct {
	Name     string `json:"name"`
	Contents []byte `json:"contents"`
	Delete   bool   `json:"delete,omitempty"`
}
