package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/riser-platform/riser-server/pkg/git"

	"github.com/riser-platform/riser-server/pkg/deployment"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/state"
)

func Test_DeleteDeployment(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/deployments/mydep/dev", nil)
	req.Header.Add("CONTENT-TYPE", "application/json")
	ctx, rec := newContextWithRecorder(req)

	deploymentService := &deployment.FakeService{
		DeleteFn: func(deploymentName, namespace, stageName string, committer state.Committer) error {
			return nil
		},
	}

	err := DeleteDeployment(ctx, nil, deploymentService)

	assert.NoError(t, err)
	assert.Equal(t, 1, deploymentService.DeleteCallCount)
	assert.Equal(t, http.StatusAccepted, rec.Result().StatusCode)
	apiResponse := model.APIResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiResponse))
	assert.Equal(t, "Deployment deletion requested", apiResponse.Message)
}

func Test_DeleteDeployment_NothingToDelete(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/deployments/mydep/dev", nil)
	req.Header.Add("CONTENT-TYPE", "application/json")
	ctx, rec := newContextWithRecorder(req)

	deploymentService := &deployment.FakeService{
		DeleteFn: func(deploymentName, namespace, stageName string, committer state.Committer) error {
			return git.ErrNoChanges
		},
	}

	err := DeleteDeployment(ctx, nil, deploymentService)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
	apiResponse := model.APIResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiResponse))
	assert.Equal(t, "Deployment not found", apiResponse.Message)
}

func Test_mapDryRunCommitsFromDomain(t *testing.T) {
	commits := []state.DryRunCommit{
		state.DryRunCommit{
			Message: "commit1",
			Files: []core.ResourceFile{
				core.ResourceFile{
					Name:     "file1",
					Contents: []byte("contents1"),
				},
				core.ResourceFile{
					Name:     "file2",
					Contents: []byte("contents2"),
				},
			},
		},
		state.DryRunCommit{
			Message: "commit2",
		},
	}

	result := mapDryRunCommitsFromDomain(commits)

	assert.Len(t, result, 2)
	assert.Equal(t, "commit1", result[0].Message)
	assert.Len(t, result[0].Files, 2)
	assert.Equal(t, "file1", result[0].Files[0].Name)
	assert.Equal(t, "contents1", result[0].Files[0].Contents)
	assert.Equal(t, "file2", result[0].Files[1].Name)
	assert.Equal(t, "contents2", result[0].Files[1].Contents)
	assert.Equal(t, result[1].Message, "commit2")
	assert.Empty(t, result[1].Files)
}

func Test_mapDeploymentRequestToDomain(t *testing.T) {
	request := &model.DeploymentRequest{
		DeploymentMeta: model.DeploymentMeta{
			Name:  "mydeployment",
			Stage: "mystage",
			Docker: model.DeploymentDocker{
				Tag: "mytag",
			},
			ManualRollout: true,
		},
		App: &model.AppConfigWithOverrides{
			AppConfig: model.AppConfig{
				Name: "myapp",
			},
		},
	}

	result, err := mapDeploymentRequestToDomain(request)

	assert.NoError(t, err)
	assert.Equal(t, "mydeployment", result.Name)
	assert.Equal(t, DefaultNamespace, result.Namespace)
	assert.Equal(t, "mystage", result.Stage)
	assert.Equal(t, "mytag", result.Docker.Tag)
	assert.Equal(t, request.App.AppConfig, *result.App)
	assert.True(t, result.ManualRollout)
}

func Test_mapDeploymentRequestToDomain_Overrides(t *testing.T) {
	request := &model.DeploymentRequest{
		DeploymentMeta: model.DeploymentMeta{
			Name:  "mydeployment",
			Stage: "mystage",
			Docker: model.DeploymentDocker{
				Tag: "mytag",
			},
		},
		App: &model.AppConfigWithOverrides{
			AppConfig: model.AppConfig{},
			Overrides: map[string]model.AppConfig{
				"mystage": model.AppConfig{
					Expose: &model.AppConfigExpose{ContainerPort: 1337},
				},
			},
		},
	}

	result, err := mapDeploymentRequestToDomain(request)

	assert.NoError(t, err)
	assert.Equal(t, int32(1337), result.App.Expose.ContainerPort)

}
