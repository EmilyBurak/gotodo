package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63"))

type Task struct {
	ID              string
	Name            string
	Status          string
	Deleted         string
	Pomodoros       string
	PomodorosNeeded string
}

func (t Task) FilterValue() string { // implements list.Item interface with FilterValue() method
	return t.Name
}

func (t Task) Title() string {
	return t.Name
}

func (t Task) Description() string {
	return t.Status
}

type Model struct {
	list     list.Model
	choice   string
	quitting bool
}

// initList initializes the list with the tasks
func (m Model) InitList(width, height int, tasks []Task) *Model { // pass in list of tasks?
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list.Title = "Tasks"
	m.list.Styles = list.Styles{
		Title: style,
	}
	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		items[i] = task
	}
	m.list.SetItems(items)
	return &m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(Task)
			if ok {
				m.choice = fmt.Sprintf("Task: %s", i.Name)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {

	if m.choice != "" {
		return style.Render(fmt.Sprintf("Selected: %s", m.choice))
	}
	if m.quitting {
		return style.Render("Quitting...")
	}
	return "\n" + m.list.View()
}
