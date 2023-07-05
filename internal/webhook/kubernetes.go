package webhook

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

type Status struct {
	Deployment  string `json:"Name"`
	RestartedAt string `json:"RestartedAt"`
}

type Statuses []Status

func appInit() *rest.Config {

	config, err := rest.InClusterConfig()
	fmt.Printf("incluster error %v\n", err)
	if err != nil {
		// fallback to kubeconfig
		home := homedir.HomeDir()
		kubeconfig := filepath.Join(home, ".kube", "config")
		if envvar := os.Getenv("KUBECONFIG"); len(envvar) > 0 {
			kubeconfig = envvar
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Printf("The kubeconfig cannot be loaded: %v\n", err)
			os.Exit(1)
		}

	}
	return config
}

func UpdateDeployment(name string) map[string]string {

	clientset, err := kubernetes.NewForConfig(appInit())
	if err != nil {
		panic(err)
	}
	namespace := os.Getenv("WORKER_NAMESPACE")
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	result, getErr := deploymentsClient.Get(context.TODO(), name, metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("failed to get latest version of deployment: %v", getErr))
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		if len(result.Spec.Template.GetAnnotations()) != 0 {
			annotate := result.Spec.Template.GetAnnotations()
			annotate["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
			result.Spec.Template.Annotations = annotate
			for i, d := range annotate {
				log.Printf("annotations at %s: %s\n", i, d)
			}

		} else {
			log.Println("No annoattions")
			annotate := make(map[string]string)
			annotate["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
			result.Spec.Template.Annotations = annotate
			for i, d := range annotate {
				log.Printf("annotations from empty at %s: %s\n", i, d)
			}
		}

		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("update failed: %v", retryErr))
	}
	return result.Spec.Template.GetAnnotations()

}

// func DeploymentGet(namespace string) []string {
// 	var deploymentList []string
// 	clientset, err := kubernetes.NewForConfig(appInit())
// 	if err != nil {
// 		panic(err)
// 	}
// 	deploymentsClient := clientset.AppsV1().Deployments(namespace)
// 	result, getErr := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
// 	if getErr != nil {
// 		panic(fmt.Errorf("failed to get list of deployments: %v", getErr))
// 	}
// 	for _, d := range result.Items {
// 		deploymentList = append(deploymentList, d.Name)
// 		// fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
// 	}
// 	return deploymentList
// }
