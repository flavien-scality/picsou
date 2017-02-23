package stats

import (
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
  CloudWatch *CloudWatch
  Instances []*ec2.Instance
  InstancesRunning []int
  InstancesUsage float64
  Volumes []*ec2.Volume
  VolumesUsage float64
  // VolumesTotalSize int
  // VolumesUsedSize int
}

type User struct {
  Reservations map[int][]int
  Volumes []int
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
  Users map[string]*User
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
      Users: make(map[string]*User),
    }
    newService = newService.getReservations().getRunningInstances().getReservationsUsage().getVolumes().getVolumesUsage().getUsers().getUsages()
    if len(newService.Reservations) != 0 {
      srv.Service[region] = newService
    }
  }
  return srv
}

func (s *EC2) getUsages() *EC2{
  for _, r := range s.Reservations {
    r.getInstancesUsage()
  }
  return s
}

func (s *EC2) RunningInstances(instances []*ec2.Instance) []int {
  var running []int
  for i, instance := range instances {
    if GetState(*instance.State.Code) == "running" {
      running = append(running, i)
    }
  }
  return running
}

func (s *EC2) getRunningInstances() *EC2 {
  var reservationsRunning []int
  for i, reservation := range s.Reservations {
    if len(reservation.InstancesRunning) != 0 {
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
    return float64(len(s.ReservationsRunning)) / float64(len(s.Reservations)) * 100.0
  }
}

func (s *Reservation) GetInstancesRatio() float64 {
  if len(s.Instances) == 0 {
    return float64(0)
  } else {
    return float64(len(s.InstancesRunning)) / float64(len(s.Instances)) * 100.0
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
    instancesRunning := s.RunningInstances(reservation.Instances)
    if len(reservation.Instances) != 0 {
      s.Reservations = append(s.Reservations, &Reservation{Instances: reservation.Instances, InstancesRunning: instancesRunning, CloudWatch: s.CloudWatch})
    } else {
      s.Reservations = append(s.Reservations, &Reservation{Instances: reservation.Instances, InstancesRunning: instancesRunning, CloudWatch: s.CloudWatch})
    }
    running += len(instancesRunning)
    count += len(reservation.Instances)
  }
  return s
}

func (s *EC2) GetVolumesSize(targets []int) int64 {
  var count int64
  if targets != nil {
    for _, target := range targets {
      count += *s.Volumes[target].Size
    }
  } else {
    for _, volume := range s.Volumes {
      count += *volume.Size
    }
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

func (s *EC2) GetReservationsUsage(targets []int) float64 {
  var total, nb float64

  startTime := time.Now().AddDate(0, 0, -1)
  delta, err := time.ParseDuration("-10m")
  if err != nil {
    fmt.Println("error during time parsing: ", err)
  }
  endTime := time.Now().Add(delta)

  for _, target := range targets {
    for _, instance := range s.Reservations[target].Instances {
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
    return total / nb
  }
  return 0.0
}

func (s *EC2) getReservationsUsage() *EC2 {
  s.ReservationsUsage = s.GetReservationsUsage(s.ReservationsRunning)
  return s
}

func (s *EC2) getUsers() *EC2 {
  var username string
  users := make(map[string]*User)
  for r, reservation := range s.Reservations {
    for i, instance := range reservation.Instances {
      username = *instance.KeyName
      if _, ok := users[username]; !ok {
        users[username] = &User{
          Reservations: make(map[int][]int),
          Volumes: []int{},
        }
      }
      users[username].Reservations[r] = append(users[username].Reservations[r], i)
    }
  }
  for i, volume := range s.Volumes {
    if len(volume.Attachments) != 0 {
      for _, v := range volume.Attachments {
        params := &ec2.DescribeInstancesInput{
          Filters: []*ec2.Filter{
            {
              Name: aws.String("instance-id"),
              Values: []*string{v.InstanceId},
            },
          },
        }
        resp, err := s.Client.DescribeInstances(params)
        if err != nil {
          fmt.Println("error when describing instances for volume: ", err)
        }
        for _, res := range resp.Reservations {
          for _, ins := range res.Instances {
            username = *ins.KeyName
            users[username].Volumes = append(users[username].Volumes, i)
          }
        }
      }
    }
  }
  s.Users = users
  return s
}

func (s *EC2) GetUsersRunning(user string) int {
  count := 0
  for r, _ := range s.Users[user].Reservations {
    for _, ri := range s.ReservationsRunning {
      if r == ri {
        count += 1
      }
    }
  }
  return count
}

func (s *Stats) Divide(a, b int) float64 {
  if b == 0 {
    return float64(0)
  } else {
    return float64(a) / float64(b) * 100
  }
}

func (s *User) ReservationsKeys() []int {
  var keys []int
  for key, _ := range s.Reservations {
    keys = append(keys, key)
  }
  return keys
}
