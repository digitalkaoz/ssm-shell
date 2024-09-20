package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"reflect"
	"testing"
)

type mockECSClientContainers struct {
	ecsiface.ECSAPI
	Response ecs.DescribeTasksOutput
	Error    error
}

func (m *mockECSClientContainers) DescribeTasks(in *ecs.DescribeTasksInput) (*ecs.DescribeTasksOutput, error) {
	if m.Response.Tasks == nil {
		return nil, errors.New("error")
	}

	return &m.Response, m.Error
}

func TestListContainers(t *testing.T) {
	tests := []struct {
		name       string
		provider   ecsiface.ECSAPI
		clusterArn string
		taskArn    string
		want       []*string
		wantErr    error
	}{
		{
			name: "EmptyContainers",
			provider: &mockECSClientContainers{
				Response: ecs.DescribeTasksOutput{},
			},
			clusterArn: "clusterArn",
			taskArn:    "taskArn",
			want:       nil,
			wantErr:    errors.New("error"),
		},
		{
			name: "SingleContainer",
			provider: &mockECSClientContainers{
				Response: ecs.DescribeTasksOutput{
					Tasks: []*ecs.Task{
						{
							Containers: []*ecs.Container{
								{
									Name: aws.String("Container"),
								},
							},
						},
					},
				},
			},
			clusterArn: "clusterArn",
			taskArn:    "taskArn",
			want:       []*string{aws.String("Container")},
		},
		{
			name: "MultipleContainers",
			provider: &mockECSClientContainers{
				Response: ecs.DescribeTasksOutput{
					Tasks: []*ecs.Task{
						{
							Containers: []*ecs.Container{
								{
									Name: aws.String("Container1"),
								},
								{
									Name: aws.String("Container2"),
								},
							},
						},
					},
				},
			},
			clusterArn: "clusterArn",
			taskArn:    "taskArn",
			want:       []*string{aws.String("Container1"), aws.String("Container2")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := listContainers(tt.clusterArn, tt.taskArn, tt.provider)

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected result: %v, got: %v", tt.want, result)
			}
		})
	}
}
