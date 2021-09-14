// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package qcloud

import (
	"context"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/ysicing/cloudsdk"
)

func (p *provider) check(ctx context.Context, id string) bool {
	request := cvm.NewDescribeInstancesRequest()
	request.Filters = []*cvm.Filter{
		&cvm.Filter{
			Name:   common.StringPtr("tag-key"),
			Values: common.StringPtrs([]string{"k3s"}),
		},
		&cvm.Filter{
			Name:   common.StringPtr("tag-value"),
			Values: common.StringPtrs([]string{"k3s"}),
		},
		&cvm.Filter{
			Name:   common.StringPtr("instance-id"),
			Values: common.StringPtrs([]string{id}),
		},
	}
	response, err := p.getClient().DescribeInstances(request)
	if err != nil {
		return false
	}
	if *response.Response.TotalCount != 1 {
		return false
	}
	return *response.Response.InstanceSet[0].InstanceChargeType == "SPOTPAID"
}

func (p *provider) Destroy(ctx context.Context, instance *cloudsdk.Instance) error {
	if !p.check(ctx, instance.ID) {
		return fmt.Errorf("只允许删除竞价实例")
	}
	request := cvm.NewTerminateInstancesRequest()
	request.InstanceIds = common.StringPtrs([]string{instance.ID})
	_, err := p.getClient().TerminateInstances(request)
	if err != nil {
		return err
	}
	return nil
}
