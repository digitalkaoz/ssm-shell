package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	tea "github.com/charmbracelet/bubbletea"
)

type clusterMsg []*string

func getClusters() tea.Msg {
	result, err := listClusters()
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

func listClusters() ([]*string, error) {
	sess, err := createAwsSession()
	if err != nil {
		return nil, err
	}
	svc := ecs.New(sess)

	result, err := svc.ListClusters(&ecs.ListClustersInput{
		MaxResults: aws.Int64(100),
	})
	//TODO paging

	if err != nil {
		return nil, err
	}
	return result.ClusterArns, nil
}
