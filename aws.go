package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var ec2Client ec2iface.EC2API
var ecsClient *ecs.ECS

func SetEcsClient(client *ecs.ECS) {
	ecsClient = client
}

func SetEc2Client(client ec2iface.EC2API) {
	ec2Client = client
}

func createAwsSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func createEc2Client() (ec2iface.EC2API, error) {
	if ec2Client != nil {
		return ec2Client, nil
	}
	sess, err := createAwsSession()
	if err != nil {
		return nil, err
	}

	return ec2.New(sess), nil
}

func createEcsService() (*ecs.ECS, error) {
	if ecsClient != nil {
		return ecsClient, nil
	}
	sess, err := createAwsSession()
	if err != nil {
		return nil, err
	}
	svc := ecs.New(sess)

	ecsClient = svc
	return svc, nil
}
