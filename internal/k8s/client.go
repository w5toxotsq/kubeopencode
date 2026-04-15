package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

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
	gvr, ok := resourceMap[kind]
	if !ok {
		return nil, fmt.Errorf("unsupported resource type: %s", kind)
	}

	obj, err := c.dynamic.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get %s/%s: %w", kind, name, err)
	}

	raw := obj.Object
	return &Resource{
		Kind:      kind,
		Name:      name,
		Namespace: namespace,
		Raw:       raw,
	}, nil
}

func (r *Resource) ToJSON() (string, error) {
	b, err := json.MarshalIndent(r.Raw, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
