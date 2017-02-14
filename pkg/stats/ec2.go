package stats

import (
	// "bytes"
	// "sync"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
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
	CloudWatch *CloudWatch
	Reservations []*Reservation
	ReservationsRunning []int
	ReservationsUsage float64
	Volumes []*ec2.Volume
	VolumesUsage float64
}

// Stats represente the differents statistics for the differents reservations
// and instances from the ec2 account
type Stats struct {
	// Differents reservations and instances
	Service map[string]*EC2
	Regions []string
}

type CloudWatch struct {
	Client *cloudwatch.CloudWatch
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
	srv := &Stats{Regions: regions, Service: make(map[string]*EC2)}
	for _, region := range regions {
		newService := &EC2{
		        Client: ec2.New(sess, &aws.Config{Region: aws.String(region)}),
			CloudWatch: &CloudWatch{Client: cloudwatch.New(sess, &aws.Config{Region: aws.String(region)}),},
			Reservations: []*Reservation{},
			Volumes: []*ec2.Volume{},
		}
		srv.Service[region] = newService.getReservations().getRunningInstances().getReservationsUsage().getVolumes().getVolumesUsage()
	}
	return srv
}

func (s *EC2) RunningInstances(instances []*ec2.Instance) (int, int) {
	var running, total int
	for _, instance := range instances {
		if GetState(*instance.State.Code) == "running" {
			running += 1
		}
		total += 1
	}
	return running, total
}

func (s *EC2) getRunningInstances() *EC2 {
	var reservationsRunning []int
	for i, reservation := range s.Reservations {
		if reservation.InstancesRunning != 0 {
			reservationsRunning = append(reservationsRunning, i)
		}
	}
	s.ReservationsRunning = reservationsRunning
	return s
}

func (s *EC2) GetRunningRatio() float64 {
	if len(s.Reservations) == 0 {
		return float64(0)
	} else {
		return float64(len(s.ReservationsRunning)) / float64(len(s.Reservations))
	}
}

func (s *EC2) getReservations() *EC2 {
	running := 0
	count := 0

	// Call the DescribeInstances Operation
	resp, err := s.Client.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	// resp has all of the response data, pull out instance IDs:
	for _, reservation := range resp.Reservations {
		instancesRunning, instancesTotal := s.RunningInstances(reservation.Instances)
		if instancesTotal != 0 && instancesRunning == instancesTotal {
			s.Reservations = append(s.Reservations, &Reservation{Instances: reservation.Instances, InstancesRunning: instancesRunning, InstancesTotal: instancesTotal})
		} else {
			s.Reservations = append(s.Reservations, &Reservation{Instances: reservation.Instances, InstancesRunning: instancesRunning, InstancesTotal: instancesTotal})
		}
		running += instancesRunning
		count += instancesTotal
	}
	return s
}

func (s *EC2) GetVolumesSize() int64 {
	var count int64
	for _, volume := range s.Volumes {
		count += *volume.Size
	}
	return count
}

func (s *EC2) getVolumes() *EC2 { 
	// var filterName = "availability-zone"
	volumesOutput, err := s.Client.DescribeVolumes(nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	s.Volumes = append(s.Volumes, volumesOutput.Volumes...)
	return s
}

func (s *EC2) getReservationsUsage() *EC2 {
	var total, nb float64
	startTime := time.Now().AddDate(0, 0, -1)
	delta, err := time.ParseDuration("-10m")
	if err != nil {
		fmt.Println("error during time parsing: ", err)
	}
	endTime := time.Now().Add(delta)
	for _, running := range s.ReservationsRunning {
			for _, instance := range s.Reservations[running].Instances {
			res, err := s.CloudWatch.Client.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
				MetricName: aws.String("CPUUtilization"),
				Namespace: aws.String("AWS/EC2"),
				Dimensions: []*cloudwatch.Dimension {
					&cloudwatch.Dimension {
						Name: aws.String("InstanceId"),
						Value: instance.InstanceId,
					},
				},
				StartTime: aws.Time(startTime),
				EndTime: aws.Time(endTime),
				Period: aws.Int64(3600),
				Statistics: []*string{ aws.String("Sum"), aws.String("SampleCount"), },
			})
			if err != nil {
				fmt.Println("error during cloudwatch getMetrics: ", err)
			}
			for _, datapoint := range res.Datapoints {
				if *datapoint.Sum > 0.0 {
					nb += *datapoint.SampleCount
					total += *datapoint.Sum
				}
			}
		}
	}
	if int(nb) != 0 {
		s.ReservationsUsage = total/nb
	} else {
		s.ReservationsUsage = 0.0
	}
	return s
}

func (s *EC2) getVolumesUsage() *EC2 {
	var total, nb float64
	startTime := time.Now().AddDate(0, 0, -1)
	delta, err := time.ParseDuration("-10m")
	if err != nil {
		fmt.Println("error during time parsing: ", err)
	}
	endTime := time.Now().Add(delta)
	for _, volume := range s.Volumes {
		res, err := s.CloudWatch.Client.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
			MetricName: aws.String("VolumeIdleTime"),
			Namespace: aws.String("AWS/EBS"),
			Dimensions: []*cloudwatch.Dimension {
				&cloudwatch.Dimension {
					Name: aws.String("VolumeId"),
					Value: volume.VolumeId,
				},
			},
			StartTime: aws.Time(startTime),
			EndTime: aws.Time(endTime),
			Period: aws.Int64(3600),
			Statistics: []*string{ aws.String("Sum"), aws.String("SampleCount") },
		})
		if err != nil {
			fmt.Println("error during cloudwatch getMetrics: ", err)
		}
		for _, datapoint := range res.Datapoints {
			if *datapoint.Sum > 0.0 {
				nb += *datapoint.SampleCount
				total += *datapoint.Sum
			}
		}
	}
	if int(nb) != 0 {
		s.VolumesUsage = 100.0 - (total/nb/3600.0)*100
	} else {
		s.VolumesUsage = 0.0
	}
	return s

}
