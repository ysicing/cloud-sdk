// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cloudsdk

import (
	"context"
	"database/sql/driver"
	"errors"
)

// ProviderType specifies the hosting provider.
type ProviderType string

// Value converts the value to a sql string.
func (s ProviderType) Value() (driver.Value, error) {
	return string(s), nil
}

// Provider type enumeration.
const (
	ProviderAliyun = ProviderType("aliyun")
	ProviderQcloud = ProviderType("qcloud")
)

// ErrInstanceNotFound is returned when the requested
// instance does not exist in the cloud provider.
var ErrInstanceNotFound = errors.New("Not Found")

// A Provider represents a hosting provider, such as
// Digital Ocean and is responsible for server management.
type Provider interface {
	// Create creates a new server.
	Create(context.Context, InstanceCreateOpts) (*Instance, error)
	// Destroy destroys an existing server.
	Destroy(context.Context, *Instance) error
	// Show
	Show(ctx context.Context) []Instance
}

// An Instance represents a server instance
// (e.g Digital Ocean Droplet).
type Instance struct {
	Provider           ProviderType
	ID                 string
	Name               string
	Region             string
	InstanceType       string // 实例规格
	ImageID            string // 实例镜像
	InstanceChargeType string // 付费模式
}

// InstanceCreateOpts define soptional instructions for
// creating server instances.
type InstanceCreateOpts struct {
	Name    string
	CAKey   []byte
	CACert  []byte
	TLSKey  []byte
	TLSCert []byte
}

// InstanceError snapshots an error creating an instance
// with server logs.
type InstanceError struct {
	Err  error
	Logs []byte
}

// Error implements the error interface.
func (e *InstanceError) Error() string {
	return e.Err.Error()
}
