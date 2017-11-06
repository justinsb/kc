package cmd

import (
	"fmt"
	"strings"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CompleteNamespaces(f Factory, prefix string) ([]string, error) {
	clientset, err := f.Clientset()
	if err != nil {
		return nil, err
	}

	list, err := clientset.CoreV1().Namespaces().List(meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var matches []string
	for _, ns := range list.Items {
		if strings.HasPrefix(ns.Name, prefix) {
			matches = append(matches, ns.Name)
		}
	}

	return matches, nil
}

func CompletePods(f Factory, prefix string) ([]string, error) {
	clientset, err := f.Clientset()
	if err != nil {
		return nil, err
	}

	namespace, err := getNamespace(f)
	if err != nil {
		return nil, err
	}

	// Find matching pods
	pods, err := clientset.CoreV1().Pods(namespace).List(meta_v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing pods: %v", err)
	}

	var matches []string
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, prefix) {
			matches = append(matches, pod.Name)
		}
	}

	return matches, nil
}
