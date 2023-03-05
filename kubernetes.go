package main

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"time"
)

func getK8sConnection() (*kubernetes.Clientset, error) {
	// To add a minimim spinner time
	sleep := make(chan string)
	go func(c chan string) {
		time.Sleep(500 * time.Millisecond)
		close(c)
	}(sleep)

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		return nil, err
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("Error getting kubernetes config: %v\n", err)
		return nil, err
	}
	<-sleep
	k8s, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {

		return nil, err
	}
	return k8s, nil
}

func getNamespaces(k8s *kubernetes.Clientset) (*v1.NamespaceList, error) {
	sleep := make(chan string)
	go func(c chan string) {
		time.Sleep(500 * time.Millisecond)
		close(c)
	}(sleep)
	nl, err := k8s.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	<-sleep
	return nl, nil
}
