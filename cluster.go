package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	tea "github.com/charmbracelet/bubbletea"
)

type clusterMsg []*string

func getClusters() tea.Msg {
	client, err := createEcsService()
	if err != nil {
		return errorMsg(fmt.Sprintf("Error listing clusters: %s", err))
	}

	result, err := listClusters(client)
	if err != nil {
		return errorMsg(fmt.Sprintf("Error listing clusters: %s", err))
	}

	if len(result) == 0 {
		return emptyMsg("No clusters found")
	}

	return clusterMsg(result)
}

func clusterView(m *State) string {
	heading := fmt.Sprintf(`service type: %s
Please choose cluster where to connect to:`, highlightStyle(m.connectService))

	return chooseView(m, heading, m.clusters)
}

func clusterUpdate(m *State, msg tea.Msg) (*State, tea.Cmd) {
	return chooseUpdate(m, msg, len(m.clusters)-1,
		func(m *State) { m.connectService = "" },
		func(m *State) (*State, tea.Cmd) {
			m.selectedCluster = *m.clusters[m.cursor]
			m.cursor = 0
			return m, getServicesCmd(m.selectedCluster)
		},
	)
}

func listClusters(client ecsiface.ECSAPI) ([]*string, error) {
	var clusters []*string

	input := &ecs.ListClustersInput{}

	for {
		resp, err := client.ListClusters(input)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, resp.ClusterArns...)

		if resp.NextToken == nil {
			break
		}
		input.NextToken = resp.NextToken
	}

	return clusters, nil
}
