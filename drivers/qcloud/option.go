// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package qcloud

type Option func(*provider)

func WithApi(key, secret string) Option {
	return func(p *provider) {
		p.apikey = key
		p.apisecret = secret
	}
}

func WithRegion(region string) Option {
	return func(p *provider) {
		p.region = region
	}
}

func WithInstanceName(name string) Option {
	return func(p *provider) {
		p.name = name
	}
}
