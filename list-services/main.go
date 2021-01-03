package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	KubeConfig = flag.String("kubeconfig", "", "Path to the kubeconfig file to use for CLI requests.")
)

func main() {
	flag.Parse()

	// create kubernetes client
	client, err := newClient(*KubeConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Get all services in all namespace
	services, err := client.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, service := range services.Items {
		fmt.Printf("Service name= %s namespace=%s type=%s \n", service.Name, service.Namespace, service.Spec.Type)
		if service.Spec.Type == "ClusterIP" {
			for _, port := range service.Spec.Ports {
				if port.Name != "" {
					fmt.Printf("ClientIP Port: name=%s port=%d\n", port.Name, port.Port)
				}
			}
		}
	}
}

func newClient(kubeConfigPath string) (kubernetes.Interface, error) {
	if kubeConfigPath == "" {
		// use default path(.kube/config) if kubeconfig path is not set
		kubeConfigPath = clientcmd.RecommendedHomeFile
	}
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(kubeConfig)
}
