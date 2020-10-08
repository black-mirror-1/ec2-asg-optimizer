package asgoptimizer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

var logger, debug *log.Logger

// var conf Config

func setupLogging(cfg *Config) {
	logger = log.New(cfg.LogFile, "", cfg.LogFlag)

	if os.Getenv("ENABLE_DEBUG") == "true" || cfg.EnableDebug {
		fmt.Println("Debug logs enabled")
		debug = log.New(cfg.LogFile, "", cfg.LogFlag)
	} else {
		debug = log.New(ioutil.Discard, "", 0)
	}

}

// func isMixedInstancePolicy(a autoscaling.Group) bool {
// 	if a.MixedInstancesPolicy != nil {
// 		return true
// 	}
// 	return false
// }

func connectEC2(region string) *ec2.EC2 {

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	return ec2.New(sess,
		aws.NewConfig().WithRegion(region))
}

// getRegions generates a list of AWS regions.
func getRegions(ec2conn ec2iface.EC2API) ([]string, error) {
	var output []string

	logger.Println("Scanning for available AWS regions")

	resp, err := ec2conn.DescribeRegions(&ec2.DescribeRegionsInput{})

	if err != nil {
		logger.Println(err.Error())
		return nil, err
	}

	debug.Println(resp)

	for _, r := range resp.Regions {

		if r != nil && r.RegionName != nil {
			debug.Println("Found region", *r.RegionName)
			output = append(output, *r.RegionName)
		}
	}
	return output, nil
}

//Run is the entry point to the asgadvisor module
func Run(cfg *Config) {

	setupLogging(cfg)

	debug.Println(*cfg)

	// use this only to list all the other regions
	ec2Conn := connectEC2("us-west-2")

	allRegions, err := getRegions(ec2Conn)

	if err != nil {
		logger.Println(err.Error())
		return
	}

	processRegions(allRegions, cfg)

}

func processRegions(regions []string, cfg *Config) {

	var wg sync.WaitGroup

	for _, r := range regions {

		wg.Add(1)
		r := region{name: r, conf: cfg}

		go func() {

			if r.enabled() {
				logger.Printf("Enabled to run in %s, processing region.\n", r.name)
				r.processRegion()
			} else {
				debug.Println("Not enabled to run in", r.name)
				debug.Println("List of enabled regions:", cfg.EnabledRegions)
			}

			wg.Done()
		}()
	}
	wg.Wait()
}
