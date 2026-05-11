package analyzer

import (
	"testing"
	"time"
)

func TestNewAnalysisCache(t *testing.T) {
	c := NewAnalysisCache(5 * time.Minute)
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
	if c.Size() != 0 {
		t.Errorf("expected empty cache, got size %d", c.Size())
	}
}

func TestCache_SetAndGet(t *testing.T) {
	c := NewAnalysisCache(5 * time.Minute)
	c.Set("Pod", "default", "my-pod", "analysis result")

	result, ok := c.Get("Pod", "default", "my-pod")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if result != "analysis result" {
		t.Errorf("expected 'analysis result', got %q", result)
	}
}

func TestCache_Miss(t *testing.T) {
	c := NewAnalysisCache(5 * time.Minute)

	_, ok := c.Get("Pod", "default", "nonexistent")
	if ok {
		t.Error("expected cache miss for nonexistent key")
	}
}

func TestCache_Expiry(t *testing.T) {
	c := NewAnalysisCache(10 * time.Millisecond)
	c.Set("Deployment", "kube-system", "coredns", "some result")

	time.Sleep(20 * time.Millisecond)

	_, ok := c.Get("Deployment", "kube-system", "coredns")
	if ok {
		t.Error("expected cache miss after TTL expiry")
	}
}

func TestCache_Invalidate(t *testing.T) {
	c := NewAnalysisCache(5 * time.Minute)
	c.Set("Service", "default", "my-svc", "result")
	c.Invalidate("Service", "default", "my-svc")

	_, ok := c.Get("Service", "default", "my-svc")
	if ok {
		t.Error("expected cache miss after invalidation")
	}
}

func TestCache_Flush(t *testing.T) {
	c := NewAnalysisCache(5 * time.Minute)
	c.Set("Pod", "default", "pod-a", "r1")
	c.Set("Pod", "default", "pod-b", "r2")
	c.Flush()

	if c.Size() != 0 {
		t.Errorf("expected empty cache after flush, got size %d", c.Size())
	}
}

func TestCacheKey_Uniqueness(t *testing.T) {
	k1 := cacheKey("Pod", "default", "my-pod")
	k2 := cacheKey("Pod", "kube-system", "my-pod")
	k3 := cacheKey("Deployment", "default", "my-pod")

	if k1 == k2 || k1 == k3 || k2 == k3 {
		t.Error("expected unique cache keys for different inputs")
	}
}
