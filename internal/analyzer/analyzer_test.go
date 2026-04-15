package analyzer

import (
	"context"
	"strings"
	"testing"

	"github.com/kubeopencode/kubeopencode/internal/k8s"
)

func TestNew(t *testing.T) {
	a := New("test-api-key")
	if a == nil {
		t.Fatal("expected non-nil analyzer")
	}
	if a.model == "" {
		t.Fatal("expected non-empty model")
	}
}

func TestNewWithModel(t *testing.T) {
	model := "gpt-4"
	a := NewWithModel("test-api-key", model)
	if a.model != model {
		t.Errorf("expected model %q, got %q", model, a.model)
	}
}

func TestAnalyze_InvalidResource(t *testing.T) {
	a := New("invalid-key")
	ctx := context.Background()

	resource := &k8s.Resource{
		Kind:      "pod",
		Name:      "test-pod",
		Namespace: "default",
		Raw: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"name":      "test-pod",
				"namespace": "default",
			},
		},
	}

	_, err := a.Analyze(ctx, resource)
	// An invalid API key should always result in an authentication error from the upstream API.
	// TODO: mock the HTTP client here so this test doesn't require network access.
	if err == nil {
		t.Fatal("expected error with invalid API key")
	}
}

func TestResource_ToJSON(t *testing.T) {
	r := &k8s.Resource{
		Kind: "pod",
		Name: "test",
		Raw: map[string]interface{}{
			"kind": "Pod",
			"metadata": map[string]interface{}{
				"name": "test",
			},
		},
	}

	jsonStr, err := r.ToJSON()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if jsonStr == "" {
		t.Fatal("expected non-empty JSON")
	}

	// Sanity check: the JSON output should at least contain the resource name.
	if !strings.Contains(jsonStr, "test") {
		t.Errorf("expected JSON to contain resource name 'test', got: %s", jsonStr)
	}
}
