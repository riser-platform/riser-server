package secret

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/riser-platform/riser-server/pkg/state"

	"github.com/pkg/errors"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

const testCert = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVyakNDQXBhZ0F3SUJBZ0lSQUswbGdFTFB0MVE1UXRTQ1BKQ0k4Y1V3RFFZSktvWklodmNOQVFFTEJRQXcKQURBZUZ3MHhPVEE1TWpReE1ERTJNemxhRncweU9UQTVNakV4TURFMk16bGFNQUF3Z2dJaU1BMEdDU3FHU0liMwpEUUVCQVFVQUE0SUNEd0F3Z2dJS0FvSUNBUURBNURtRGt4dzlzRFhtVjJ4M0EyL3JBYUpHc1NVZzJ2Uk5NWFVtCndkdWJlTVB0eDhicHU4bEtCZVJlQTFuSk5QT1diQ0tlMGxTeVVQTUw4RFBRWmF5eEtIMUhPQnd3MDVrdmhJQ1YKLzFMcGlZbzRNVC91WWlLcS9vODhMY0JDTTBYdHBzZGFjL05LZExyalg2VWhoaEVPQ2VmWnIwSEJZbTBOTDdEUwp5OFVZN0ZucnZ0YUo0aVBCNTcrK0R3bUJkS3lzMGJzK3ZtVXp5L0Q2UjdHUGxPTTF4YmFrWjZ4NElrUXNmbVlTClVRR1RNMzFCb3Y5d2x6dksyS0pET2NNNU9JdG1TdGVTekMvbExzdllPWVpqSE95Q01iVG8veStMdldlZXFCTHYKNU9GZytTelVLSWNOT2RQWUdOSVNwRk9UNUtzVTlrMU55UXVNbWhKZ1dEdWpaTkdxT1IvWmNTK0M5RjBLTXlWcgpZSkdQWS9rYjRLNFZ3bUlkbTdodk1OZkU4QnRjS0UxVzJaYzNPOWZZeFN3ZDByQk9oZ0J4WndHUU4zTE8rbkVFCnA5TGt6UjhxMVlCV3VYa0VJZGd0Y1l1UkJEb0ZUcTRMUXNzazdQY2JjUWx6dHdNbUFPTWtKNVVGcktjZVU3YzgKVnlnOHphSVcxekNySVdUV2RMeEJIa0FZL0ZURVRGMkNVL3NTOWxTd090elRLdDVPVUdld3EzOS94TlpPVmlVLwp0RnhzZ3AxSjFYNk16SGlOTzdGVHVaTitKcXpmcnI4ekZzcEVDR29vVmZMd0I1MERnSnZXeFc4bWptTkplQWN6Ck9PVXR2TEpoK29TK3NJRTdqYmhROTJpTW42V2VNdUhMSEJxL2RXWmFqYmozbDJhaHF0emNPLy95ejB1VVdUclUKVnZqVkxRSURBUUFCb3lNd0lUQU9CZ05WSFE4QkFmOEVCQU1DQUFFd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBTgpCZ2txaGtpRzl3MEJBUXNGQUFPQ0FnRUFKL2lNdEJ3Z3RaY1JpNlRudnVKQ3IvT2RYblZFaVRDeWtkOVREeWxRCnBMRXNqUWt4TjhwT2x3R1V5K3FOeGM0SjJNSmo2WG9KK3RrVXNRSVZTOHpzcTNWblU3YzBpd1VRUWZSZUdXNmcKRlVLTnNUaVpwTUplVnp0aUpRT3FJNlZBUnNEMFZNSzVXWnUzc2w0VCsrVUl6SXVsWUJLUWx2UzFGNmpYNjdXUQplenVLR250NC85OUFYTG53cjBFaVh5U1hQcFhzVklJNXBlSlU4K1BoN1dHRkFxb04wMFlIMUFaSEJaU2JpYjRvClVXU21NWlhlSDJ4aXVIUUVKRCtkdENGeUNYU1R6OFRNeXpaYy90NGtaZFFvWGVST3FxZ09QR0FCcmR5R2lkNHkKY2tjeGcxbDRTWXZJeHlmNktJWFZGQ2VaUVRFcU9xRzQyQzR2WjZSTlEzR093VGFzQSsyK1ZNNUNZQlVkSlIzWgpzZkI4djVFK3VPZ3RrQXpGNUV0SjgxODcwUXE3ZXgwVDUzREEvd2JoVDR5RkJWb25UU1p4cGFUaGlyUEN1RUVXClErOGpPWWhPWTNZeDVnekwzY1puclVLRk9XaVltOWVVVWpWZnB1R3JBdmJzK0svZlJMNWp4UDFFUWhydDhDRnoKOHlxc2lrYm1NM0FDeGxpakJqRWtrM0Zsa3ZjRTdpa1FqVFd2SVE1QWxHVDRXTHNnQVVHcVJzcDBBNzB4V3NSdwphMW9ZUWIzZngzWG0zZ0hwVXVyVGtjSVF2Q0xkS2NPZTlNbnB3UVBtRCtpK3V3THh5LzJjN004VzZaa2dXTlZJCmU3and6bW11SzNiWmI4ZmpBS0NHRDB4dkJOU0Y4N1hZZU9qN0d3Ykd3bjVPMFFYaXJ2dE1lZFZ6c1NuUTA2ZS8KNFVrPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="

