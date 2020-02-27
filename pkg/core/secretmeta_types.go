package core

import "github.com/google/uuid"

type SecretMeta struct {
	Name      string
	AppId     uuid.UUID
	StageName string
	Revision  int64
}
