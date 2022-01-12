package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	apps_v1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	storage_v1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

/*
   当前使用kubernetes版本: 1.23.1
   当前使用go版本：1.17.5
*/

const (
	TestNamespace           = "test-namespace" //测试使用的命名空间
	TestDockerConfigJsonKey = "docker-harbor"  //docker仓库密文key
)

func main() {
	clientSet := initClient()

	//createOrUpdateNamespace(clientSet)
	//listNamespace(clientSet)
	//deleteNamespace(clientSet)

	//createOrUpdateConfigMap(clientSet)
	//listConfigMap(clientSet)
	//deleteConfigMap(clientSet)

	//createOrUpdateSecret(clientSet)
	//listSecret(clientSet)
	//deleteSecret(clientSet)

	createOrUpdateDeployment(clientSet)
	//listDeployment(clientSet)
	//deleteDeployment(clientSet)

	//createOrUpdateService(clientSet)
	//listService(clientSet)
	//deleteService(clientSet)

	//createOrUpdateStorage(clientSet)
	//listStorage(clientSet)
	//deleteStorage(clientSet)

	//createOrUpdatePV(clientSet)
	//listPV(clientSet)
	//deletePV(clientSet)

	//createOrUpdatePVC(clientSet)
	//listPVC(clientSet)
	//deletePVC(clientSet)
}

/*
   创建Namespace,已存在则更新
   源码位置:K8s.io/client-go/kubernetes/typed/core/v1/namespace.go
*/
func createOrUpdateNamespace(clientSet *kubernetes.Clientset) {
	yamlFile, err := ioutil.ReadFile("./yaml/namespace.yaml")
	if err != nil {
		panic(err)
	}
	jsonBytes := yaml2Json(yamlFile)
	namespace := core_v1.Namespace{}
	err = json.Unmarshal(jsonBytes, &namespace)
	if err != nil {
		panic(err)
	}
	client := clientSet.CoreV1().Namespaces()
	if _, err = client.Get(context.TODO(), namespace.ObjectMeta.Name, meta_v1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			if _, err := client.Create(context.TODO(), &namespace, meta_v1.CreateOptions{}); err != nil {
				panic(err)
			}
			fmt.Println("Namespace创建成功")
			return
		}
		panic(err)
	}
	if _, err := client.Update(context.TODO(), &namespace, meta_v1.UpdateOptions{}); err != nil {
		panic(err)
	}
	fmt.Println("Namespace更新成功")
}

/*
   获取命名空间列表
*/
func listNamespace(clientSet *kubernetes.Clientset) {
	client := clientSet.CoreV1().Namespaces()
	namespaceList, err := client.List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	marshal, _ := json.Marshal(namespaceList)
	fmt.Println(string(marshal))
}

/*
   删除Namespace
*/
func deleteNamespace(clientSet *kubernetes.Clientset) {
	client := clientSet.CoreV1().Namespaces()
	deletePolicy := meta_v1.DeletePropagationForeground
	err := client.Delete(context.TODO(), TestNamespace, meta_v1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Namespace删除成功")
}

//*************************分割线****************************

/*
    创建密文,已存在则更新
	源码位置:K8s.io/client-go/kubernetes/typed/core/v1/secret.go
*/
func createOrUpdateSecret(clientSet *kubernetes.Clientset) {
	secret := core_v1.Secret{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      TestDockerConfigJsonKey,
			Namespace: TestNamespace,
		},
		StringData: map[string]string{
			core_v1.DockerConfigJsonKey: "{\"auths\":{\"https://registry.dockerhubar.com/\":{\"username\":\"admin\",\"password\":\"123456\"}}}",
		},
		Type: core_v1.SecretTypeDockerConfigJson,
	}
	client := clientSet.CoreV1().Secrets(TestNamespace)
	if _, err := client.Get(context.TODO(), secret.ObjectMeta.Name, meta_v1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			if _, err := client.Create(context.TODO(), &secret, meta_v1.CreateOptions{}); err != nil {
				panic(err)
			}
			fmt.Println("Secret创建成功")
			return
		}
		panic(err)
	}
	if _, err := client.Update(context.TODO(), &secret, meta_v1.UpdateOptions{}); err != nil {
		panic(err)
	}
	fmt.Println("Secret更新成功")
}

