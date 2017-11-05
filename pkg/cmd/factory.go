package cmd

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Factory provides what is effectively injection for the commands
type Factory interface {
	// Clientset returns the interface to the kubernetes API clients
	Clientset() (kubernetes.Interface, error)

	// ConfigAccess returns the ConfigAccess interface for working with kubeconfig files
	ConfigAccess() (clientcmd.ConfigAccess, error)
}
