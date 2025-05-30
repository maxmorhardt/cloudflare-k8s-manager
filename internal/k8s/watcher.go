package k8s

import (
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func Watcher() {
	clientset, err := GetClientset()
	if err != nil {
		log.Error()
		panic(err)
	}

	ingressWatchList := createIngressListWatch(clientset)
	serviceWatchList := createServiceListWatch(clientset)

	stopCh := make(chan struct{})
	defer close(stopCh)

	go watchIngresses(ingressWatchList, stopCh)
	go watchServices(serviceWatchList, stopCh)

	select {}
}

func createIngressListWatch(clientset *kubernetes.Clientset) *cache.ListWatch {
	log.Info("Creating list watch for ingresses in all namespaces")
	watchList := cache.NewListWatchFromClient(
		clientset.NetworkingV1().RESTClient(),
		"ingresses",
		metav1.NamespaceAll,
		fields.Everything(),
	)

	return watchList
}

func createServiceListWatch(clientset *kubernetes.Clientset) *cache.ListWatch {
	log.Info("Creating list watch for services in all namespaces")
	watchList := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		"services",
		metav1.NamespaceAll,
		fields.Everything(),
	)

	return watchList
}

func watchIngresses(watchList *cache.ListWatch, stopCh <-chan struct{}) {
	options := cache.InformerOptions{
		ListerWatcher: watchList, 
		ObjectType: &networkingv1.Ingress{},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj any) {
				ing := obj.(*networkingv1.Ingress)
				log.Infof("Ingress added or exists namespace=%v name=%v", ing.Namespace, ing.Name)
			},
			UpdateFunc: func(_, newObj any) {
				ing := newObj.(*networkingv1.Ingress)
				log.Infof("Ingress updated namespace=%v name=%v", ing.Namespace, ing.Name)
			},
			DeleteFunc: func(obj any) {
				ing := obj.(*networkingv1.Ingress)
				log.Infof("Ingress deleted namespace=%v name=%v", ing.Namespace, ing.Name)
			},
		},
		ResyncPeriod: 0,
		MinWatchTimeout: 0,
		Indexers: nil,
		Transform: nil,
	}

	_, controller := cache.NewInformerWithOptions(options)

	controller.Run(stopCh)
}

func watchServices(watchList *cache.ListWatch, stopCh <-chan struct{}) {
	options := cache.InformerOptions{
		ListerWatcher: watchList, 
		ObjectType: &corev1.Service{},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj any) {
				svc := obj.(*corev1.Service)
				log.Infof("Service added or exists namespace=%v name=%v", svc.Namespace, svc.Name)
			},
			UpdateFunc: func(_, newObj any) {
				svc := newObj.(*corev1.Service)
				log.Infof("Service updated namespace=%v name=%v", svc.Namespace, svc.Name)
			},
			DeleteFunc: func(obj any) {
				svc := obj.(*corev1.Service)
				log.Infof("Service deleted namespace=%v name=%v", svc.Namespace, svc.Name)
			},
		},
		ResyncPeriod: 0,
		MinWatchTimeout: 0,
		Indexers: nil,
		Transform: nil,
	}

	_, controller := cache.NewInformerWithOptions(options)

	controller.Run(stopCh)
}