/*
	获取Secret列表,若不指定Namespace则获取所有的
*/
func listSecret(clientSet *kubernetes.Clientset) {
	client := clientSet.CoreV1().Secrets(TestNamespace)
	secretList, err := client.List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	marshal, _ := json.Marshal(secretList)
	fmt.Println(string(marshal))
}

/*
   删除Secret
*/
func deleteSecret(clientSet *kubernetes.Clientset) {
	client := clientSet.CoreV1().Secrets(TestNamespace)
	deletePolicy := meta_v1.DeletePropagationForeground
	err := client.Delete(context.TODO(), TestDockerConfigJsonKey, meta_v1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Secret删除成功")
}

//*************************分割线****************************

/*
   创建Deployment,已存在则更新
   源码位置:K8s.io/client-go/kubernetes/typed/apps/v1/deployment.go
*/
func createOrUpdateDeployment(clientSet *kubernetes.Clientset) {
	yamlFile, err := ioutil.ReadFile("./yaml/deployment.yaml")
	if err != nil {
		panic(err)
	}
	jsonBytes := yaml2Json(yamlFile)
	deployment := apps_v1.Deployment{}
	err = json.Unmarshal(jsonBytes, &deployment)
	if err != nil {
		panic(err)
	}
	deploymentClient := clientSet.AppsV1().Deployments(TestNamespace)
	if _, err = deploymentClient.Get(context.TODO(), deployment.ObjectMeta.Name, meta_v1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			if _, err := deploymentClient.Create(context.TODO(), &deployment, meta_v1.CreateOptions{}); err != nil {
				panic(err)
			}
			fmt.Println("Deployment创建成功")
			return
		}
		panic(err)
	}
	if _, err := deploymentClient.Update(context.TODO(), &deployment, meta_v1.UpdateOptions{}); err != nil {
		panic(err)
	}
	fmt.Println("Deployment更新成功")
}

/*
   获取Deployment列表,若不指定namespace则获取所有的
*/
func listDeployment(clientSet *kubernetes.Clientset) {
	client := clientSet.AppsV1().Deployments(TestNamespace)
	deploymentList, err := client.List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	marshal, _ := json.Marshal(deploymentList)
	fmt.Println(string(marshal))
}

/*
   删除Deployment
*/
func deleteDeployment(clientSet *kubernetes.Clientset) {
	client := clientSet.AppsV1().Deployments(TestNamespace)
	deletePolicy := meta_v1.DeletePropagationForeground
	err := client.Delete(context.TODO(), "svc-cloud-resourceserver", meta_v1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Deployment删除成功")
}

//*************************分割线****************************

/*
    创建Service,已存在则更新
	源码位置:K8s.io/client-go/kubernetes/typed/core/v1/service.go
*/
func createOrUpdateService(clientSet *kubernetes.Clientset) {
	yamlFile, err := ioutil.ReadFile("./yaml/service.yaml")
	if err != nil {
		panic(err)
	}
	jsonBytes := yaml2Json(yamlFile)
	service := core_v1.Service{}
	err = json.Unmarshal(jsonBytes, &service)
	if err != nil {
		panic(err)
	}
	client := clientSet.CoreV1().Services(TestNamespace)
	existService, err := client.Get(context.TODO(), service.ObjectMeta.Name, meta_v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			if _, err := client.Create(context.TODO(), &service, meta_v1.CreateOptions{}); err != nil {
				panic(err)
			}
			fmt.Println("Service创建成功")
			return
		}
		panic(err)
	}
	if _, err := client.Update(context.TODO(), existService, meta_v1.UpdateOptions{}); err != nil {
		panic(err)
	}
	fmt.Println("service更新成功")
}

/*
   获取Service列表,若不指定namespace则获取所有的
*/
func listService(clientSet *kubernetes.Clientset) {
	client := clientSet.CoreV1().Services(TestNamespace)
	serviceList, err := client.List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	marshal, _ := json.Marshal(serviceList)
	fmt.Println(string(marshal))
}

