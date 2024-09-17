package main

import tea "github.com/charmbracelet/bubbletea"

func connectTypeView(m State) string {
	return chooseView(&m, "Please choose service where to connect to:", m.availableConnections)
}

func connectTypeUpdate(m State, msg tea.Msg) (*State, tea.Cmd) {
	return chooseUpdate(&m, msg, len(m.availableConnections)-1,
		func(m *State) { m.connectService = "" },
		func(m *State) (*State, tea.Cmd) {
			m.connectService = *m.availableConnections[m.cursor]
			if m.connectService == "ECS" {
				return m, getClusters
			}
			return m, getInstances
		},
	)
}
