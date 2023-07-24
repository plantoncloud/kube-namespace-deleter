package deleter

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type Deleter interface {
	RunCommand(cmd *exec.Cmd) error
	GetNamespace(namespaceName string) (*v1.Namespace, error)
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte) error
	RemoveFile(filename string) error
	RemoveFinalizers(namespace *v1.Namespace) error
}

type RealDeleter struct{}

func (e RealDeleter) RunCommand(cmd *exec.Cmd) error {
	return cmd.Run()
}

func (e RealDeleter) GetNamespace(namespaceName string) (*v1.Namespace, error) {
	// Use the current context in kube-config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user home-directory")
	}

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homeDir, ".kube", "config"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to build config from flags")
	}

	// Create a Kubernetes client
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client-set")
	}

	// Get the namespace
	namespace, err := clientSet.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get namespace %s", namespaceName)
	}

	return namespace, nil
}

func (e RealDeleter) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (e RealDeleter) WriteFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

func (e RealDeleter) RemoveFile(filename string) error {
	return os.Remove(filename)
}

func (e RealDeleter) RemoveFinalizers(namespace *v1.Namespace) error {
	namespace.SetFinalizers([]string{})
	// Create an HTTP client
	client := &http.Client{}

	namespaceJsonBytes, err := json.Marshal(namespace)
	if err != nil {
		return errors.Wrap(err, "failed to marshal namespace to JSON")
	}

	// Create the HTTP request
	req, err := http.NewRequest("PUT", "http://127.0.0.1:8001/api/v1/namespaces/"+namespace.Name+"/finalize", bytes.NewBuffer(namespaceJsonBytes))
	if err != nil {
		return errors.Wrap(err, "failed to create HTTP request")
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send HTTP request")
	}
	defer resp.Body.Close()

	// Check the HTTP response status code
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("received non-OK HTTP status: %s", resp.Status)
	}
	return nil
}

func Delete(namespace string, executor Deleter) error {
	n, err := executor.GetNamespace(namespace)
	if err != nil {
		return errors.Wrapf(err, "failed to get namespace %s", namespace)
	}

	if err := executor.RemoveFinalizers(n); err != nil {
		return errors.Wrapf(err, "failed to remove finalizers from namespace %s", namespace)
	}
	return nil
}
