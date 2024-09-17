package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

type State struct {
	connectService       string
	availableConnections []*string

	clusters        []*string
	selectedCluster string

	services        []*string
	selectedService string

	tasks        []*string
	selectedTask string

	containers        []*string
	selectedContainer string

	instances        []*string
	selectedInstance string

	cursor int
	error  string
}

type errorMsg string
type emptyMsg string

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FDFDF5"}).Render

	errorMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
			Render

	headlineStyle = lipgloss.NewStyle().
			Padding(2, 2).
			Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
			Render

	footerStyle = lipgloss.NewStyle().
			Padding(2, 2).
			Foreground(lipgloss.AdaptiveColor{Light: "#A3A3A3", Dark: "#8A8A8A"}).
			Render

	header = figure.NewFigure("SSM Shell", "", true)
)

func InitialState() *State {
	ptr := func(s string) *string {
		return &s
	}

	return &State{
		availableConnections: []*string{ptr("ECS"), ptr("EC2")},
	}
}

func (m *State) Init() tea.Cmd {
	return nil
}

func (m *State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.connectService == "" {
		return connectTypeUpdate(*m, msg)
	}
	if m.connectService == "ECS" {
		return ecsUpdate(m, msg)
	}

	if m.connectService == "EC2" {
		return ec2Update(m, msg)
	}

	return m, nil
}

func (m *State) View() string {
	s := ""

	if m.connectService == "" {
		s = connectTypeView(*m)
	}
	if m.connectService == "ECS" {
		return ecsView(m)
	}
	if m.connectService == "EC2" {
		return ec2View(m)
	}

	return s
}

func chooseView(m *State, headline string, choices []*string) string {
	s := headlineStyle(header.String()) + "\n" + headline + " \n\n"

	for i, choice := range choices {
		if m.cursor == i {
			s += fmt.Sprintf("* %s\n", highlightStyle(*choice))
		} else {
			s += fmt.Sprintf("  %s\n", titleStyle(*choice))
		}
	}

	if m.error != "" {
		s += errorMessageStyle(m.error) + "\n\n"
	}

	s += footerStyle("backspace: back, enter/space: select, ctrl+c/q: quit, up/k: up, down/j: down")

	return s
}

func chooseUpdate(m *State, msg tea.Msg, max int, back func(*State), setter func(*State) (*State, tea.Cmd)) (*State, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			m.error = ""
			back(m)
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < max {
				m.cursor++
			}
		case "enter", " ":
			return setter(m)
		}
	}
	return m, nil
}
