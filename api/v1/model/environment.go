package model

type EnvironmentMeta struct {
	Name string
}

type EnvironmentConfig struct {
	SealedSecretCert  []byte `json:"sealedSecretCert,omitempty"`
	PublicGatewayHost string `json:"publicGatewayHost,omitempty"`
}
