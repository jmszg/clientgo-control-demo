package main

import (
	"clientgo-crd-demo/pkg"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 获取config
	// 先尝试从集群外部获取，获取不到则从集群内部获取
	var config, err = clientcmd.BuildConfigFromFlags("", "./config")
	if err != nil {
		clusterConfig, err := rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
		config = clusterConfig
	}

	// 通过config创建 clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 通过 client 创建 informer,添加事件处理函数
	factory := informers.NewSharedInformerFactory(clientSet, 0)
	serviceInformer := factory.Core().V1().Services()
	ingressInformer := factory.Networking().V1().Ingresses()
	newController := pkg.NewController(clientSet, serviceInformer, ingressInformer)

	// 启动 informer
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)
	newController.Run(stopCh)
}
