package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
)

type containersMsg []*string

func getEcsConnectCmd(clusterArn string, taskArn string, containerName string, callback func(err error) tea.Msg) tea.Cmd {
	cmd := exec.Command("aws", "ecs", "execute-command", "--cluster="+clusterArn, "--task="+taskArn, "--container="+containerName, "--command=\"/bin/sh\"", "--interactive")
	return tea.ExecProcess(cmd, callback)
}

func getContainersCmd(clusterArn string, tasArn string) tea.Cmd {
	return func() tea.Msg {
		return getContainersMsg(clusterArn, tasArn)
	}
}

func getContainersMsg(clusterArn string, taskArn string) tea.Msg {
	client, err := createEcsService()
	if err != nil {
		return errorMsg(fmt.Sprintf("Error creating ECS client: %s", err))
	}
	result, err := listContainers(clusterArn, taskArn, client)
	if err != nil {
		return errorMsg(fmt.Sprintf("Error listing containers: %s", err))
	}

	if len(result) == 0 {
		return emptyMsg("No containers found")
	}
	return containersMsg(result)
}

// ecs containers
func containerView(m *State) string {
	heading := fmt.Sprintf(`service type: %s
cluster: %s
service: %s
task: %s
Please choose container where to connect to:
`,
		highlightStyle(m.connectService),
		highlightStyle(m.selectedCluster),
		highlightStyle(m.selectedService),
		highlightStyle(m.selectedTask),
	)

	return chooseView(m, heading, m.containers)
}

func containerUpdate(m *State, msg tea.Msg) (*State, tea.Cmd) {
	return chooseUpdate(m, msg, len(m.containers)-1,
		func(m *State) { m.selectedTask = "" },
		func(m *State) (*State, tea.Cmd) {
			m.selectedContainer = *m.containers[m.cursor]
			m.cursor = 0

			// final state, connect to container
			return m, getEcsConnectCmd(m.selectedCluster, m.selectedTask, m.selectedContainer, func(err error) tea.Msg {
				if err == nil {
					m.selectedContainer = ""
					return nil
				}
				return errorMessageStyle(err.Error())
			})
		},
	)
}
func listContainers(clusterArn string, taskArn string, client ecsiface.ECSAPI) ([]*string, error) {
	result, err := client.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: aws.String(clusterArn),
		Tasks:   []*string{&taskArn},
	})

	if err != nil {
		return nil, err
	}

	containers := result.Tasks[0].Containers
	var containerNames []*string
	for _, container := range containers {
		containerNames = append(containerNames, container.Name)
	}

	return containerNames, nil
}
