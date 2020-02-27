package resources

import (
	"testing"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSealedSecretCert = `-----BEGIN CERTIFICATE-----
MIIErjCCApagAwIBAgIRAKVbCeiOLBFrQkTE0gDzz8owDQYJKoZIhvcNAQELBQAw
ADAeFw0xOTA3MTMwODI3MzFaFw0yOTA3MTAwODI3MzFaMAAwggIiMA0GCSqGSIb3
DQEBAQUAA4ICDwAwggIKAoICAQDr/b56OJJZjTp7yAY3JpjWa77RxZ9hiGR3ffPG
PjhUPPkj1FjV/y+krhGCvHkDhEPg8ccNWxEEz959hRCPKxw2t1UeGxqJeDtJ690t
IaY/h0tSQqKr5neE2TXdGtsMciVAwBHnbl5xX0UFzqhSDMmeraaoDQEdbbe/I4ym
fZ1okHYXjySFXOggBmZ63YD7DkpIV6/Mu2cZqgkWNtvfYe3zpbZBm8kyIHLw5Dk0
5GqoC2xdWMTMiKb1k90qQY3TRkXceTeH6v4uxMTZEbpfznZ2tcU3c1ooxSMxI7gO
LZ7HtCW3MDxVtEjvJWnq9xUFioS1Hq9YReDJ9a4fgL9QeZLTVSL+oeAvBhunvAgd
zfEcrgYIAULoRZVWwEGQDe5MPWEfdoQAEybpjxCxnjon+HckLoq7zbrTyCwcvADo
VRqger2Bm2mu7m3vHpmQO1xTFomzF/73Q84Dstdutd60RJbGpg7k2oODZgzQ9vR8
9Ybm2qB0NIIk1/FXi9SUmP61B9vtxWffUu22V9MUaYSW9NCtGA4t+AGN6GfXdVsg
LLz4I4HVWcs1yOzrzTGcIF5lSsNBQMJyU7C4z8Izhzs2+IfFig8Zq5PDX31nDNxE
zxf7XYhvIKmK09BlUaAakW9kTLPjdHJySEbKs1MU6G9pmmfxnGCH3CnRsCRt/c4f
nITOqQIDAQABoyMwITAOBgNVHQ8BAf8EBAMCAAEwDwYDVR0TAQH/BAUwAwEB/zAN
BgkqhkiG9w0BAQsFAAOCAgEA6K3uQjddhVwqsEOkdEmqKHRWxoMsTKWAvhyOXBy9
Czl/8/F3Rb/rTMoTbM5xmmLSKidycezY4M275GmgdXd0Y+ygXcukznrk4wFqwkDE
41Tm+k0B6KHHOMwMVo0HB5JvadOeDUt1TFuFbN8JNhah2h8Nx7piTBPbiTYo0Cg5
uuPLtSAsKocs/PszFtbJRHfBAFa28xnhuIFH+Lsguc0AQHYmqdDiOZao3aM2Kh4+
n4Z8fcMGFVHPQ3sbZJOWrxD6WYSlHN333kOGblU31pIOZFnXH13mdUdE5uVxvL1i
6e0cVnPdNvg1uTbW2rfXiC69rtlq69LZprzBqGhulJSgAXH5JFpLU/hbbks4kmWd
M9LqlnuhqI6jMbB54TNfxKwJSkBoqjbBZ7e5FUqbKkYFKH9PnRuG0O8J/XTFDdCv
4CuHIjs5D/hDjShK22w/jroRKwPZs8g3XtlAzrDP7xqz6hDBMs8BbyvD1Oc4YzWI
hoqkLmoUCsH7mleLl8n23+tzi51sRp7L3MKNHiwyzPp6nSxUfa0iDeHO8baF+E0D
ip4BNt2piGOeDlqJOaBL06yUHPuxsxaEbEyndL7KQjMKUuIKUVz0OZ4Br5wB8l0x
PNijHhubYJ4Tg9+qeIIP55U4lToTjIxh0gFJqwBoDY3ta054LKc4s0hE0FmwtMVQ
xbU=
-----END CERTIFICATE-----`

func Test_CreateSealedSecret(t *testing.T) {
	secret := &core.SecretMeta{
		Name:      "mysecretname",
		StageName: "dev",
		Revision:  1,
	}

	result, err := CreateSealedSecret("mysecretvalue", "myapp", secret, "apps", []byte(testSealedSecretCert))

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "myapp-mysecretname-1", result.Name)
	assert.Equal(t, "apps", result.Namespace)
	assert.Equal(t, "SealedSecret", result.Kind)
	assert.Equal(t, "bitnami.com/v1alpha1", result.APIVersion)
	assert.Equal(t, "1", result.Annotations["riser.dev/revision"])
	assert.Equal(t, "myapp", result.Labels["riser.dev/app"])
	// Sanity check that we're setting the encrypted data field.
	// We'll use e2e integration testing that tests the mounting of secrets into a pod for better coverage.
	assert.Len(t, result.Spec.EncryptedData, 1)
	assert.NotEmpty(t, result.Spec.EncryptedData["data"])
}

func Test_parsePublicKey(t *testing.T) {
	result, err := parsePublicKey([]byte(testSealedSecretCert))

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 512, result.Size())
}

func Test_parsePublicKey_badCert(t *testing.T) {
	result, err := parsePublicKey([]byte("notacert"))

	require.Nil(t, result)
	assert.Error(t, err)
}
