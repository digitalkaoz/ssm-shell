package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	tea "github.com/charmbracelet/bubbletea"
)

type servicesMsg []*string

func getServicesCmd(clusterArn string) tea.Cmd {
	return func() tea.Msg {
		return getServicesMsg(clusterArn)
	}
}
func getServicesMsg(clusterArn string) tea.Msg {
	result, err := listServices(clusterArn)
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

func listServices(clusterArn string) ([]*string, error) {
	sess, err := createAwsSession()
	if err != nil {
		return nil, err
	}
	svc := ecs.New(sess)

	result, err := svc.ListServices(&ecs.ListServicesInput{
		Cluster:    aws.String(clusterArn),
		MaxResults: aws.Int64(100),
	})
	//TODO paging

	if err != nil {
		return nil, err
	}
	return result.ServiceArns, nil
}
