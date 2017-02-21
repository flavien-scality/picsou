package stats

import (
 . "github.com/smartystreets/goconvey/convey"
 "github.com/aws/aws-sdk-go/service/ec2"
 "github.com/aws/aws-sdk-go/service/ec2/ec2iface"
 "errors"
 "strconv"
 "testing"
)

var FuncShouldPanic bool

type mockEC2Client struct {
 ec2iface.EC2API
}

var errDescribe = errors.New("describeError")

func (m *mockEC2Client) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
 // Handle should panic or not
 err := errDescribe
 if !FuncShouldPanic {
 err = nil
 }
 return &ec2.DescribeInstancesOutput{
 Reservations: []*ec2.Reservation{
  {
  Instances: []*ec2.Instance{
   // Set 2 instances for a simple reservation
   {}, {},
  },
  },
 },
 }, err
}

func TestListingInstances(t *testing.T) {
 // Setup test
 mockSvc := &mockEC2Client{}
 srv := Stats{Service: EC2{Regions: []string{"us-west-1",},},}

 Convey("Testing instances listing", t, func() {
 Convey("Should be equal to '2'", func() {
  FuncShouldPanic = false
  So(srv.listInstances(mockSvc), ShouldEqual, 2)
 })
 Convey("Should panic when aws API call fails", func() {
  FuncShouldPanic = true
  So(func () { srv.listInstances(mockSvc) }, ShouldPanicWith, errDescribe)
 })
 })
}


func testGetState(k int64, expected string) {
  Convey("Testing with: "+strconv.FormatInt(k, 10)+" - "+expected, func() {
    So(GetState(k), ShouldEqual, expected)
  })
}

func TestGetStateMapping(t *testing.T) {
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
