package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"reflect"
	"testing"
)

type mockedEc2Client struct {
	ec2iface.EC2API
	Response      ec2.DescribeInstancesOutput
	PagedResponse ec2.DescribeInstancesOutput
}

func (m *mockedEc2Client) DescribeInstances(in *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	if m.Response.Reservations == nil {
		return nil, errors.New("error")
	}
	if in.NextToken != nil {
		return &m.PagedResponse, nil
	}
	return &m.Response, nil
}

func TestListInstances(t *testing.T) {
	tests := []struct {
		name          string
		want          []*string
		wantErr       error
		response      ec2.DescribeInstancesOutput
		pagedResponse ec2.DescribeInstancesOutput
		provider      *mockedEc2Client
	}{
		{
			name: "ValidInstanceList",
			provider: &mockedEc2Client{
				Response: ec2.DescribeInstancesOutput{
					Reservations: []*ec2.Reservation{
						{
							Instances: []*ec2.Instance{
								{
									InstanceId: aws.String("i-abc123"),
									Tags: []*ec2.Tag{
										{
											Key:   aws.String("Name"),
											Value: aws.String("TestInstance"),
										},
									},
								},
							},
						},
					},
				},
			},
			want: []*string{aws.String("TestInstance (i-abc123)")},
		},
		{
			name:     "DescribeInstancesError",
			provider: &mockedEc2Client{},
			wantErr:  errors.New("error"),
		},
		{
			name: "MultiplePages",
			provider: &mockedEc2Client{
				Response: ec2.DescribeInstancesOutput{
					NextToken: aws.String("next"),
					Reservations: []*ec2.Reservation{
						{
							Instances: []*ec2.Instance{
								{
									InstanceId: aws.String("i-123"),
									Tags: []*ec2.Tag{
										{
											Key:   aws.String("Name"),
											Value: aws.String("TestInstance1"),
										},
									},
								},
							},
						},
					},
				},
				PagedResponse: ec2.DescribeInstancesOutput{
					Reservations: []*ec2.Reservation{
						{
							Instances: []*ec2.Instance{
								{
									InstanceId: aws.String("i-456"),
									Tags: []*ec2.Tag{
										{
											Key:   aws.String("Name"),
											Value: aws.String("TestInstance2"),
										},
									},
								},
							},
						},
					},
				},
			},
			want: []*string{aws.String("TestInstance1 (i-123)"), aws.String("TestInstance2 (i-456)")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := listInstances(tt.provider)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("listInstances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(response, tt.want) {
				t.Errorf("listInstances() = %v, want %v", response, tt.want)
			}
		})
	}
}
