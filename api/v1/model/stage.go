package model

type StageMeta struct {
	Name string
}

type StageConfig struct {
	SealedSecretCert  []byte `json:"sealedSecretCert,omitempty"`
	PublicGatewayHost string `json:"publicGatewayHost,omitempty"`
}
