package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"reflect"
	"testing"
)

type mockECSClientTasks struct {
	ecsiface.ECSAPI
	Response      ecs.ListTasksOutput
	PagedResponse ecs.ListTasksOutput
	Error         error
}

func (m *mockECSClientTasks) ListTasks(in *ecs.ListTasksInput) (*ecs.ListTasksOutput, error) {
	if m.Response.TaskArns == nil {
		return nil, errors.New("error")
	}
	if in.NextToken != nil {
		return &m.PagedResponse, m.Error
	}

	return &m.Response, m.Error
}

func TestListTasks(t *testing.T) {
	tests := []struct {
		name       string
		provider   ecsiface.ECSAPI
		clusterArn string
		serviceArn string
		want       []*string
		wantErr    error
	}{
		{
			name: "EmptyTasks",
			provider: &mockECSClientTasks{
				Response: ecs.ListTasksOutput{},
			},
			clusterArn: "clusterArn",
			serviceArn: "taskArn",
			want:       nil,
			wantErr:    errors.New("error"),
		},
		{
			name: "SingleTask",
			provider: &mockECSClientTasks{
				Response: ecs.ListTasksOutput{
					TaskArns: []*string{aws.String("taskArn1")},
				},
			},
			clusterArn: "clusterArn",
			serviceArn: "taskArn",
			want:       []*string{aws.String("taskArn1")},
		},
		{
			name: "MultipleTasks",
			provider: &mockECSClientTasks{
				Response: ecs.ListTasksOutput{
					TaskArns: []*string{aws.String("taskArn1"), aws.String("taskArn2")},
				},
			},
			clusterArn: "clusterArn",
			serviceArn: "taskArn",
			want:       []*string{aws.String("taskArn1"), aws.String("taskArn2")},
		},
		{
			name: "Paged Results",
			provider: &mockECSClientTasks{
				Response: ecs.ListTasksOutput{
					TaskArns:  []*string{aws.String("taskArn1")},
					NextToken: aws.String("nextToken"),
				},
				PagedResponse: ecs.ListTasksOutput{
					TaskArns: []*string{aws.String("taskArn2")},
				},
			},
			clusterArn: "clusterArn",
			serviceArn: "taskArn",
			want:       []*string{aws.String("taskArn1"), aws.String("taskArn2")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := listTasks(tt.clusterArn, tt.serviceArn, tt.provider)

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected result: %v, got: %v", tt.want, result)
			}
		})
	}
}
