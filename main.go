package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

func listInstances(svc ec2iface.EC2API) int {
	count := 0

	// Call the DescribeInstances Operation
	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	// resp has all of the response data, pull out instance IDs:
	for _, res := range resp.Reservations {
		count += len(res.Instances)
	}
	return count
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	// Create an EC2 service object in the "us-west-2" region
	// Note that you can also configure your region globally by
	// exporting the AWS_REGION environment variable
	regions := []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "ca-central-1", "eu-west-1", "eu-west-2", "ap-northeast-1", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2", "ap-south-1", "sa-east-1"}
	for i := range regions {
		region := regions[i]
	        svc := ec2.New(sess, &aws.Config{Region: aws.String(region)})
		fmt.Println("> Region name: ", region)
		fmt.Println("  - Number of instances: ", listInstances(svc))
	}
}
