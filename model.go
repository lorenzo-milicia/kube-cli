package main

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lorenzo-milicia/bubbles/list"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type errMsg error

type k8s struct {
	clientset  *kubernetes.Clientset
	namespaces *v1.NamespaceList
}

type model struct {
	loading           spinner.Model
	filtering         list.Model
	choice            string
	connected         bool
	itemsFetched      bool
	itemsLoaded       bool
	kube              k8s
	namespaceSelected bool
	opSelection       list.Model
	operations        []NamespaceOperation
	selectedOperation string
	operationSelected bool
	operationOutput   string
	error             errMsg
}

func initialModel(ops []NamespaceOperation) model {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#CD931A"))
	const defaultWidth = 50

	l := list.New([]list.Item{}, itemDelegate{}, defaultWidth, 15)
	l.Title = "Select a namespace"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	//	l.SetShowFilter(false)
	l.Styles.TitleBar = lipgloss.NewStyle()
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.KeyMap.Filter.SetEnabled(false)
	l.SetFilteringMode(list.Filtering)

	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = "Search..."
	ti.Prompt = "Select a namespace: "
	ti.PromptStyle = lipgloss.NewStyle().MarginLeft(2).Bold(true)
	ti.TextStyle = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("#63d0ff"))

	l.FilterInput = ti

	var o []list.Item
	for _, op := range ops {
		fmt.Printf("op: %v\n", op)
		o = append(o, op)
	}

	oi := list.New(o, itemDelegate{}, defaultWidth, 15)

	oi.Title = "Select an operation"

	return model{
		loading:           s,
		kube:              k8s{},
		connected:         false,
		itemsFetched:      false,
		filtering:         l,
		namespaceSelected: false,
		operations:        ops,
		opSelection:       oi,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.loading.Tick, clusterConnect)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case errMsg:
		m.error = msg
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	case connectedMsg:
		m.connected = true
		m.kube.clientset = msg.clientset
		cmds = append(cmds, namespacesFetch(m))
	case namespacesMsg:
		m.itemsFetched = true
		m.kube.namespaces = msg.namespaces
		cmds = append(cmds, loadItems(&m))
		cmds = append(cmds, func() tea.Msg {
			return itemsLoadedMsg{}
		})
	case itemsLoadedMsg:
		m.itemsLoaded = true
	case KubernetesOperationMsg:
		m.operationOutput = msg.View
	}

	if !m.itemsFetched {
		s, cmd := m.loading.Update(msg)
		m.loading = s
		cmds = append(cmds, cmd)
	}
	if m.itemsLoaded && !m.namespaceSelected {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.filtering.SetWidth(msg.Width)

		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.filtering.SelectedItem().(item)
				if ok {
					m.choice = string(i)
				}
				m.namespaceSelected = true
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.filtering, cmd = m.filtering.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.namespaceSelected {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.opSelection.SetWidth(msg.Width)

		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.opSelection.SelectedItem().(item)
				if ok {
					m.selectedOperation = string(i)
				}
				if string(i) == "" {
					m.error = errMsg(errors.New("no operations present"))
					return m, tea.Quit
				}
				m.operationSelected = true
				var op NamespaceOperation
				for _, o := range m.operations {
					if o.Name == string(i) {
						op = o
					}
				}
				return m, op.Command(m.kube.clientset, m.choice)
			}
		}
		var cmd tea.Cmd
		m.opSelection, cmd = m.opSelection.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := ""
	s += clusterConnectionView(m)
	s += namespacesFetchView(m)
	s += filteringView(m)
	s += operationsView(m)
	if m.error != nil {
		s += fmt.Sprintf("Error: %v", m.error)
	}
	s += m.operationOutput
	return s
}
