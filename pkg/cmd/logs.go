package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type LogsOptions struct {
	Name string

	Follow bool
}

func RunLogs(ctx context.Context, f Factory, out io.Writer, o *LogsOptions) error {
	if len(o.Name) == 0 {
		return fmt.Errorf("you must specify a non-empty name")
	}

	clientset, err := f.Clientset()
	if err != nil {
		return err
	}

	configAccess, err := f.ConfigAccess()
	if err != nil {
		return err
	}

	configs, err := configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	if configs.CurrentContext == "" {
		return fmt.Errorf("current-context not set")
	}

	context := configs.Contexts[configs.CurrentContext]
	if context == nil {
		return fmt.Errorf("context %q not found", configs.CurrentContext)
	}

	namespace := context.Namespace
	if namespace == "" {
		namespace = "default"
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(o.Name, meta_v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			pod = nil
		} else {
			return fmt.Errorf("error getting pod %q: %v", o.Name, err)
		}
	}

	podLogOptions := &v1.PodLogOptions{
		Follow: o.Follow,
	}

	if pod != nil {
		logsRequest := clientset.CoreV1().Pods(namespace).GetLogs(o.Name, podLogOptions)
		if err != nil {
			return fmt.Errorf("error getting logs for pod %q: %v", o.Name, err)
		}

		prefix := ""
		return pipeLogs(ctx, logsRequest, prefix, out)
	}

	// Find matching pods
	pods, err := clientset.CoreV1().Pods(namespace).List(meta_v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing pods: %v", err)
	}

	var wg sync.WaitGroup
	count := 0
	for _, pod := range pods.Items {
		if !strings.HasPrefix(pod.Name, o.Name) {
			glog.V(4).Infof("rejecting pod %q", pod.Name)
			continue
		}

		count++
		logsRequest := clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, podLogOptions)
		if err != nil {
			// TODO: Don't bail if a pod exits quickly
			return fmt.Errorf("error getting logs for pod %q: %v", pod.Name, err)
		}

		wg.Add(1)

		prefix := pod.Name + ": "
		go func(prefix string) {
			_ = pipeLogs(ctx, logsRequest, prefix, out)
			wg.Done()
		}(prefix)
	}
	wg.Wait()

	if count == 0 {
		return fmt.Errorf("found no pods matching %q", o.Name)
	}

	return nil
}

func pipeLogs(ctx context.Context, req *rest.Request, prefix string, out io.Writer) error {
	readCloser, err := req.Stream()
	if err != nil {
		return err
	}
	defer readCloser.Close()

	scanner := bufio.NewScanner(readCloser)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		// TODO: Do we want to use context here?  And if so, we should force close the channel on context.Done() instead of this
		if ctx.Err() != nil {
			return ctx.Err()
		}
		line := scanner.Text()
		if prefix != "" {
			line = prefix + line
		}
		line += "\n"

		// TODO: Lock on out?
		// TODO: Push logs into channel or something or otherwise batch writes?
		_, err := out.Write([]byte(line))
		if err != nil {
			return err
		}
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	return nil
}
