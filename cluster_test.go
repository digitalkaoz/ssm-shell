package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"reflect"
	"testing"
)

type mockECSClientCluster struct {
	ecsiface.ECSAPI
	Response      ecs.ListClustersOutput
	PagedResponse ecs.ListClustersOutput
}

func (m *mockECSClientCluster) ListClusters(in *ecs.ListClustersInput) (*ecs.ListClustersOutput, error) {
	if m.Response.ClusterArns == nil {
		return nil, errors.New("error")
	}
	if in.NextToken != nil {
		return &m.PagedResponse, nil
	}
	return &m.Response, nil
}

func TestListClusters(t *testing.T) {
	tests := []struct {
		name     string
		want     []*string
		wantErr  error
		provider *mockECSClientCluster
	}{
		{
			name: "empty Clusters",
			provider: &mockECSClientCluster{
				Response: ecs.ListClustersOutput{},
			},
			want:    nil,
			wantErr: errors.New("error"),
		},
		{
			name: "single Cluster",
			provider: &mockECSClientCluster{
				Response: ecs.ListClustersOutput{ClusterArns: []*string{aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/test")}},
			},
			want:    []*string{aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/test")},
			wantErr: nil,
		},
		{
			name: "multiple Clusters",
			provider: &mockECSClientCluster{
				Response: ecs.ListClustersOutput{ClusterArns: []*string{aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/test"), aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/test2")}},
			},
			want:    []*string{aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/test"), aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/test2")},
			wantErr: nil,
		},
		{
			name: "paged Response",
			provider: &mockECSClientCluster{
				Response:      ecs.ListClustersOutput{ClusterArns: []*string{aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/test"), aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/test2")}, NextToken: aws.String("token")},
				PagedResponse: ecs.ListClustersOutput{ClusterArns: []*string{aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/test"), aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/test3")}},
			},
			want:    []*string{aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/test"), aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/test2"), aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/test"), aws.String("arn:aws:ecs:us-west-2:123456789012:cluster/test3")},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := listClusters(tt.provider)

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected result: %v, got: %v", tt.want, result)
			}
		})
	}
}
