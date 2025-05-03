package k8s

import (
	"os"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientset() *kubernetes.Clientset {
	config := getKubeConfig()
	return createKubeClient(config)
}

func getKubeConfig() *rest.Config {
	log.Info("Attempting to use in cluster config")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Warn("In cluster config failed. Attempting to use kube config")
		kubeconfig := os.Getenv("KUBECONFIG")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Error("Failed to load kubeconfig", "error", err)
			os.Exit(1)
		}
	}

	return config
}

func createKubeClient(config *rest.Config) *kubernetes.Clientset {
	log.Info("Creating k8s client")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error("Failed to create Kubernetes client", "error", err)
		os.Exit(1)
	}

	return clientset
}