/*
    删除Service
    三种删除策略：
	Orphan：    只删除当前对象，不删除其所管理的资源对象
	Background：删除之后，所管理的资源对象由GC删除
	Foreground：删除之前所管理的资源对象必须先删除
*/
func deleteService(clientSet *kubernetes.Clientset) {
	client := clientSet.CoreV1().Services(TestNamespace)
	deletePolicy := meta_v1.DeletePropagationForeground
	err := client.Delete(context.TODO(), "svc-cloud-resourceserver", meta_v1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Service删除成功")
}

//*************************分割线****************************

/*
    创建Storage,已存在则更新
	源码位置:K8s.io/client-go/kubernetes/typed/storage/v1/storageclass.go
*/
func createOrUpdateStorage(clientSet *kubernetes.Clientset) {
	yamlFile, err := ioutil.ReadFile("./yaml/storageClass.yaml")
	if err != nil {
		panic(err)
	}
	jsonBytes := yaml2Json(yamlFile)
	storageClass := storage_v1.StorageClass{}
	err = json.Unmarshal(jsonBytes, &storageClass)
	if err != nil {
		panic(err)
	}
	client := clientSet.StorageV1().StorageClasses()
	if _, err := client.Get(context.TODO(), storageClass.ObjectMeta.Name, meta_v1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			if _, err := client.Create(context.TODO(), &storageClass, meta_v1.CreateOptions{}); err != nil {
				panic(err)
			}
			fmt.Println("StorageClass创建成功")
			return
		}
		panic(err)
	}
	if _, err := client.Update(context.TODO(), &storageClass, meta_v1.UpdateOptions{}); err != nil {
		panic(err)
	}
	fmt.Println("StorageClass更新成功")
}

/*
  获取storageClass列表
*/
func listStorage(clientSet *kubernetes.Clientset) {
	client := clientSet.StorageV1().StorageClasses()
	storageClassList, err := client.List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	marshal, _ := json.Marshal(storageClassList)
	fmt.Println(string(marshal))
}

/*
	删除storageClass
*/
func deleteStorage(clientSet *kubernetes.Clientset) {
	client := clientSet.StorageV1().StorageClasses()
	deletePolicy := meta_v1.DeletePropagationForeground
	err := client.Delete(context.TODO(), "test-storage-class", meta_v1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("StorageClass删除成功")
}

//*************************分割线****************************

/*
    创建ConfigMap,已存在则更新
	源码位置:K8s.io/client-go/kubernetes/typed/core/v1/configmap.go
*/
func createOrUpdateConfigMap(clientSet *kubernetes.Clientset) {
	yamlFile, err := ioutil.ReadFile("./yaml/configMap.yaml")
	if err != nil {
		panic(err)
	}
	jsonBytes := yaml2Json(yamlFile)
	configMap := core_v1.ConfigMap{}
	err = json.Unmarshal(jsonBytes, &configMap)
	if err != nil {
		panic(err)
	}
	client := clientSet.CoreV1().ConfigMaps(TestNamespace)
	if _, err := client.Get(context.TODO(), configMap.ObjectMeta.Name, meta_v1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			if _, err := client.Create(context.TODO(), &configMap, meta_v1.CreateOptions{}); err != nil {
				panic(err)
			}
			fmt.Println("ConfigMap创建成功")
			return
		}
		panic(err)
	}
	if _, err := client.Update(context.TODO(), &configMap, meta_v1.UpdateOptions{}); err != nil {
		panic(err)
	}
	fmt.Println("ConfigMap更新成功")
}

/*
   获取ConfigMap列表,若不指定namespace则获取所有的
*/
func listConfigMap(clientSet *kubernetes.Clientset) {
	client := clientSet.CoreV1().ConfigMaps(TestNamespace)
	configMapList, err := client.List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	marshal, _ := json.Marshal(configMapList)
	fmt.Println(string(marshal))
}

