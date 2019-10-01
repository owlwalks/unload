package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/google/logger"
)

var (
	lbv2          *elasticloadbalancingv2.Client
	setupLbv2Once sync.Once
)

var (
	errSetupLbv2 = fmt.Errorf("lbv2 is not setup")
)

func setupLbv2() error {
	setupLbv2Once.Do(func() {
		cfg, err := external.LoadDefaultAWSConfig()
		if err != nil {
			logger.Errorln(err)
			return
		}
		// work out aws current region
		meta := ec2metadata.New(cfg)
		cfg.Region, err = meta.Region()
		if err != nil {
			logger.Errorln(err)
			return
		}
		lbv2 = elasticloadbalancingv2.New(cfg)
		// watch and reconcile
		go watchLbv2()
	})
	if lbv2 == nil {
		return errSetupLbv2
	}
	return nil
}

func regPod(arn *string, ip string, port int64) {
	if err := setupLbv2(); err != nil {
		logger.Warningln(err)
		return
	}
	req := lbv2.RegisterTargetsRequest(&elasticloadbalancingv2.RegisterTargetsInput{
		TargetGroupArn: arn,
		Targets: []elasticloadbalancingv2.TargetDescription{{
			Id:   &ip,
			Port: &port,
		}},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	_, err := req.Send(ctx)
	if err != nil {
		logger.Errorln(err)
	}
}

// this will remove out-of-synced unhealthy targets
func reconcile(arn *string) {
	if err := setupLbv2(); err != nil {
		logger.Warningln(err)
		return
	}
	des := lbv2.DescribeTargetHealthRequest(&elasticloadbalancingv2.DescribeTargetHealthInput{
		TargetGroupArn: arn,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	res, err := des.Send(ctx)
	if err != nil {
		logger.Errorln(err)
		return
	}
	var targets []elasticloadbalancingv2.TargetDescription
	for _, desc := range res.TargetHealthDescriptions {
		if desc.TargetHealth.State == elasticloadbalancingv2.TargetHealthStateEnumUnhealthy {
			targets = append(targets, *desc.Target)
		}
	}
	if len(targets) > 0 {
		dereg := lbv2.DeregisterTargetsRequest(&elasticloadbalancingv2.DeregisterTargetsInput{
			TargetGroupArn: arn,
			Targets:        targets,
		})
		if _, err := dereg.Send(ctx); err != nil {
			logger.Errorln(err)
		}
	}
}

func watchLbv2() {
	// max interval of nlb healthcheck
	for range time.Tick(30 * time.Second) {
		reconcile(targetGroupArn)
	}
}
