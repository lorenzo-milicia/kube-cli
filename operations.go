package main

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceOperation struct {
	Name    string
	Command KubernetesNamespaceCommand
}

func (o NamespaceOperation) FilterValue() string { return o.Name }

type KubernetesOperationMsg struct {
	View string
}

type KubernetesNamespaceCommand func(clientset *kubernetes.Clientset, namespace string) tea.Cmd

var PodsOperation = NamespaceOperation{
	Name: "Get list of pods",
	Command: func(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
		return func() tea.Msg {
			p, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return errMsg(err)
			}
			s := ""
			for _, pod := range p.Items {
				s += pod.Name + "\n"
			}
			return KubernetesOperationMsg{View: s}
		}
	},
}
