package k8s

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientset() (*kubernetes.Clientset, error) {
	log.Info("Attempting to use in cluster config")
	config, err := rest.InClusterConfig()
	if err == nil {
		log.Info("Successfully got in-cluster kube config")
		return createClientset(config), nil
	}

	log.Warn("In cluster config failed. Attempting to use kube config")
	kubeConfigEnvPath := os.Getenv("KUBECONFIG")
	config, err = clientcmd.BuildConfigFromFlags("", kubeConfigEnvPath)
	if err == nil {
		log.Info("Successfully got env var kube config")
		return createClientset(config), nil
	}

	log.Warn("Kube config env var failed. Attempting to use ~/.kube/config")
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("Could not get user home directory: %v", err)
		panic(err)
	}

	kubeConfigUserPath := filepath.Join(userHome, ".kube", "config")
	config, err = clientcmd.BuildConfigFromFlags("", kubeConfigUserPath)
	if err == nil {
		log.Info("Successfully got user home kube config")
		return createClientset(config), nil
	}

	return nil, err
}

func createClientset(config *rest.Config) *kubernetes.Clientset {
	log.Info("Creating k8s client")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error("Failed to create Kubernetes client", "error", err)
		os.Exit(1)
	}

	return clientset
}