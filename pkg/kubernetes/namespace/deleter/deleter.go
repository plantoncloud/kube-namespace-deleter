package deleter

import (
	"context"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

type Deleter interface {
	GetNamespace(namespaceName string) (*v1.Namespace, error)
	RemoveFinalizers(namespace *v1.Namespace) error
	DeleteNamespace(namespace *v1.Namespace) error
}

type RealDeleter struct{}

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

func (e RealDeleter) RemoveFinalizers(namespace *v1.Namespace) error {
	// Use the current context in kube-config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "failed to get user home-directory")
	}

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homeDir, ".kube", "config"))
	if err != nil {
		return errors.Wrap(err, "failed to build config from flags")
	}

	// Create a Kubernetes client
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create client-set")
	}

	// Remove the finalizers
	namespace.SetFinalizers([]string{})

	// Update the namespace
	_, err = clientSet.CoreV1().Namespaces().Update(context.TODO(), namespace, metav1.UpdateOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to update namespace")
	}

	return nil
}

func (e RealDeleter) DeleteNamespace(namespace *v1.Namespace) error {
	// Use the current context in kube-config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "failed to get user home-directory")
	}

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homeDir, ".kube", "config"))
	if err != nil {
		return errors.Wrap(err, "failed to build config from flags")
	}

	// Create a Kubernetes client
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create client-set")
	}

	// Remove the finalizers
	namespace.SetFinalizers([]string{})

	// Update the namespace
	if err = clientSet.CoreV1().Namespaces().Delete(context.TODO(), namespace.Name, metav1.DeleteOptions{}); err != nil {
		return errors.Wrap(err, "failed to update namespace")
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

	if err := executor.DeleteNamespace(n); err != nil {
		return errors.Wrapf(err, "failed to delete namespace %s", namespace)
	}
	return nil
}
