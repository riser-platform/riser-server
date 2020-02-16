package resources

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/util/cert"

	sealedCrypto "github.com/bitnami-labs/sealed-secrets/pkg/crypto"
)

type SealedSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SealedSecretSpec `json:"spec"`
}

type SealedSecretSpec struct {
	EncryptedData map[string][]byte `json:"encryptedData"`
}

// TODO: Consider using something like https://github.com/awnumar/memguard instead of passing the secret as a string
func CreateSealedSecret(plaintextSecret string, secretMeta *core.SecretMeta, namespace string, certBytes []byte) (*SealedSecret, error) {
	publicKey, err := parsePublicKey(certBytes)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing public key")
	}
	objectMeta := metav1.ObjectMeta{
		// Note: For the moment secrets are shared between deployments in the same stage and namespace. Need to validate this requirement
		Name:      fmt.Sprintf("%s-%s", secretMeta.AppName, secretMeta.SecretName),
		Namespace: namespace,
		Annotations: map[string]string{
			riserLabel("revision"): fmt.Sprintf("%d", secretMeta.Revision),
		},
		Labels: map[string]string{
			riserLabel("app"): secretMeta.AppName,
		},
	}
	ciphertext, err := sealSecret(objectMeta, publicKey, []byte(plaintextSecret))
	if err != nil {
		return nil, errors.Wrap(err, "Error sealing secret")
	}
	return &SealedSecret{
		ObjectMeta: objectMeta,
		TypeMeta: metav1.TypeMeta{
			Kind:       "SealedSecret",
			APIVersion: "bitnami.com/v1alpha1",
		},
		Spec: SealedSecretSpec{
			EncryptedData: map[string][]byte{
				"data": ciphertext,
			},
		},
	}, nil
}

// Derived from https://github.com/bitnami-labs/sealed-secrets/blob/d875137740275f7dea36c54f981a90c795e7e681/cmd/kubeseal/main.go#L75
func parsePublicKey(certBytes []byte) (*rsa.PublicKey, error) {
	certs, err := cert.ParseCertsPEM(certBytes)
	if err != nil {
		return nil, err
	}

	// ParseCertsPem returns error if len(certs) == 0, but best to be sure...
	if len(certs) == 0 {
		return nil, errors.New("Failed to read any certificates")
	}

	cert, ok := certs[0].PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Expected RSA public key but found %v", certs[0].PublicKey)
	}

	return cert, nil
}

func sealSecret(secretMeta metav1.ObjectMeta, publicKey *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	// Simplified version of labelFor (https://github.com/bitnami-labs/sealed-secrets/blob/d875137740275f7dea36c54f981a90c795e7e681/pkg/apis/sealed-secrets/v1alpha1/sealedsecret_expansion.go#L22)
	// We don't support namespace or cluster wide annotations
	label := []byte(fmt.Sprintf("%s/%s", secretMeta.GetNamespace(), secretMeta.GetName()))
	return sealedCrypto.HybridEncrypt(rand.Reader, publicKey, plaintext, label)
}
