package main

import (
	"context"
	"fmt"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type controller struct {
	clientset             kubernetes.Interface
	deploymentLister      appslisters.DeploymentLister
	deploymentCacheSynced cache.InformerSynced
	queue                 workqueue.RateLimitingInterface
}

func main() {
	clientset := initClient()

	factory := informers.NewSharedInformerFactory(clientset, 10*time.Minute)

	deploymentInformer := factory.Apps().V1().Deployments()

	ctx := context.Background()
	c := newController(clientset, deploymentInformer)

	factory.Start(ctx.Done())

	c.run(ctx)
}

func newController(clientset kubernetes.Interface, deploymentInformer appsinformers.DeploymentInformer) *controller {
	c := &controller{
		clientset:             clientset,
		deploymentLister:      deploymentInformer.Lister(),
		deploymentCacheSynced: deploymentInformer.Informer().HasSynced, //注册缓存同步信息
		queue:                 workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "my-deployment"),
	}

	deploymentInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDel,
		},
	)

	return c
}

func (c *controller) run(ctx context.Context) {
	fmt.Println("starting controller")
	if !cache.WaitForCacheSync(ctx.Done(), c.deploymentCacheSynced) { //需要传递给连接的informer
		fmt.Print("waiting for cache to be synced\n")
	}
	go wait.Until(func() { c.runWorker(ctx) }, time.Second, ctx.Done())
	//go wait.Until(c.worker, 1*time.Second, ch)
	<-ctx.Done()
}

func (c *controller) handleAdd(obj interface{}) {
	fmt.Println("handleAdd ...")
	c.queue.Add(obj)

}

func (c *controller) handleDel(obj interface{}) {
	fmt.Println("handleDel ...")
	c.queue.Add(obj)
}

func (c *controller) runWorker(ctx context.Context) {
	for c.processItem(ctx) {

	}

}

func (c *controller) processItem(ctx context.Context) bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	defer c.queue.Forget(item)
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		fmt.Printf("getting key from cache %s\n", err.Error())
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Printf("sllitting key into namespace and name %s\n", err.Error())
		return false
	}

	err = c.syncDeployment(ns, name)
	if err != nil {
		fmt.Printf("syncing deployment %s\n", err.Error())
		return false
	}
	return true
}

func (c *controller) syncDeployment(ns, name string) error {
	if ns != "csi" {
		return nil
	}

	ctx := context.Background()

	deployment, err := c.deploymentLister.Deployments(ns).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("name '%s' in work queue no longer exists", name))
			return nil
		}
	}

	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Selector: depLabels(*deployment),
			Type:     corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 8080,
				},
			},
		},
	}
	_, err = c.clientset.CoreV1().Services(ns).Create(ctx, &svc, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("create service %s\n", err.Error())
	}
	return nil
}

func depLabels(dep appsv1.Deployment) map[string]string {
	return dep.Spec.Template.Labels

}

func initClient() *kubernetes.Clientset {
	var err error
	kubeConfig, err := ioutil.ReadFile("./admin.conf")
	restConf, err := clientcmd.RESTConfigFromKubeConfig(kubeConfig)
	clientSet, err := kubernetes.NewForConfig(restConf)
	if err != nil {
		panic(err)
	}
	return clientSet
}
