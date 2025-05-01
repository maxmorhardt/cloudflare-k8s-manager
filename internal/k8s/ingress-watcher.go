package k8s

import (
	"log/slog"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func Run(logger *slog.Logger) {
	watchIngress(logger)
}

func watchIngress(logger *slog.Logger) {
	config := getKubeConfig(logger)
	clientset := createKubeClient(logger, config)
	watchList := createIngressListWatch(logger, clientset)
	controller := createInformer(logger, watchList)

	stopCh := make(chan struct{})
	defer close(stopCh)

	logger.Info("Listening to ingress resources")
	controller.Run(stopCh)
}

func getKubeConfig(logger *slog.Logger) *rest.Config {
	logger.Info("Attempting to use in cluster config")
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Warn("In cluster config failed. Attempting to use kube config")
		kubeconfig := os.Getenv("KUBECONFIG")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			logger.Error("Failed to load kubeconfig", "error", err)
			os.Exit(1)
		}
	}

	return config
}

func createKubeClient(logger *slog.Logger, config *rest.Config) *kubernetes.Clientset {
	logger.Info("Creating k8s client")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error("Failed to create Kubernetes client", "error", err)
		os.Exit(1)
	}

	return clientset
}

func createIngressListWatch(logger *slog.Logger, clientset *kubernetes.Clientset) *cache.ListWatch {
	logger.Info("Creating list watch for ingress in all namespaces")
	watchList := cache.NewListWatchFromClient(
		clientset.NetworkingV1().RESTClient(),
		"ingresses",
		metav1.NamespaceAll,
		fields.Everything(),
	)

	return watchList
}

func createInformer(logger *slog.Logger, watchList *cache.ListWatch) cache.Controller {
	options := cache.InformerOptions{
		ListerWatcher: watchList, 
		ObjectType: &networkingv1.Ingress{},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj any) {
				ingress := obj.(*networkingv1.Ingress)
				logger.Info("Ingress added or exists", "namespace", ingress.Namespace, "name", ingress.Name)
			},
			UpdateFunc: func(_, newObj any) {
				ingress := newObj.(*networkingv1.Ingress)
				logger.Info("Ingress updated", "namespace", ingress.Namespace, "name", ingress.Name)
			},
			DeleteFunc: func(obj any) {
				ingress := obj.(*networkingv1.Ingress)
				logger.Info("Ingress deleted", "namespace", ingress.Namespace, "name", ingress.Name)
			},
		},
		ResyncPeriod: 0,
		MinWatchTimeout: 0,
		Indexers: nil,
		Transform: nil,
	}

	_, controller := cache.NewInformerWithOptions(options)

	return controller
}