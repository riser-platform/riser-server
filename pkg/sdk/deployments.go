package sdk

import (
	"fmt"
	"net/http"

	"github.com/riser-platform/riser-server/api/v1/model"
)

type DeploymentsClient interface {
	Delete(deploymentName, namespace, envName string) (*model.SaveDeploymentResponse, error)
	Save(deployment *model.SaveDeploymentRequest, dryRun bool) (*model.SaveDeploymentResponse, error)
	SaveStatus(deploymentName, namespace, envName string, status *model.DeploymentStatusMutable) (statusCode int, err error)
}

type deploymentsClient struct {
	client *Client
}

func (c *deploymentsClient) Delete(deploymentName, namespace, envName string) (*model.SaveDeploymentResponse, error) {
	request, err := c.client.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/deployments/%s/%s/%s", envName, namespace, deploymentName), nil)
	if err != nil {
		return nil, err
	}

	responseModel := &model.SaveDeploymentResponse{}
	_, err = c.client.Do(request, responseModel)
	if err != nil {
		return nil, err
	}

	return responseModel, nil
}

func (c *deploymentsClient) Save(deployment *model.SaveDeploymentRequest, dryRun bool) (*model.SaveDeploymentResponse, error) {
	request, err := c.client.NewRequest(http.MethodPut, "/api/v1/deployments", deployment)
	if err != nil {
		return nil, err
	}

	if dryRun {
		q := request.URL.Query()
		q.Add("dryRun", "true")
		request.URL.RawQuery = q.Encode()
	}

	responseModel := &model.SaveDeploymentResponse{}
	_, err = c.client.Do(request, responseModel)
	if err != nil {
		return nil, err
	}

	return responseModel, nil
}

func (c *deploymentsClient) SaveStatus(deploymentName, namespace, envName string, status *model.DeploymentStatusMutable) (statusCode int, err error) {
	request, err := c.client.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/deployments/%s/%s/%s/status", envName, namespace, deploymentName), status)
	if err != nil {
		return 0, err
	}
	response, err := c.client.Do(request, nil)
	if err != nil {
		return response.StatusCode, err
	}
	return response.StatusCode, nil
}
