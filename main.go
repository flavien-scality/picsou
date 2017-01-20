package main

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type Instance struct {
	Id string
}

type Reservation struct {
	Instances []Instance
}

type EC2 struct {
	Regions []string
	Reservations []Reservation
}

type Cloud struct {
	Service EC2
}

func getInstances(reservation int, instances []*ec2.Instance) {
	fmt.Println("reservation: ", reservation, " instances: ", len(instances))
}

func listInstances(svc ec2iface.EC2API) int {
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
		getInstances(reservation, res.Instances)
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

	var wg sync.WaitGroup
	nums := make(chan int)
	regions := []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "ca-central-1", "eu-west-1", "eu-west-2", "ap-northeast-1", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2", "ap-south-1", "sa-east-1"}
	srv := &Cloud{EC2{Regions: regions}}
	for i := range srv.Service.Regions {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			region := srv.Service.Regions[i]
			svc := ec2.New(sess, &aws.Config{Region: aws.String(region)})
			nm := listInstances(svc)
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
}
