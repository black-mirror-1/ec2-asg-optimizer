package asgoptimizer

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type autoScalingGroup struct {
	*autoscaling.Group
	name            string
	poolStrength    int64
	recommendations []recommendation
}

type recommendation struct {
	severity string
	message  string
	// higher the number higher the importance
	priority int64
}

func (a *autoScalingGroup) isMixedInstancePolicy() bool {
	if a.MixedInstancesPolicy != nil {
		return true
	}
	return false
}

func (a *autoScalingGroup) getInstanceOverrideCount() int {
	if a.isMixedInstancePolicy() {
		return len(a.MixedInstancesPolicy.LaunchTemplate.Overrides)
	}
	fmt.Println("ASGName: ", a.AutoScalingGroupName, ", No MIG")
	return 1
}

func (a *autoScalingGroup) isAZFlexible() bool {
	if a.getAZCount() > 1 {
		return true
	}
	return false
}

func (a *autoScalingGroup) getAZCount() int {
	return len(a.AvailabilityZones)
}

func (a *autoScalingGroup) getSpotPoolCount() int {
	return a.getInstanceOverrideCount() * a.getAZCount()
}
