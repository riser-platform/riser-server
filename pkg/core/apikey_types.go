package core

const (
	LoginTypeAPIKey = "APIKey"
)

type ApiKey struct {
	UserId  int    `json:"userId"`
	KeyHash []byte `json:"keyHash"`
}
