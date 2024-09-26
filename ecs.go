package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func ecsUpdate(m *State, msg tea.Msg) (tea.Model, tea.Cmd) {
	// persist state
	switch msg := msg.(type) {
	case clusterMsg:
		m.clusters = msg
		if len(m.clusters) == 1 {
			m.selectedCluster = *m.clusters[0]
			return m, getServicesCmd(m.selectedCluster)
		}
	case servicesMsg:
		m.services = msg
		if len(m.services) == 1 {
			m.selectedService = *m.services[0]
			return m, getTasksCmd(m.selectedCluster, m.selectedService)
		}
	case tasksMsg:
		m.tasks = msg
		if len(m.tasks) == 1 {
			m.selectedTask = *m.tasks[0]
			return m, getContainersCmd(m.selectedCluster, m.selectedTask)
		}
	case containersMsg:
		m.containers = msg
		if len(m.containers) == 1 {
			m.selectedContainer = *m.containers[0]
			return m, getEcsConnectCmd(m.selectedCluster, m.selectedTask, m.selectedContainer, func(err error) tea.Msg {
				if err == nil {
					m.selectedContainer = ""
					return nil
				}
				return errorMessageStyle(err.Error())
			})
		}
	case errorMsg:
		m.error = string(msg)
		return m, nil
	case emptyMsg:
		m.error = string(msg)
		return m, nil
	}

	// handle selection
	if m.selectedCluster == "" {
		return clusterUpdate(m, msg)
	}
	if m.selectedCluster != "" && m.selectedService == "" {
		return serviceUpdate(m, msg)
	}
	if m.selectedService != "" && m.selectedTask == "" {
		return taskUpdate(m, msg)
	}
	if m.selectedTask != "" && m.selectedContainer == "" {
		return containerUpdate(m, msg)
	}

	return m, nil
}

func ecsView(m *State) string {
	if m.selectedCluster == "" && len(m.clusters) > 0 {
		return clusterView(m)
	}
	if m.selectedService == "" && m.selectedCluster != "" && len(m.services) > 0 {
		return serviceView(m)
	}
	if m.selectedTask == "" && m.selectedService != "" && len(m.tasks) > 0 {
		return taskView(m)
	}
	if m.selectedContainer == "" && m.selectedTask != "" && len(m.containers) > 0 {
		return containerView(m)
	}

	if m.error != "" {
		return headlineStyle(header.String()) + "\n" + errorMessageStyle(m.error) + footerStyle("backspace: back, enter/space: select, ctrl+c/q: quit, up/k: up, down/j: down")
	}

	return "loading..."
}
