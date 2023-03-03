package main

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"time"
)

func kubernetesClientset() (*kubernetes.Clientset, error) {
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

	return kubernetes.NewForConfig(kubeConfig)
}

func connectCmd() tea.Msg {
	sleep := make(chan string)
	go func(c chan string) {
		time.Sleep(300 * time.Millisecond)
		close(c)
	}(sleep)
	cs, err := kubernetesClientset()
	if err != nil {
		return errMsg(err)
	}
	<-sleep
	return cs
}

func (m model) getNamespaces() (*v1.NamespaceList, error) {
	return m.clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
}
