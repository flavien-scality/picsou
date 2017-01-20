package main

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"testing"
)

type mockEC2Client struct {
	ec2iface.EC2API
}

func (m *mockEC2Client) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{Instances: []*ec2.Instance{
				{"salut", "ca"},
			}},
		},
	}
}

func TestListingInstances(t *testing.T) {
	// Setup test
	mockSvc := &mockEC2Client{}

	listInstances(mockSvc)
}
