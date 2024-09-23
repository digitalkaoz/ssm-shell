package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
	"regexp"
	"strings"
)

type instancesMsg []*string

// bubbletea command generator
func getEc2ConnectCmd(instanceId string, callback func(err error) tea.Msg) tea.Cmd {
	re := regexp.MustCompile(`\((.+)\)`)
	matches := re.FindStringSubmatch(instanceId)
	cmd := exec.Command("aws", "ssm", "start-session", "--target="+matches[1])
	return tea.ExecProcess(cmd, callback)
}

func getInstances() tea.Msg {
	client, err := createEc2Client()
	if err != nil {
		return errorMsg(fmt.Sprintf("Error listing instances: %s", err))
	}
	result, err := listInstances(client)
	if err != nil {
		return errorMsg(fmt.Sprintf("Error listing instances: %s", err))
	}

	if len(result) == 0 {
		return emptyMsg("No instances found")
	}

	return instancesMsg(result)
}

// bubbletea view handler
func instanceView(m *State) string {
	heading := fmt.Sprintf(`service type: %s
Please choose instance where to connect to:
`,
		highlightStyle(m.connectService),
	)

	return chooseView(m, heading, m.instances)
}

// bubbletea update handler
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

func listInstances(client ec2iface.EC2API) ([]*string, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
		},
	}

	var reservations []*ec2.Reservation
	for {
		output, err := client.DescribeInstances(input)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, output.Reservations...)

		if output.NextToken == nil {
			break
		}
		input.NextToken = output.NextToken
	}

	return getInstanceNames(reservations), nil
}

func getInstanceNames(reservations []*ec2.Reservation) []*string {
	instances := make([]*string, len(reservations))

	for i, reservation := range reservations {
		for _, instance := range reservation.Instances {
			prettyName := strings.TrimLeft(fmt.Sprintf("%s (%s)", getNameFromInstance(instance), *instance.InstanceId), " ")
			instances[i] = &prettyName
		}
	}

	return instances
}

func getNameFromInstance(instance *ec2.Instance) string {
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			if tag.Value != nil {
				return *tag.Value
			}
		}
	}
	return ""
}
