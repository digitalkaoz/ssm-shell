package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
	"regexp"
)

type instancesMsg []*string

func getEc2ConnectCmd(instanceId string, callback func(err error) tea.Msg) tea.Cmd {
	re := regexp.MustCompile(`\((.+)\)`)
	matches := re.FindStringSubmatch(instanceId)
	cmd := exec.Command("aws", "ssm", "start-session", "--target="+matches[1])
	return tea.ExecProcess(cmd, callback)
}

func getInstances() tea.Msg {
	result, err := listInstances()
	if err != nil {
		return errorMsg(fmt.Sprintf("Error listing instances: %s", err))
	}

	if len(result) == 0 {
		return emptyMsg("No instances found")
	}

	return instancesMsg(result)
}

func instanceView(m *State) string {
	heading := fmt.Sprintf(`service type: %s
Please choose instance where to connect to:
`,
		highlightStyle(m.connectService),
	)

	return chooseView(m, heading, m.instances)
}

func instanceUpdate(m *State, msg tea.Msg) (*State, tea.Cmd) {
	return chooseUpdate(m, msg, len(m.instances)-1,
		func(m *State) { m.connectService = "" },
		func(m *State) (*State, tea.Cmd) {
			m.selectedInstance = *m.instances[m.cursor]
			m.cursor = 0

			// final ec2 state
			return m, getEc2ConnectCmd(m.selectedInstance, func(err error) tea.Msg {
				m.selectedInstance = ""
				if err != nil {
					return errorMsg(fmt.Sprintf("Error connecting to instance: %s", err))
				}
				return nil
			})
		},
	)
}

func listInstances() ([]*string, error) {
	sess, err := createAwsSession()
	if err != nil {
		return nil, err
	}
	svc := ec2.New(sess)

	result, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
		},
		MaxResults: aws.Int64(100),
	})
	//TODO paging

	if err != nil {
		return nil, err
	}
	instances := make([]*string, len(result.Reservations))
	for i, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			prettyName := fmt.Sprintf("%s (%s)", *instance.Tags[0].Value, *instance.InstanceId)
			instances[i] = &prettyName
		}
	}

	return instances, nil
}
