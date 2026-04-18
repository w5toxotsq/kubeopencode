package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/apimachinery/pkg/runtime/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Client struct {
	dynamic dynamic.Interface
}

type Resource struct {
	Kind      string
	Name      string
	Namespace string
	Raw       map[string]interface{}
}

var resourceMap = map[string]schema.GroupVersionResource{
	"pod":         {Group: "", Version: "v1", Resource: "pods"},
	"deployment":  {Group: "apps", Version: "v1", Resource: "deployments"},
	"service":     {Group: "", Version: "v1", Resource: "services"},
	"configmap":   {Group: "", Version: "v1", Resource: "configmaps"},
	"statefulset": {Group: "apps", Version: "v1", Resource: "statefulsets"},
	// added ingress since I use it frequently in my homelab setup
	"ingress":     {Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"},
	// added daemonset for monitoring agents (node-exporter, etc.)
	"daemonset":   {Group: "apps", Version: "v1", Resource: "daemonsets"},
	// added secret - useful for quickly inspecting (non-sensitive) metadata
	"secret":      {Group: "", Version: "v1", Resource: "secrets"},
	// added namespace for cluster-wide queries
	"namespace":   {Group: "", Version: "v1", Resource: "namespaces"},
	// added job/cronjob - I run a bunch of batch workloads in my cluster
	"job":         {Group: "batch", Version: "v1", Resource: "jobs"},
	"cronjob":     {Group: "batch", Version: "v1", Resource: "cronjobs"},
	// added persistentvolumeclaim - handy for debugging storage issues in my homelab
	"persistentvolumeclaim": {Group: "", Version: "v1", Resource: "persistentvolumeclaims"},
	// shorthand aliases
	"pvc":         {Group: "", Version: "v1", Resource: "persistentvolumeclaims"},
	// added replicaset - occasionally useful when debugging deployment rollouts
	"replicaset":  {Group: "apps", Version: "v1", Resource: "replicasets"},
	"rs":          {Group: "apps", Version: "v1", Resource: "replicasets"},
	// added serviceaccount - useful when debugging RBAC issues
	"serviceaccount": {Group: "", Version: "v1", Resource: "serviceaccounts"},
	"sa":             {Group: "", Version: "v1", Resource: "serviceaccounts"},
}

func NewClient(kubeconfigPath string) (*Client, error) {
	if kubeconfigPath == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfigPath = filepath.Join(home, ".kube", "config")
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build config: %w", err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	return &Client{dynamic: dynClient}, nil
}

func (c *Client) GetResource(ctx context.Context, kind, name, namespace string) (*Resource, error) {
	// normalize kind to lowercase so callers don't have to worry about casing
	kind = strings.ToLower(kind)

	gvr, ok := resourceMap[kind]
	if !ok {
		return nil, fmt.Errorf("unsupported resource type: %s", kind)
	}

	obj, err := c.dynamic.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
