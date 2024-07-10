package ec2collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Data struct {
	Name     string
	Priority string `default: "low"`
	PIC      string `default: "undefined"`
}

func setupAWSConf(profile string, region string) aws.Config {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(profile), config.WithRegion(region))
	if err != nil {
		panic(err)
	}
	return cfg
}

func getAttribute(profile string, region string) []Data {
	ec2Client := *ec2.NewFromConfig(setupAWSConf(profile, region))
	ec2GetIns, err := ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{})
	if err != nil {
		panic(err)
	}

	// get service EC2
	var listData []Data
	for _, ec2Instance := range ec2GetIns.Reservations {
		// get the instance Describes (many)
		var resourceData Data
		for _, ec2InstanceDesc := range ec2Instance.Instances {
			for _, tag := range ec2InstanceDesc.Tags {
				if *tag.Key == "priority" {
					resourceData.Priority = *tag.Value
				}
				if *tag.Key == "Name" {
					resourceData.Name = *tag.Value
				}
				if *tag.Key == "PIC" {
					resourceData.PIC = *tag.Value
				}
			}
		}
		listData = append(listData, resourceData)
	}
	return listData
}

func Ec2Collector(profile string, region string) []Data {
	listData := getAttribute(profile, region)
	for _, v := range listData {
		if v.PIC == "" {
			v.PIC = "undefined"
		}
		if v.Priority == "" {
			v.Priority = "low"
		}
		listData = append(listData, v)
	}
	return listData
}
