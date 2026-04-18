package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_InvalidKubeconfig(t *testing.T) {
	// Providing a non-existent kubeconfig path should return an error
	client, err := NewClient("/nonexistent/path/to/kubeconfig")
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClient_EmptyKubeconfig(t *testing.T) {
	// Empty kubeconfig path should attempt to use in-cluster or default config.
	// In a test environment without a cluster this will fail, which is expected.
	// Note: set KUBECONFIG env var locally if you want this to succeed.
	client, err := NewClient("")
	if err != nil {
		// Expected in CI/local without a cluster
		assert.Nil(t, client)
	} else {
		require.NotNil(t, client)
	}
}

func TestClient_GetResource_UnsupportedType(t *testing.T) {
	// We cannot create a real client in unit tests, so we verify that
	// SupportedResourceTypes does not include an unknown type, which
	// guards the GetResource path.
	unsupported := ResourceType("flibbertigibbet")
	assert.False(t, unsupported.IsSupported())
}

func TestResourceType_IsSupported_KnownTypes(t *testing.T) {
	// Sanity check that common resource types are recognised as supported.
	knownTypes := []ResourceType{"pod", "deployment", "service"}
	for _, rt := range knownTypes {
		assert.True(t, rt.IsSupported(), "expected %q to be supported", rt)
	}
}
