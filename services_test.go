package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

type mockECSClientServices struct {
	ecsiface.ECSAPI
	Response      ecs.ListServicesOutput
	PagedResponse ecs.ListServicesOutput
}

func (m *mockECSClientServices) ListServices(in *ecs.ListServicesInput) (*ecs.ListServicesOutput, error) {
	if m.Response.ServiceArns == nil {
		return nil, errors.New("error")
	}
	if in.NextToken != nil {
		return &m.PagedResponse, nil
	}

	return &m.Response, nil
}

func TestListServices(t *testing.T) {
	testCases := []struct {
		name     string
		cluster  string
		want     []*string
		wantErr  error
		provider *mockECSClientServices
	}{
		{
			name:    "Result with single service",
			cluster: "cluster1",
			provider: &mockECSClientServices{
				Response: ecs.ListServicesOutput{
					ServiceArns: []*string{aws.String("service1")},
				},
			},
			want: []*string{aws.String("service1")},
		},
		{
			name:    "Result with multiple services",
			cluster: "cluster2",
			provider: &mockECSClientServices{
				Response: ecs.ListServicesOutput{
					ServiceArns: []*string{aws.String("service1"), aws.String("service2")},
				},
			},
			want: []*string{aws.String("service1"), aws.String("service2")},
		},
		{
			name:    "Empty result",
			cluster: "cluster3",
			provider: &mockECSClientServices{
				Response: ecs.ListServicesOutput{},
			},
			wantErr: errors.New("error"),
		},
		{
			name:    "paged result",
			cluster: "cluster3",
			provider: &mockECSClientServices{
				Response: ecs.ListServicesOutput{
					NextToken:   aws.String("token"),
					ServiceArns: []*string{aws.String("service1"), aws.String("service2")},
				},
				PagedResponse: ecs.ListServicesOutput{
					ServiceArns: []*string{aws.String("service3"), aws.String("service4")},
				},
			},
			want: []*string{aws.String("service1"), aws.String("service2"), aws.String("service3"), aws.String("service4")},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			services, err := listServices(tt.cluster, tt.provider)

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(services, tt.want) {
				t.Errorf("unexpected services result. want: %v, got: %v", tt.want, services)
			}

		})
	}
}
