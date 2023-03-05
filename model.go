package main

import (
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
	loading      spinner.Model
	filtering    list.Model
	choice       string
	connected    bool
	itemsFetched bool
	itemsLoaded  bool
	kube         k8s
	done         bool
}

func initialModel() model {
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

	return model{loading: s, kube: k8s{}, connected: false, itemsFetched: false, filtering: l, done: false}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.loading.Tick, clusterConnect)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case errMsg:
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
	}

	if !m.itemsFetched {
		s, cmd := m.loading.Update(msg)
		m.loading = s
		cmds = append(cmds, cmd)
	}
	if m.itemsLoaded {
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
				m.done = true
				return m, tea.Quit
			}
		}
		var cmd tea.Cmd
		m.filtering, cmd = m.filtering.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := ""
	s += clusterConnectionView(m)
	s += namespacesFetchView(m)
	s += filteringView(m)
	if m.done {
		s += fmt.Sprintf("You chose... %s\n", m.choice)
	}
	return s
}
