package main

import (
	//"github.com/aws/aws-sdk-go/service/ec2"
	//"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
)

/*
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
*/

func testGetState(k int64, expected string) {
	Convey("Testing with: "+strconv.FormatInt(k, 10)+" - "+expected, func() {
		So(GetState(k), ShouldEqual, expected)
	})
}

func TestGetState(t *testing.T) {
	maps := make(map[int64]string)
	maps[0] = "pending"
	maps[16] = "running"
	maps[32] = "shutting-down"
	maps[48] = "terminated"
	maps[64] = "stopping"
	maps[80] = "stopped"
	maps[42] = ""

	Convey("Testing GetState function", t, func() {
		for k, v := range maps {
			testGetState(k, v)
		}
	})
}
