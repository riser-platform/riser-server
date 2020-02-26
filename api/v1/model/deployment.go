package model

import validation "github.com/go-ozzo/ozzo-validation/v3"

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

func (d DeploymentRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.DeploymentMeta),
		validation.Field(&d.App, validation.Required))
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

func (d DeploymentMeta) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, append(RulesNamingIdentifier(), validation.Required)...),
		validation.Field(&d.Stage, validation.Required))
}

type DeploymentDocker struct {
	Tag string `json:"tag"`
}
