package model

type DeploymentRequest struct {
	DeploymentMeta `json:",inline"`
	App            *AppConfigWithOverrides `json:"app"`
}

type DeploymentResponse struct {
	Message       string         `json:"message"`
	DryRunCommits []DryRunCommit `json:"dryRunCommits,omitempty"`
}

type DryRunCommit struct {
	Message string       `json:"message"`
	Files   []DryRunFile `json:"files"`
}

type DryRunFile struct {
	Name     string `json:"name"`
	Contents string `json:"contents"`
}

type Deployment struct {
	DeploymentMeta `json:",inline"`
	App            *AppConfig `json:"app"`
}

type DeploymentMeta struct {
	Name      string           `json:"name"`
	Namespace string           `json:"namespace"`
	Stage     string           `json:"stage"`
	Docker    DeploymentDocker `json:"docker"`
}

type DeploymentDocker struct {
	Tag string `json:"tag"`
}
