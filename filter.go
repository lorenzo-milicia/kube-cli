package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lorenzo-milicia/bubbles/list"
	"io"
	"strings"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(0).MarginBottom(0)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	matchStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#63d0ff"))
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(0)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type itemsLoadedMsg struct{}

type item string

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := string(i)

	filterString := m.FilterInput.Value()
	str = strings.ReplaceAll(str, filterString, fmt.Sprintf("%s", matchStyle.Render(filterString)))

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render(matchStyle.Render("> ") + str)
		}
	}

	fmt.Fprint(w, fn(str))
}

func loadItems(m *model) tea.Cmd {
	var items []list.Item
	for _, ns := range m.kube.namespaces.Items {
		items = append(items, item(ns.Name))
	}
	return m.filtering.SetItems(items)
}

func filteringView(m model) string {
	if !m.itemsLoaded || m.done {
		return ""
	}
	return m.filtering.View()
}
