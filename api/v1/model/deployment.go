package model

import validation "github.com/go-ozzo/ozzo-validation/v3"

type SaveDeploymentRequest struct {
	DeploymentMeta `json:",inline"`
	App            *AppConfigWithOverrides `json:"app"`
}

func (d *SaveDeploymentRequest) ApplyDefaults() error {
	if d.App == nil {
		d.App = &AppConfigWithOverrides{}
	}
	return d.App.ApplyDefaults()
}

func (d SaveDeploymentRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.DeploymentMeta),
		validation.Field(&d.App, validation.Required))
}

type SaveDeploymentResponse struct {
	RiserRevision int64          `json:"riserRevision"`
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
	Name string `json:"name"`
	// Namespace is an intentional omission. We always use the app's namespace as we do not allow an app to deploy to multiple namespaces at
	// this time.
	Environment   string           `json:"environment"`
	Docker        DeploymentDocker `json:"docker"`
	ManualRollout bool             `json:"manualRollout"`
}

func (d DeploymentMeta) Validate() error {
	return validation.ValidateStruct(&d,
		// There's a separate RuneLength rule here to reserve 8 characters for the deployment prefix (e.g. for myapp: r100-myapp)
		validation.Field(&d.Name, append(RulesNamingIdentifier(), validation.RuneLength(3, 55), validation.Required)...),
		validation.Field(&d.Environment, validation.Required))
}

type DeploymentDocker struct {
	Tag string `json:"tag"`
}
