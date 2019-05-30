package core

const (
	LoginTypeAPIKey = "APIKey"
)

type ApiKey struct {
	Id      int    `json:"id"`
	UserId  int    `json:"userId"`
	KeyHash []byte `json:"keyHash"`
}
