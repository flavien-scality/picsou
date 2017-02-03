package stats

import (
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// Regions defines all amazon regions to check for ec2 usage
var Regions = []string{
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

// Instance Represente a minimuse the informations needed from ec2.Instance
//
// For more documentation see:
// https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/Instance
type Instance struct {
	// Unique id of the ec2.Instance
	ID string
	// Curremt state of the ec2.Instance
	State uint
	// Time when the instance have been lanch
	LaunchTime time.Time
	// Type of the ec2 instance
	Type string
}

// Reservation is a collection of EC2 instances started as part of the same launch request.
//
// For more documentation see:
// http://docs.aws.amazon.com/general/latest/gr/glos-chap.html#reservation
type Reservation struct {
	// Instances started
	Instances []Instance
}

// EC2 represente a collection of Regions and reservations from aws-ec2
//
// For more documentation see:
// https://aws.amazon.com/ec2/
type EC2 struct {
	// Array of region of the reservations
	Regions []string
	// Array of reservations for the ec2 account
	Reservations []Reservation
}

// Stats represente the differents statistics for the differents reservations
// and instances from the ec2 account
type Stats struct {
	// Differents reservations and instances
	Service EC2
}

// GetState match the status code with the representing status string
func GetState(code int64) string {
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

// New will get from aws all the reservations and then the instances
// from that informations the function will compute the appropriate stats
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
			rn, nm := srv.listInstances(svc)
			fmt.Println("Region Name: ", region, " | Number of instances: ", rn, "/", nm)
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

func (s Stats) getInstances(reservation int, instances []*ec2.Instance) int {
	// fmt.Println("reservation: ", reservation, " instances: ", len(instances))
	running := 0
	// t := time.Now()
	for _, instance := range instances {
		// fmt.Print("id: ", *instance.InstanceId)
		// fmt.Print(" | deltatime: ", t.Sub(*instance.LaunchTime))
		// fmt.Print(
		// 	" | LaunchTime: ",
		// 	instance.LaunchTime.Hour(),
		// 	":",
		// 	instance.LaunchTime.Minute(),
		// )
		// fmt.Print(" | ClientToken: ", *instance.ClientToken)
		// fmt.Println(" | KeyName: ", *instance.KeyName)
		// fmt.Println(" | State: ", GetState(*instance.State.Code))
		if GetState(*instance.State.Code) == "running" {
			running += 1
		}
	}
	return running
}

func (s Stats) listInstances(svc ec2iface.EC2API) (int, int) {
	running := 0
	count := 0

	// Call the DescribeInstances Operation
	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	// resp has all of the response data, pull out instance IDs:
	reservation := 1
	for _, res := range resp.Reservations {
		// fmt.Println("\n\n** instances: ", res.Instances)
		// fmt.Println(res)
		// fmt.Print("Owner: ", *res.OwnerId)
		// fmt.Println(" | ReservationId: ", *res.ReservationId)
		runs := s.getInstances(reservation, res.Instances)
		ttl := len(res.Instances)
		count += ttl
		running += runs
		// If runs != ttl || runs != 0 then send warning incomplete reservation shutdown
		fmt.Println("Reservation: ", reservation, " | instances: ", runs, "/", ttl)
		reservation++
	}
	return running, count
}
