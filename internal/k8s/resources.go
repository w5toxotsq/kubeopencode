package k8s

import (
	"context"
	"fmt"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ResourceType represents a Kubernetes resource type
type ResourceType string

const (
	ResourceTypePod         ResourceType = "pods"
	ResourceTypeDeployment  ResourceType = "deployments"
	ResourceTypeService     ResourceType = "services"
	ResourceTypeConfigMap   ResourceType = "configmaps"
	ResourceTypeSecret      ResourceType = "secrets"
	ResourceTypeStatefulSet ResourceType = "statefulsets"
	ResourceTypeDaemonSet   ResourceType = "daemonsets"
	// Added Job support for my batch workload use case
	ResourceTypeJob ResourceType = "jobs"
)

// groupVersionResources maps resource types to their GVR
var groupVersionResources = map[ResourceType]schema.GroupVersionResource{
	ResourceTypePod:         {Group: "", Version: "v1", Resource: "pods"},
	ResourceTypeDeployment:  {Group: "apps", Version: "v1", Resource: "deployments"},
	ResourceTypeService:     {Group: "", Version: "v1", Resource: "services"},
	ResourceTypeConfigMap:   {Group: "", Version: "v1", Resource: "configmaps"},
	ResourceTypeSecret:      {Group: "", Version: "v1", Resource: "secrets"},
	ResourceTypeStatefulSet: {Group: "apps", Version: "v1", Resource: "statefulsets"},
	ResourceTypeDaemonSet:   {Group: "apps", Version: "v1", Resource: "daemonsets"},
	ResourceTypeJob:         {Group: "batch", Version: "v1", Resource: "jobs"},
}

// GetResource retrieves a single Kubernetes resource by name and namespace.
func (c *Client) GetResource(ctx context.Context, resourceType ResourceType, namespace, name string) (*unstructured.Unstructured, error) {
	gvr, ok := groupVersionResources[resourceType]
	if !ok {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	resource, err := c.dynamic.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get %s/%s in namespace %s: %w", resourceType, name, namespace, err)
	}

	return resource, nil
}

// ListResources retrieves all Kubernetes resources of a given type in a namespace.
// If namespace is empty, it lists resources across all namespaces.
func (c *Client) ListResources(ctx context.Context, resourceType ResourceType, namespace string) ([]unstructured.Unstructured, error) {
	gvr, ok := groupVersionResources[resourceType]
	if !ok {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	var list *unstructured.UnstructuredList
	var err error

	if namespace == "" {
		list, err = c.dynamic.Resource(gvr).List(ctx, metav1.ListOptions{})
	} else {
		list, err = c.dynamic.Resource(gvr).Namespace(namespace).List(ctx, metav1.ListOptions{})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list %s in namespace %q: %w", resourceType, namespace, err)
	}

	return list.Items, nil
}

// SupportedResourceTypes returns the list of resource types supported for analysis.
// Returns types in sorted order for consistent, predictable output.
func SupportedResourceTypes() []ResourceType {
	types := make([]ResourceType, 0, len(groupVersionResources))
	for k := range groupVersionResources {
		types = append(types, k)
	}
	sort.Slice(types, func(i, j int) bool {
		return types[i] < types[j]
	})
	return types
}
