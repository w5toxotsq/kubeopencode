package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestSupportedResourceTypes(t *testing.T) {
	// Ensure SupportedResourceTypes is not empty
	assert.NotEmpty(t, SupportedResourceTypes, "SupportedResourceTypes should not be empty")

	// Verify common Kubernetes resource types are present
	expectedTypes := []string{"pod", "deployment", "service", "configmap"}
	for _, expected := range expectedTypes {
		assert.Contains(t, SupportedResourceTypes, expected,
			"SupportedResourceTypes should contain %q", expected)
	}
}

func TestResourceType_IsSupported(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		expected     bool
	}{
		{
			name:         "supported resource - pod",
			resourceType: "pod",
			expected:     true,
		},
		{
			name:         "supported resource - deployment",
			resourceType: "deployment",
			expected:     true,
		},
		{
			name:         "unsupported resource - unknown",
			resourceType: "unknownresource",
			expected:     false,
		},
		{
			name:         "case sensitivity - uppercase Pod",
			resourceType: "Pod",
			expected:     false,
		},
		{
			name:         "empty string",
			resourceType: "",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := SupportedResourceTypes[tt.resourceType]
			assert.Equal(t, tt.expected, ok)
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name: "k8s not found error",
			err: k8serrors.NewNotFound(
				schema.GroupResource{Group: "", Resource: "pods"},
				"test-pod",
			),
			expected: true,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "generic error",
			err:      assert.AnError,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := k8serrors.IsNotFound(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResourceMetadata(t *testing.T) {
	// Test that resource metadata is properly structured
	meta := metav1.ObjectMeta{
		Name:      "test-resource",
		Namespace: "default",
		Labels: map[string]string{
			"app": "test",
		},
	}

	require.Equal(t, "test-resource", meta.Name)
	require.Equal(t, "default", meta.Namespace)
	assert.Equal(t, "test", meta.Labels["app"])
}
