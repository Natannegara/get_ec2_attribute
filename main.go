package main

import (
	col "github.com/Natannegara/get_ec2_attribute/ec2collector"
)

func main() {
	col.Ec2Collector("samdev", "ap-southeast-3")
}
