// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package qcloud

import (
	"context"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func (p *provider) setup(ctx context.Context) error {
	return p.setupKeypair(ctx)
}

func (p *provider) setupKeypair(ctx context.Context) error {
	p.sshkey = p.listKeypair(ctx)
	return nil
}

func (p *provider) listKeypair(ctx context.Context) []string {
	request := cvm.NewDescribeKeyPairsRequest()
	response, err := p.getClient().DescribeKeyPairs(request)
	if err != nil {
		return nil
	}
	if *response.Response.TotalCount == 0 {
		request := cvm.NewCreateKeyPairRequest()
		_, err := p.getClient().CreateKeyPair(request)
		if err != nil {
			return nil
		}
		return p.listKeypair(ctx)
	}
	var keys []string
	for _, k := range response.Response.KeyPairSet {
		keys = append(keys, *k.KeyId)
	}
	return keys
}

func (p *provider) listSubnet(ctx context.Context) (*cvm.VirtualPrivateCloud, string) {
	request := vpc.NewDescribeSubnetsRequest()
	if len(p.subnetname) != 0 {
		request.Filters = []*vpc.Filter{
			&vpc.Filter{
				Name:   common.StringPtr("subnet-name"),
				Values: common.StringPtrs([]string{p.subnetname}),
			},
		}
	}
	response, err := p.getVpcClient().DescribeSubnets(request)

	if err != nil || *response.Response.TotalCount < 1 {
		return &cvm.VirtualPrivateCloud{
			VpcId:    common.StringPtr("DEFAULT"),
			SubnetId: common.StringPtr("DEFAULT"),
		}, fmt.Sprintf("%v-1", p.region)
	}
	return &cvm.VirtualPrivateCloud{
		VpcId:    response.Response.SubnetSet[0].VpcId,
		SubnetId: response.Response.SubnetSet[0].SubnetId,
	}, *response.Response.SubnetSet[0].Zone
}
