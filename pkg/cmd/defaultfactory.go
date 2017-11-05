package cmd

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// DefaultFactory providers the default implementation of Factory
type DefaultFactory struct {
	clientset kubernetes.Interface

	options *FactoryOptions
}

var _ Factory = &DefaultFactory{}

type FactoryOptions struct {
	Kubeconfig string
}

// Clientset implements Factory::Clientset
func (f *DefaultFactory) Clientset() (kubernetes.Interface, error) {
	if f.clientset == nil {
		kubeconfig := f.options.Kubeconfig

		if kubeconfig == "" {
			return nil, fmt.Errorf("kubeconfig path must be provided")
		}

		// use the current context in kubeconfig

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("error reading kubeconfig %q: %v", kubeconfig, err)
		}

		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, fmt.Errorf("error building client: %v", err)
		}

		f.clientset = client
	}

	return f.clientset, nil
}

// Clientset implements Factory::ConfigAccess
func (f *DefaultFactory) ConfigAccess() (clientcmd.ConfigAccess, error) {
	kubeconfig := f.options.Kubeconfig
	if kubeconfig == "" {
		return nil, fmt.Errorf("kubeconfig path must be provided")
	}

	c := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
	return c, nil
}

func NewDefaultFactory(options *FactoryOptions) Factory {
	f := &DefaultFactory{
		options: options,
	}
	return f
}
