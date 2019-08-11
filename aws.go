package main

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"k8s.io/klog"
)

var (
	lbv2          *elasticloadbalancingv2.Client
	setupLbv2Once sync.Once
)

func setupLbv2() {
	setupLbv2Once.Do(func() {
		cfg, err := external.LoadDefaultAWSConfig()
		if err != nil {
			klog.Errorln(err)
			return
		}
		// work out aws current region
		meta := ec2metadata.New(cfg)
		cfg.Region, err = meta.Region()
		if err != nil {
			klog.Errorln(err)
			return
		}
		lbv2 = elasticloadbalancingv2.New(cfg)
	})
}

func regPod(targetGroupArn string, ip string, port int64) {
	setupLbv2()
	if lbv2 == nil {
		klog.Warningln("lbv2 is not setup")
		return
	}
	req := lbv2.RegisterTargetsRequest(&elasticloadbalancingv2.RegisterTargetsInput{
		TargetGroupArn: &targetGroupArn,
		Targets: []elasticloadbalancingv2.TargetDescription{{
			Id:   &ip,
			Port: &port,
		}},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := req.Send(ctx)
	if err != nil {
		klog.Errorln(err)
	}
}

func deregPod(targetGroupArn string, ip string, port int64) {
	setupLbv2()
	if lbv2 == nil {
		klog.Warningln("lbv2 is not setup")
		return
	}
	req := lbv2.DeregisterTargetsRequest(&elasticloadbalancingv2.DeregisterTargetsInput{
		TargetGroupArn: &targetGroupArn,
		Targets: []elasticloadbalancingv2.TargetDescription{{
			Id:   &ip,
			Port: &port,
		}},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := req.Send(ctx)
	if err != nil {
		klog.Errorln(err)
	}
}
