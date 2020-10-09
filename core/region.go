package asgoptimizer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/service/autoscaling"
)

// Tag represents an Asg Tag: Key, Value
type Tag struct {
	Key   string
	Value string
}

// data structure that stores information about a region
type region struct {
	name string

	conf *Config
	// The key in this map is the instance type.
	// instanceTypeInformation map[string]instanceTypeInformation

	// instances instances

	enabledASGs []autoScalingGroup

	services connections

	// tagsToFilterASGsBy []Tag

	// wg sync.WaitGroup
}

func (r *region) enabled() bool {

	var enabledRegions []string

	if r.conf.EnabledRegions != "" {
		// Allow both space- and comma-separated values for the region list.
		csv := strings.Replace(r.conf.EnabledRegions, " ", ",", -1)
		enabledRegions = strings.Split(csv, ",")
		debug.Println("Enabled Regions: ", enabledRegions)
	} else {
		return true
	}

	for _, region := range enabledRegions {

		// glob matching for region names
		if match, _ := filepath.Match(region, r.name); match {
			return true
		}
	}

	return false
}

func (r *region) describeAllASGs() {
	svc := r.services.autoScaling
	// fmt.Println(svc.Config.Credentials.Get())
	pageNum := 0
	// var asgs []autoScalingGroup
	err := svc.DescribeAutoScalingGroupsPages(
		&autoscaling.DescribeAutoScalingGroupsInput{},
		func(page *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) bool {
			pageNum++
			for _, group := range page.AutoScalingGroups {
				asg := autoScalingGroup{
					Group: group,
					name:  *group.AutoScalingGroupName,
				}
				r.enabledASGs = append(r.enabledASGs, asg)
			}
			return true
		},
	)
	if err != nil {
		fmt.Println(err)
	}

}

func (r *region) getAllASGsInscope() {

}

func (r *region) processRegion() {

	logger.Println("Creating connections to the required AWS services in", r.name)
	r.services.connect(r.name)

	r.describeAllASGs()

	if r.enabledASGs != nil {
		for _, asg := range r.enabledASGs {
			logger.Println("Region: ", r.name, "ASGName: ", *asg.AutoScalingGroupName, "Instance Overrides: ", asg.getInstanceOverrideCount(), "Subnet Count: ", asg.getAZCount(), "ASG Spot Pools: ", asg.getSpotPoolCount())
		}
	}
}
