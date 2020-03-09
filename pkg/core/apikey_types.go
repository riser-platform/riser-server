package core

import "github.com/google/uuid"

const (
	LoginTypeAPIKey = "APIKey"
)

type ApiKey struct {
	UserId  uuid.UUID `json:"userId"`
	KeyHash []byte    `json:"keyHash"`
}
