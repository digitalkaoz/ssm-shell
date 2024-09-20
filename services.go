package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	tea "github.com/charmbracelet/bubbletea"
)

type servicesMsg []*string

func getServicesCmd(clusterArn string) tea.Cmd {
	return func() tea.Msg {
		return getServicesMsg(clusterArn)
	}
}
func getServicesMsg(clusterArn string) tea.Msg {
	client, err := createEcsService()
	if err != nil {
		return errorMsg(fmt.Sprintf("Error creating ECS client: %s", err))
	}
	result, err := listServices(clusterArn, client)
	if err != nil {
		return errorMsg(fmt.Sprintf("Error listing services: %s", err))
	}

	if len(result) == 0 {
		return emptyMsg("No services found")
	}

	return servicesMsg(result)
}

func serviceView(m *State) string {
	heading := fmt.Sprintf(`service type: %s
cluster: %s
Please choose service where to connect to:
`,
		highlightStyle(m.connectService),
		highlightStyle(m.selectedCluster),
	)

	return chooseView(m, heading, m.services)
}

func serviceUpdate(m *State, msg tea.Msg) (*State, tea.Cmd) {
	return chooseUpdate(m, msg, len(m.services)-1,
		func(m *State) { m.selectedCluster = "" },
		func(m *State) (*State, tea.Cmd) {
			m.selectedService = *m.services[m.cursor]
			m.cursor = 0

			return m, getTasksCmd(m.selectedCluster, m.selectedService)
		},
	)
}

func listServices(clusterArn string, client ecsiface.ECSAPI) ([]*string, error) {
	var services []*string
	input := &ecs.ListServicesInput{
		Cluster: aws.String(clusterArn),
	}
	for {
		result, err := client.ListServices(input)
		if err != nil {
			return nil, err
		}
		services = append(services, result.ServiceArns...)
		if result.NextToken == nil {
			break
		}
		input.NextToken = result.NextToken
	}
	return services, nil
}