func Test_FindByStage(t *testing.T) {
	now := time.Now().UTC()
	secretMetas := []core.SecretMeta{
		core.SecretMeta{AppName: "myapp", StageName: "dev", SecretName: "mysecret", Doc: core.SecretMetaDoc{LastUpdated: now}},
	}
	secretMetaRepository := &core.FakeSecretMetaRepository{
		FindByStageFn: func(appName string, stageName string) ([]core.SecretMeta, error) {
			assert.Equal(t, "myapp", appName)
			assert.Equal(t, "dev", stageName)
			return secretMetas, nil
		},
	}

	service := service{secretMetas: secretMetaRepository}

	result, err := service.FindByStage("myapp", "dev")

	assert.NoError(t, err)
	assert.Equal(t, secretMetas, result)
}

func Test_FindByStage_SecretRepositoryErr_ReturnsErr(t *testing.T) {
	expectedErr := errors.New("test")
	secretMetaRepository := &core.FakeSecretMetaRepository{
		FindByStageFn: func(string, string) ([]core.SecretMeta, error) {
			return nil, expectedErr
		},
	}

	service := service{secretMetas: secretMetaRepository}

	result, err := service.FindByStage("myapp", "dev")

	assert.Nil(t, result)
	assert.Equal(t, "Error retrieving secret metadata: test", err.Error())
}

func Test_FindNamesByStage(t *testing.T) {
	secretMetas := []core.SecretMeta{
		core.SecretMeta{SecretName: "mysecret"},
	}
	secretMetaRepository := &core.FakeSecretMetaRepository{
		FindByStageFn: func(string, string) ([]core.SecretMeta, error) {
			return secretMetas, nil
		},
	}

	service := service{secretMetas: secretMetaRepository}

	result, err := service.FindNamesByStage("myapp", "dev")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "mysecret", result[0])
}

func Test_SealAndSave(t *testing.T) {
	testCertBytes, _ := base64.StdEncoding.DecodeString(testCert)
	stageRepository := &core.FakeStageRepository{
		GetFn: func(stageName string) (*core.Stage, error) {
			assert.Equal(t, "mystage", stageName)
			return &core.Stage{
				Name: "mystage",
				Doc: core.StageDoc{
					Config: core.StageConfig{
						SealedSecretCert: testCertBytes,
					},
				},
			}, nil
		},
	}

	secretMetaRepository := &core.FakeSecretMetaRepository{
		SaveFn: func(secretMeta *core.SecretMeta) error {
			assert.Equal(t, "myapp", secretMeta.AppName)
			assert.Equal(t, "mystage", secretMeta.StageName)
			assert.Equal(t, "mysecret", secretMeta.SecretName)
			assert.InDelta(t, time.Now().UTC().Unix(), secretMeta.Doc.LastUpdated.Unix(), 3)
			return nil
		},
	}

	meta := &core.SecretMeta{
		AppName:    "myapp",
		StageName:  "mystage",
		SecretName: "mysecret",
	}

	comitter := state.NewDryRunComitter()

	service := service{secretMetas: secretMetaRepository, stages: stageRepository}

	result := service.SealAndSave("plain", meta, "myns", comitter)

	assert.NoError(t, result)
	assert.Equal(t, 1, stageRepository.GetCallCount)
	assert.Equal(t, 1, secretMetaRepository.SaveCallCount)
	assert.Len(t, comitter.Commits, 1)
	assert.Equal(t, "Updating secret \"myapp-mysecret\" in stage \"mystage\"", comitter.Commits[0].Message)
	assert.Len(t, comitter.Commits[0].Files, 1)
	assert.Equal(t, "stages/mystage/kube-resources/riser-managed/myns/secrets/myapp/sealedsecret.myapp-mysecret.yaml", comitter.Commits[0].Files[0].Name)
}
