// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package qcloud

import (
	"context"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/ysicing/cloudsdk"
)

type CVMFilterType string

const (
	CVMFilterVPC          = CVMFilterType("vpc-id")
	CVMFilterSubnet       = CVMFilterType("subnet-id")
	CVMFilterInstanceID   = CVMFilterType("instance-id")
	CVMFilterInstanceName = CVMFilterType("instance-name")
	CVMFilterChargeType   = CVMFilterType("instance-charge-type")
	CVMFilterDefault      = CVMFilterType("none")
)

func (p *provider) Show(ctx context.Context) []cloudsdk.Instance {
	request := cvm.NewDescribeInstancesRequest()
	//if k != CVMFilterDefault {
	//	request.Filters = []*cvm.Filter {
	//		&cvm.Filter {
	//			Name: common.StringPtr(string(k)),
	//			Values: common.StringPtrs([]string{ v }),
	//		},
	//	}
	//}
	response, err := p.getClient().DescribeInstances(request)
	if err != nil {
		return nil
	}
	var instances []cloudsdk.Instance
	for _, i := range response.Response.InstanceSet {
		// 仅列出竞价实例
		if *i.InstanceChargeType != "SPOTPAID" {
			continue
		}
		instances = append(instances, cloudsdk.Instance{
			Provider:           cloudsdk.ProviderQcloud,
			ID:                 *i.InstanceId,
			Name:               *i.InstanceName,
			Region:             p.region,
			InstanceType:       *i.InstanceType,
			ImageID:            *i.ImageId,
			InstanceChargeType: *i.InstanceChargeType,
		})
	}
	return instances
}
