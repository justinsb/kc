package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SshOptions struct {
	NodeName string

	// Username is the SSH username to use when SSH-ing to the node
	Username string
}

func RunSsh(ctx context.Context, f Factory, out io.Writer, o *SshOptions) error {
	if len(o.NodeName) == 0 {
		return fmt.Errorf("must specify NodeName")
	}

	clientset, err := f.Clientset()
	if err != nil {
		return err
	}

	node, err := clientset.CoreV1().Nodes().Get(o.NodeName, meta_v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("node %s not found", err)
		} else {
			return fmt.Errorf("error getting node %s: %v", o.NodeName, err)
		}
	}

	var externalIP []string
	var externalDNS []string

	for _, address := range node.Status.Addresses {
		switch address.Type {
		case v1.NodeExternalIP:
			externalIP = append(externalIP, address.Address)
		case v1.NodeExternalDNS:
			externalDNS = append(externalDNS, address.Address)
		}
	}

	best := ""
	if len(externalIP) != 0 {
		best = externalIP[0]
	} else if len(externalDNS) != 0 {
		best = externalDNS[0]
	}

	if best == "" {
		return fmt.Errorf("cannot find an external IP or DNS address for node %s", o.NodeName)
	}

	var args []string
	if o.Username != "" {
		best = o.Username + "@" + best
	}

	args = append(args, best)
	c := exec.Command("ssh", args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	return c.Run()
}
