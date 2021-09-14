// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package qcloud

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/ysicing/cloudsdk"
	"sync"
)

type provider struct {
	init       sync.Once
	region     string
	apikey     string
	apisecret  string
	name       string
	sshkey     []string
	subnetname string
}

func (p *provider) getClient() *cvm.Client {
	credential := common.NewCredential(
		p.apikey,
		p.apisecret,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	client, _ := cvm.NewClient(credential, p.region, cpf)
	return client
}

func (p *provider) getVpcClient() *vpc.Client {
	credential := common.NewCredential(
		p.apikey,
		p.apisecret,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"
	client, _ := vpc.NewClient(credential, p.region, cpf)
	return client
}

func New(opts ...Option) (cloudsdk.Provider, error) {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}
