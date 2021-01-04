package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	// Get specific pod info
	namespace := "default"
	podName := "busybox"
	_, err = client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", podName, namespace)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", podName, namespace)
	}

	// List all pods
	namespace = "" // empty - all namespace
	list, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing pods: %v", err)
		os.Exit(1)
	}

	for _, pod := range list.Items {
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
