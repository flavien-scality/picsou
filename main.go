package main

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

var regions = []string{
	"us-east-1",
	"us-east-2",
	"us-west-1",
	"us-west-2",
	"ca-central-1",
	"eu-west-1",
	"eu-west-2",
	"ap-northeast-1",
	"ap-northeast-2",
	"ap-southeast-1",
	"ap-southeast-2",
	"ap-south-1",
	"sa-east-1",
}

type Instance struct {
	Id         string
	State      uint
	LaunchTime string
	Type       string
}

type Reservation struct {
	Instances []Instance
}

type EC2 struct {
	Regions      []string
	Reservations []Reservation
}

type Stats struct {
	Service EC2
}

func (s Stats) getState(code int64) string {
	switch code {
	case 0:
		return "pending"
	case 16:
		return "running"
	case 32:
		return "shutting-down"
	case 48:
		return "terminated"
	case 64:
		return "stopping"
	case 80:
		return "stopped"
	}
	return ""
}

func New(sess *session.Session, regions []string) *Stats {
	var wg sync.WaitGroup
	nums := make(chan int)
	srv := &Stats{EC2{Regions: regions}}
	for i := range srv.Service.Regions {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			region := srv.Service.Regions[i]
			svc := ec2.New(sess, &aws.Config{Region: aws.String(region)})
			nm := srv.listInstances(svc)
			fmt.Println("Region Name: ", region, " | Number of instances: ", nm)
			nums <- nm
		}(i)
	}
	go func() {
		wg.Wait()
		close(nums)
	}()
	sum := 0
	for num := range nums {
		sum += num
	}
	fmt.Println("done\ntotal instances: ", sum)
	return srv
}

func (s Stats) getInstances(reservation int, instances []*ec2.Instance) {
	// fmt.Println("reservation: ", reservation, " instances: ", len(instances))
	for _, instance := range instances {
		fmt.Print("id: ", *instance.InstanceId)
		// fmt.Print(" | LaunchTime: ", instance.LaunchTime)
		// fmt.Print(" | ClientToken: ", *instance.ClientToken)
		fmt.Println(" | KeyName: ", *instance.KeyName)
		// fmt.Println(" | State: ", s.getState(*instance.State.Code))
	}
}

func (s Stats) listInstances(svc ec2iface.EC2API) int {
	count := 0

	// Call the DescribeInstances Operation
	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	// resp has all of the response data, pull out instance IDs:
	reservation := 1
	for _, res := range resp.Reservations {
		//fmt.Println("\n\n** instances: ", res.Instances)
		// fmt.Println(res)
		fmt.Print("Owner: ", *res.OwnerId)
		fmt.Println(" | ReservationId: ", *res.ReservationId)
		s.getInstances(reservation, res.Instances)
		count += len(res.Instances)
		reservation += 1
	}
	return count
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	_ = New(sess, regions)
}
