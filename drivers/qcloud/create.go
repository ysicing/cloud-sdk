// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package qcloud

import (
	"context"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/ysicing/cloudsdk"
	"log"
)

// Create creates an OpenStack instance
func (p *provider) Create(ctx context.Context, opts cloudsdk.InstanceCreateOpts) (*cloudsdk.Instance, error) {
	p.init.Do(func() {
		if err := p.setup(ctx); err != nil {
			log.Println(err)
		}
	})
	instance, err := p.create(ctx, opts)
	return instance, err
}

func (p *provider) create(ctx context.Context, opts cloudsdk.InstanceCreateOpts) (*cloudsdk.Instance, error) {
	vpc, zone := p.listSubnet(ctx)
	request := cvm.NewRunInstancesRequest()
	request.InstanceChargeType = common.StringPtr("SPOTPAID") // 创建竞价实例
	request.Placement = &cvm.Placement{
		Zone: common.StringPtr(zone), // 可用区
	}
	request.InstanceType = common.StringPtr("SA2.SMALL1") // 实例规格 SA2.SMALL1 1核1g
	request.ImageId = common.StringPtr("img-h1yvvfw1")    // 镜像ID debian 10.2
	request.SystemDisk = &cvm.SystemDisk{                 // 系统盘
		DiskType: common.StringPtr("CLOUD_PREMIUM"), // 普通高效云盘
		DiskSize: common.Int64Ptr(50),               // 50GB
	}
	request.VirtualPrivateCloud = vpc                     // 网络设置
	request.InternetAccessible = &cvm.InternetAccessible{ // 流量设置
		InternetChargeType:      common.StringPtr("TRAFFIC_POSTPAID_BY_HOUR"),
		InternetMaxBandwidthOut: common.Int64Ptr(100),
		PublicIpAssigned:        common.BoolPtr(true),
	}
	request.InstanceName = common.StringPtr(p.name)
	// 登录设置
	if len(p.sshkey) == 0 {
		request.LoginSettings = &cvm.LoginSettings{
			Password: common.StringPtr("aeghiedelohn8pu0quihaeghee6Ahtai"),
		}
	} else {
		request.LoginSettings = &cvm.LoginSettings{
			KeyIds: common.StringPtrs(p.sshkey),
		}
	}
	// 默认服务开启 云安全、云监控
	request.EnhancedService = &cvm.EnhancedService{
		SecurityService: &cvm.RunSecurityServiceEnabled{
			Enabled: common.BoolPtr(true),
		},
		MonitorService: &cvm.RunMonitorServiceEnabled{
			Enabled: common.BoolPtr(true),
		},
		AutomationService: &cvm.RunAutomationServiceEnabled{
			Enabled: common.BoolPtr(true),
		},
	}
	// 标签
	request.TagSpecification = []*cvm.TagSpecification{
		&cvm.TagSpecification{
			ResourceType: common.StringPtr("instance"), // 固定不可修改
			Tags: []*cvm.Tag{
				&cvm.Tag{
					Key:   common.StringPtr("k3s"),
					Value: common.StringPtr("k3s"),
				},
			},
		},
	}
	// 竞价设置
	request.InstanceMarketOptions = &cvm.InstanceMarketOptionsRequest{
		MarketType: common.StringPtr("spot"),
		SpotOptions: &cvm.SpotMarketOptions{
			MaxPrice:         common.StringPtr("0.5"), // 竞价价格
			SpotInstanceType: common.StringPtr("one-time"),
		},
	}
	// cloudinit 初始化脚本
	request.UserData = common.StringPtr("xxxx")
	// debug 模式
	// request.DryRun = common.BoolPtr(true)
	response, err := p.getClient().RunInstances(request)
	if err != nil {
		return nil, err
	}

	instance := cloudsdk.Instance{}
	//if *request.DryRun {
	//	log.Println(*response.Response.RequestId)
	//	return &instance, nil
	//}
	instance.ID = *response.Response.InstanceIdSet[0]
	return &instance, nil
}

func (p *provider) get(ctx context.Context) (*cloudsdk.Instance, error) {
	request := cvm.NewDescribeInstancesRequest()
	request.Filters = []*cvm.Filter{
		//&cvm.Filter {
		//	Name: common.StringPtr("tag-key"),
		//	Values: common.StringPtrs([]string{ "xxxx" }),
		//},
		&cvm.Filter{
			Name: common.StringPtr(p.name),
		},
	}
	response, err := p.getClient().DescribeInstances(request)
	if err != nil {
		return nil, err
	}
	if *response.Response.TotalCount != 1 {
		return nil, fmt.Errorf("实例不唯一")
	}
	instance := &cloudsdk.Instance{}
	instance.ID = *response.Response.InstanceSet[0].InstanceId
	instance.InstanceType = *response.Response.InstanceSet[0].InstanceType
	instance.ImageID = *response.Response.InstanceSet[0].ImageId
	instance.InstanceChargeType = *response.Response.InstanceSet[0].InstanceChargeType
	return instance, nil
}
