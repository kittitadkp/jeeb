package kube

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/ptr"

	"k8s-manager/internal/logger"
	"k8s-manager/internal/printer"
)

var jeebNamespaces = []string{"jeeb-dev", "jeeb-infra", "jeeb-obs"}

type Client struct {
	cs kubernetes.Interface
}

func NewClient(kubeconfig string) (*Client, error) {
	if kubeconfig == "" {
		home, _ := os.UserHomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("build kubeconfig: %w", err)
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("create clientset: %w", err)
	}

	return &Client{cs: cs}, nil
}

func (c *Client) CreateNamespaces(ctx context.Context, namespaces []string) error {
	for _, ns := range namespaces {
		_, err := c.cs.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: ns},
		}, metav1.CreateOptions{})
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				logger.Debug("namespace %q already exists", ns)
				continue
			}
			return fmt.Errorf("create namespace %s: %w", ns, err)
		}
		logger.Info("namespace %q created", ns)
	}
	return nil
}

func (c *Client) PrintStatus(ctx context.Context, namespace string) error {
	namespaces := jeebNamespaces
	if namespace != "" {
		namespaces = []string{namespace}
	}

	for _, ns := range namespaces {
		pods, err := c.cs.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("list pods in %s: %w", ns, err)
		}
		printer.PrintPods(ns, pods.Items)
	}
	return nil
}

func (c *Client) RestartDeployment(ctx context.Context, namespace, name string) error {
	dep, err := c.cs.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get deployment %s/%s: %w", namespace, name, err)
	}

	if dep.Spec.Template.Annotations == nil {
		dep.Spec.Template.Annotations = map[string]string{}
	}
	dep.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = c.cs.AppsV1().Deployments(namespace).Update(ctx, dep, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("update deployment: %w", err)
	}

	logger.Info("restarted deployment %s/%s", namespace, name)
	return nil
}

func (c *Client) StreamLogs(ctx context.Context, namespace, deployment string, follow bool, tail int64) error {
	pods, err := c.cs.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deployment),
	})
	if err != nil {
		return fmt.Errorf("list pods: %w", err)
	}

	if len(pods.Items) == 0 {
		return fmt.Errorf("no pods found for deployment %s in %s", deployment, namespace)
	}

	pod := pods.Items[0]
	opts := &corev1.PodLogOptions{
		Follow:    follow,
		TailLines: ptr.To(tail),
	}

	req := c.cs.CoreV1().Pods(namespace).GetLogs(pod.Name, opts)
	stream, err := req.Stream(ctx)
	if err != nil {
		return fmt.Errorf("open log stream: %w", err)
	}
	defer stream.Close()

	_, err = io.Copy(os.Stdout, stream)
	return err
}
