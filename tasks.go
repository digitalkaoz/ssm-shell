package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	tea "github.com/charmbracelet/bubbletea"
)

type tasksMsg []*string

func getTasksCmd(clusterArn string, serviceArn string) tea.Cmd {
	return func() tea.Msg {
		return getTasksMsg(clusterArn, serviceArn)
	}
}

func getTasksMsg(clusterArn string, serviceArn string) tea.Msg {
	client, err := createEcsService()
	if err != nil {
		return errorMsg(fmt.Sprintf("Error creating ECS client: %s", err))
	}
	result, err := listTasks(clusterArn, serviceArn, client)
	if err != nil {
		return errorMsg(fmt.Sprintf("Error listing tasks: %s", err))
	}
	if len(result) == 0 {
		return emptyMsg("No tasks found")
	}

	return tasksMsg(result)
}

func taskView(m *State) string {
	heading := fmt.Sprintf(`service type: %s
cluster: %s
service: %s
Please choose task where to connect to:
`,
		highlightStyle(m.connectService),
		highlightStyle(m.selectedCluster),
		highlightStyle(m.selectedService),
	)

	return chooseView(m, heading, m.tasks)
}

func taskUpdate(m *State, msg tea.Msg) (*State, tea.Cmd) {
	return chooseUpdate(m, msg, len(m.tasks)-1,
		func(m *State) { m.selectedService = "" },
		func(m *State) (*State, tea.Cmd) {
			m.selectedTask = *m.tasks[m.cursor]
			m.cursor = 0

			return m, getContainersCmd(m.selectedCluster, m.selectedTask)
		},
	)
}

func listTasks(clusterArn string, serviceArn string, client ecsiface.ECSAPI) ([]*string, error) {
	input := &ecs.ListTasksInput{
		Cluster:     aws.String(clusterArn),
		ServiceName: aws.String(serviceArn),
		MaxResults:  aws.Int64(100),
	}

	var tasks []*string
	for {
		result, err := client.ListTasks(input)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, result.TaskArns...)
		if result.NextToken == nil {
			break
		}
		input.NextToken = result.NextToken
	}

	return tasks, nil
}
