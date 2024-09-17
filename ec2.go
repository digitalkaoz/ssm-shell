package main

import tea "github.com/charmbracelet/bubbletea"

func ec2Update(m *State, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case instancesMsg:
		m.instances = msg
		if len(m.instances) == 1 {
			m.selectedInstance = *m.instances[0]
			return m, getEc2ConnectCmd(m.selectedInstance, func(err error) tea.Msg {
				if err == nil {
					m.selectedInstance = ""
					return nil
				}
				return errorMessageStyle(err.Error())
			})
		}

	case errorMsg, emptyMsg:
		m.error = msg.(string)
		return m, nil
	}

	if m.selectedInstance == "" {
		return instanceUpdate(m, msg)
	}

	return m, nil
}

func ec2View(m *State) string {
	if m.selectedInstance == "" && len(m.instances) > 0 {
		return instanceView(m)
	}

	if m.error != "" {
		return headlineStyle(header.String()) + "\n" + errorMessageStyle(m.error) + footerStyle("backspace: back, enter/space: select, ctrl+c/q: quit, up/k: up, down/j: down")
	}

	return "loading..."
}
