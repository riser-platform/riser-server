package model

type DeploymentRequest struct {
	DeploymentMeta `json:",inline"`
	App            *AppConfigWithOverrides `json:"app"`
}

func (d *DeploymentRequest) ApplyDefaults() error {
	if d.App == nil {
		d.App = &AppConfigWithOverrides{}
	}
	return d.App.ApplyDefaults()
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

type DeploymentMeta struct {
	Name          string           `json:"name"`
	Stage         string           `json:"stage"`
	Docker        DeploymentDocker `json:"docker"`
	ManualRollout bool             `json:"manualRollout"`
}

type DeploymentDocker struct {
	Tag string `json:"tag"`
}
