package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	KubeConfig = flag.String("kubeconfig", "", "kubeconfig file")
)

func main() {
	// create kubernetes client
	client, err := newClient(*KubeConfig)
	if err != nil {
		log.Fatal(err)
	}

	namespace := "default"
	appname := "httpbin"

	svc, err := getService(appname, namespace, client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	pods, err := getPodsForService(svc, namespace, client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	for _, pod := range pods.Items {
		fmt.Printf("Pod Name %s \n", pod.Name)
		fmt.Printf("Pod Namespace %s \n", pod.Namespace)
		//Dump labels
		labels := pod.Labels
		for k, v := range labels {
			fmt.Printf("Pod %s Namespace %s Label %s:%s\n", pod.Name, pod.Namespace, k, v)
		}
	}

}

func newClient(kubeConfigPath string) (kubernetes.Interface, error) {
	if kubeConfigPath == "" {
		kubeConfigPath = os.Getenv("KUBECONFIG")
	}
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

// ref: https://stackoverflow.com/questions/41545123/how-to-get-pods-under-the-service-with-client-go-the-client-library-of-kubernete
func getService(deployment string, namespace string, client kubernetes.Interface) (*v1.Service, error) {
	listOptions := metav1.ListOptions{}
	services, err := client.CoreV1().Services(namespace).List(context.TODO(), listOptions)
	if err != nil {
		log.Fatal(err)
	}
	for _, service := range services.Items {
		if strings.Contains(service.Name, deployment) {
			fmt.Fprintf(os.Stdout, "service name: %v\n", service.Name)
			return &service, nil
		}
	}
	return nil, errors.New("cannot find service for deployment")
}

func getPodsForService(svc *v1.Service, namespace string, client kubernetes.Interface) (*v1.PodList, error) {
	set := labels.Set(svc.Spec.Selector)
	listOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
	return pods, err
}
