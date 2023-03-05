package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type connectedMsg struct {
	clientset *kubernetes.Clientset
}
type namespacesMsg struct {
	namespaces *v1.NamespaceList
}

func clusterConnect() tea.Msg {
	cs, err := getK8sConnection()
	if err != nil {
		return errMsg(err)
	}
	return connectedMsg{cs}
}

func namespacesFetch(m model) tea.Cmd {
	return func() tea.Msg {
		n, err := getNamespaces(m.kube.clientset)
		if err != nil {
			return errMsg(err)
		}
		return namespacesMsg{n}
	}
}

var (
	kubernetesStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#326CE5"))
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#b0ffa6"))
)

func styleString(s string, style lipgloss.Style) lipgloss.Style {
	return style.SetString(s)
}

func clusterConnectionView(m model) string {
	kString := styleString("Kubernetes cluster", kubernetesStyle)
	if !m.connected {
		return fmt.Sprintf("%s Connecting to the %s...\n", m.loading.View(), kString)
	} else {
		return fmt.Sprintf("%s Connected to the %s\n", styleString("✔", successStyle), kString)
	}
}

func namespacesFetchView(m model) string {
	if !m.connected {
		return ""
	}
	nsString := styleString("namespaces", kubernetesStyle)
	if !m.itemsFetched {
		return fmt.Sprintf("%s Fetching %s...\n", m.loading.View(), nsString)
	} else {
		return fmt.Sprintf("%s Fetched %s\n", styleString("✔", successStyle), nsString)
	}
}
