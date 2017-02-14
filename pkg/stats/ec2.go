package stats

import (
	// "bytes"
	// "sync"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	//"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// Regions defines all amazon regions to check for ec2 usage
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
	Instances []*ec2.Instance
	InstancesRunning int
	InstancesTotal int
	// Volumes []*ec2.Volume
	// VolumesTotalSize int
	// VolumesUsedSize int
	// VolumesUsageRatio float64
}

// EC2 represente a collection of Regions and reservations from aws-ec2
//
// For more documentation see:
// https://aws.amazon.com/ec2/
type EC2 struct {
	Client *ec2.EC2
	// Array of reservations for the ec2 account
	Reservations []*Reservation
	ReservationsRunning []int
	RunningRatio float64
	VolumesTotal int
}

// Stats represente the differents statistics for the differents reservations
// and instances from the ec2 account
type Stats struct {
	// Differents reservations and instances
	Service map[string]*EC2
	Regions []string
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
func New(sess *session.Session) *Stats {
	// var wg sync.WaitGroup
	// nums := make(chan int)
	srv := &Stats{Regions: regions, Service: make(map[string]*EC2)}
	for _, region := range regions {
		// wg.Add(1)
		newService := &EC2{
		        Client: ec2.New(sess, &aws.Config{Region: aws.String(region)}),
			// defer wg.Done()
			// fmt.Printf("Region Name: %s | Number of instances: %d/%d\n", region, rn, nm)
	                // fmt.Printf("running ratio: %f\n", srv.Service[region].RunningRatio)
			// nums <- nm
			Reservations: []*Reservation{},
		}
		srv.Service[region] = newService
		srv = srv.getReservations(region)
	}
	// func() {
	// 	wg.Wait()
	// 	close(nums)
	// }()
	// sum := 0
	// for num := range nums {
	// 	sum += num
	// }
	// fmt.Printf("Total instances: %d", sum)
	return srv.getVolumes()
}

func (s *Stats) RunningInstances(instances []*ec2.Instance) (int, int) {
	// fmt.Println("reservation: ", reservation, " instances: ", len(instances))
	var running, total int
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
		total += 1
	}
	return running, total
}

func (s *EC2) setRunningRatio() (float64, []int) {
	var reservationsRunning []int
	var running, count = 0, 0
	for i, reservation := range s.Reservations {
		if reservation.InstancesRunning != 0 {
			reservationsRunning = append(reservationsRunning, i)
		}
		if s.Reservations[i].InstancesTotal != 0 {
			running += reservation.InstancesRunning
			count += reservation.InstancesTotal
		}
	}
	if count != 0 {
		return float64(running) * 100.0 / float64(count), reservationsRunning
	}
	return float64(0), reservationsRunning
}

func (s *Stats) getReservations(region string) *Stats {
	running := 0
	count := 0

	// Call the DescribeInstances Operation
	resp, err := s.Service[region].Client.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	// resp has all of the response data, pull out instance IDs:
	for _, reservation := range resp.Reservations {
		// fmt.Println("\n\n** instances: ", res.Instances)
		// fmt.Println(res)
		// fmt.Print("Owner: ", *res.OwnerId)
		// fmt.Println(" | ReservationId: ", *res.ReservationId)
		instancesRunning, instancesTotal := s.RunningInstances(reservation.Instances)
		// fmt.Println("instancesTotal: ", instancesTotal, " instancesRunning: ", instancesRunning, " instances: ", len(reservation.Instances))
		if instancesTotal != 0 && instancesRunning == instancesTotal {
			s.Service[region].Reservations = append(s.Service[region].Reservations, &Reservation{Instances: reservation.Instances, InstancesRunning: instancesRunning, InstancesTotal: instancesTotal})
		} else {
			s.Service[region].Reservations = append(s.Service[region].Reservations, &Reservation{Instances: reservation.Instances, InstancesRunning: instancesRunning, InstancesTotal: instancesTotal})
		}
		running += instancesRunning
		count += instancesTotal
		// If instancesRunning != instancesTotal || instancesRunning != 0 then send warning incomplete reservation shutdown
		// fmt.Printf("Reservation: %d | instances: %d/%d\n", reservationNb, instancesRunning, instancesTotal)
	}
	return s
}

func (s *Stats) getVolumes() *Stats {
	var filterName = "availability-zone"
	fmt.Println("len service: ", len(s.Service))
	for region, _ := range s.Service {
		fmt.Println(region)
		var filter = &ec2.Filter{
			Name: aws.String(filterName),
			Values: []*string{ aws.String(region), },
		}
		volumesOutput, err := s.Service[region].Client.DescribeVolumes(&ec2.DescribeVolumesInput{Filters: []*ec2.Filter{ filter}})
		if err != nil {
			fmt.Println("Error: ", err)
		}
		s.Service[region].VolumesTotal = len(volumesOutput.Volumes)
	}
	return s
}
