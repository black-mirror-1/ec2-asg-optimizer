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
		// if a.MixedInstancesPolicy.LaunchTemplate != nil {
		// 	if a.MixedInstancesPolicy.LaunchTemplate.Overrides != nil {
		// 		fmt.Println("ASGName: ", a.AutoScalingGroupName, ", #AvailabilityZone: ", len(a.AvailabilityZones), ", InstanceOverrides= ", len(a.MixedInstancesPolicy.LaunchTemplate.Overrides))
		// 	}
		// } else {
		// 	fmt.Println("ASGName: ", a.AutoScalingGroupName, ", No Launch Template")
		// }
		fmt.Println("ASGName: ", a.AutoScalingGroupName, ", #AvailabilityZone: ", len(a.AvailabilityZones), ", InstanceOverrides= ", len(a.MixedInstancesPolicy.LaunchTemplate.Overrides))
		return len(a.MixedInstancesPolicy.LaunchTemplate.Overrides)
	}
	fmt.Println("ASGName: ", a.AutoScalingGroupName, ", No MIG")
	return 1
}