/*
   删除ConfigMap
*/
func deleteConfigMap(clientSet *kubernetes.Clientset) {
	client := clientSet.CoreV1().ConfigMaps(TestNamespace)
	deletePolicy := meta_v1.DeletePropagationForeground
	err := client.Delete(context.TODO(), "test-configmap-nginx", meta_v1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("ConfigMap删除成功")
}

//*************************分割线****************************

/*
    创建PersistentVolume,已存在则更新
	源码位置:K8s.io/client-go/kubernetes/typed/core/v1/persistentvolume.go
*/
func createOrUpdatePV(clientSet *kubernetes.Clientset) {
	yamlFile, err := ioutil.ReadFile("./yaml/persistentVolume.yaml")
	if err != nil {
		panic(err)
	}
	jsonBytes := yaml2Json(yamlFile)
	pv := core_v1.PersistentVolume{}
	err = json.Unmarshal(jsonBytes, &pv)
	if err != nil {
		panic(err)
	}
	client := clientSet.CoreV1().PersistentVolumes()
	if _, err := client.Get(context.TODO(), pv.ObjectMeta.Name, meta_v1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			if _, err := client.Create(context.TODO(), &pv, meta_v1.CreateOptions{}); err != nil {
				panic(err)
			}
			fmt.Println("PV创建成功")
			return
		}
		panic(err)
	}
	if _, err := client.Update(context.TODO(), &pv, meta_v1.UpdateOptions{}); err != nil {
		panic(err)
	}
	fmt.Println("PV更新成功")
}

/*
  获取PersistentVolume列表
*/
func listPV(clientSet *kubernetes.Clientset) {
	client := clientSet.CoreV1().PersistentVolumes()
	pvList, err := client.List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	marshal, _ := json.Marshal(pvList)
	fmt.Println(string(marshal))
}

/*
	删除PersistentVolume
*/
func deletePV(clientSet *kubernetes.Clientset) {
	client := clientSet.CoreV1().PersistentVolumes()
	deletePolicy := meta_v1.DeletePropagationForeground
	err := client.Delete(context.TODO(), "test-pv", meta_v1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("PV删除成功")
}

//*************************分割线****************************

/*
    创建PersistentVolumeClaim,已存在则更新
	源码位置:K8s.io/client-go/kubernetes/typed/core/v1/persistentvolumeclaim.go
*/
func createOrUpdatePVC(clientSet *kubernetes.Clientset) {
	yamlFile, err := ioutil.ReadFile("./yaml/persistentVolumeClaim.yaml")
	if err != nil {
		panic(err)
	}
	jsonBytes := yaml2Json(yamlFile)
	pvc := core_v1.PersistentVolumeClaim{}
	err = json.Unmarshal(jsonBytes, &pvc)
	if err != nil {
		panic(err)
	}
	client := clientSet.CoreV1().PersistentVolumeClaims(TestNamespace)
	if _, err := client.Get(context.TODO(), pvc.ObjectMeta.Name, meta_v1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			if _, err := client.Create(context.TODO(), &pvc, meta_v1.CreateOptions{}); err != nil {
				panic(err)
			}
			fmt.Println("PVC创建成功")
			return
		}
		panic(err)
	}
	if _, err := client.Update(context.TODO(), &pvc, meta_v1.UpdateOptions{}); err != nil {
		panic(err)
	}
	fmt.Println("PVC更新成功")
}

/*
  获取PersistentVolumeClaim列表
*/
func listPVC(clientSet *kubernetes.Clientset) {
	client := clientSet.CoreV1().PersistentVolumeClaims(TestNamespace)
	pvList, err := client.List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	marshal, _ := json.Marshal(pvList)
	fmt.Println(string(marshal))
}

/*
	删除PersistentVolumeClaim
*/
func deletePVC(clientSet *kubernetes.Clientset) {
	client := clientSet.CoreV1().PersistentVolumeClaims(TestNamespace)
	deletePolicy := meta_v1.DeletePropagationForeground
	err := client.Delete(context.TODO(), "test-pvc", meta_v1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("PVC删除成功")
}

//*************************分割线****************************

/*
   读取配置文件并且初始化客户端
   kubeconfig 默认在主节点 /etc/kubernetes/admin.conf
   一般在 $HOME/.kube/config 也会复制一份用于身份认证
*/
func initClient() *kubernetes.Clientset {
	var err error
	kubeConfig, err := ioutil.ReadFile("./config")
	restConf, err := clientcmd.RESTConfigFromKubeConfig(kubeConfig)
	clientSet, err := kubernetes.NewForConfig(restConf)
	if err != nil {
		panic(err)
	}
	return clientSet
}

/*
   yaml转json
*/
func yaml2Json(yamlBytes []byte) (jsonBytes []byte) {
	toJSON, err := yaml.ToJSON(yamlBytes)
	if err != nil {
		panic(err)
	}
	return toJSON
}
