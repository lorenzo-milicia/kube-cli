package operations

import (
	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/client-go/kubernetes"
)

type NamespaceOperation struct {
	Name    string
	Command KubernetesNamespaceCommand
}

func (o NamespaceOperation) FilterValue() string { return o.Name }

type KubernetesOperationMsg struct{
	View string
}

type KubernetesNamespaceCommand func(clientset *kubernetes.Clientset, namespace string) tea.Cmd
