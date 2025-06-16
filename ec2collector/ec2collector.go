package ec2collector

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	identityAWS "github.com/aws/aws-sdk-go-v2/service/sts"
)

type Data struct {
	Name            string
	Priority        string
	PIC             string
	AccountId       string
	BackupRetention string
	BackupStatus    string
}

func setupAWSConf(profile string, region string) aws.Config {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(profile), config.WithRegion(region))
	if err != nil {
		panic(err)
	}
	return cfg
}

// use this function to get all raw attributes of object in AWS (currently: EC2 & AWS Account Info)
func getAttribute(profile string, region string) []Data {

	awsConf := setupAWSConf(profile, region)

	// make Client Session for EC2 Instance
	ec2Client := *ec2.NewFromConfig(awsConf)
	ec2GetIns, err := ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{})
	if err != nil {
		panic(err)
	}

	// make Client Session for Information AWS Account
	infoAWSClient := identityAWS.NewFromConfig(awsConf)
	res, err := infoAWSClient.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		panic(err)
	}

	// get service EC2
	var listData []Data
	for _, ec2Instance := range ec2GetIns.Reservations {
		// get the instance Describes (many)
		var resourceData Data
		resourceData.AccountId = *res.Account
		for _, ec2InstanceDesc := range ec2Instance.Instances {
			for _, tag := range ec2InstanceDesc.Tags {
				if *tag.Key == "priority" {
					resourceData.Priority = *tag.Value
				}
				if *tag.Key == "Name" {
					resourceData.Name = *tag.Value
				}
				if *tag.Key == "pic" {
					resourceData.PIC = *tag.Value
				}
				if *tag.Key == "backup_retention" {
					resourceData.BackupRetention = *tag.Value
				}
				if *tag.Key == "backup_status" {
					resourceData.BackupStatus = *tag.Value
				}
			}
		}
		listData = append(listData, resourceData)
	}
	return listData
}

// use this function to fill empty value of attribute, fill it with default value
func Ec2Collector(profile string, region string) []Data {
	listData := getAttribute(profile, region)
	for _, v := range listData {
		if strings.TrimSpace(v.PIC) == "" {
			v.PIC = "undefined"
		}
		if strings.TrimSpace(v.Priority) == "" {
			v.Priority = "low"
		}
		if strings.TrimSpace(v.BackupRetention) == "" {
			v.BackupRetention = "false"
		}
		if strings.TrimSpace(v.BackupStatus) == "" {
			v.BackupStatus = "false"
		}
		listData = append(listData, v)
	}
	return listData
}
