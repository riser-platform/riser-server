package sdk

import (
	"fmt"
	"net/http"

	"github.com/riser-platform/riser-server/api/v1/model"
)

type EnvironmentsClient interface {
	Ping(envName string) error
	List() ([]model.EnvironmentMeta, error)
	GetConfig(envName string) (*model.EnvironmentConfig, error)
	SetConfig(envName string, config *model.EnvironmentConfig) error
}

type environmentsClient struct {
	client *Client
}

func (c *environmentsClient) Ping(envName string) error {
	request, err := c.client.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/environments/%s/ping", envName), nil)
	if err != nil {
		return err
	}

	_, err = c.client.Do(request, nil)
	return err
}

func (c *environmentsClient) List() ([]model.EnvironmentMeta, error) {
	request, err := c.client.NewGetRequest("/api/v1/environments")
	if err != nil {
		return nil, err
	}

	environments := []model.EnvironmentMeta{}
	_, err = c.client.Do(request, &environments)
	if err != nil {
		return nil, err
	}

	return environments, nil
}

// GetConfig gets the configuration for a environment.
func (c *environmentsClient) GetConfig(envName string) (*model.EnvironmentConfig, error) {
	request, err := c.client.NewGetRequest(fmt.Sprintf("/api/v1/environments/%s/config", envName))
	if err != nil {
		return nil, err
	}

	responseModel := &model.EnvironmentConfig{}
	_, err = c.client.Do(request, responseModel)
	if err != nil {
		return nil, err
	}

	return responseModel, nil
}

// SetConfig sets configuration for a environment. Empty values are ignored and merged with existing config values.
func (c *environmentsClient) SetConfig(envName string, config *model.EnvironmentConfig) error {
	request, err := c.client.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/environments/%s/config", envName), config)
	if err != nil {
		return err
	}

	_, err = c.client.Do(request, nil)
	return err
}
