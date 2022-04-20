package client

import (
	"net"

	"github.com/yametech/global-ipam/pkg/allocator"
	"github.com/yametech/global-ipam/pkg/store"
)

var _ store.Store = &Client{}

type Client struct{}

func New(name string, IPAMConfig *allocator.IPAMConfig) (*store.Store, error) {
	return nil, nil
}

// Close implements store.Store
func (*Client) Close() error {
	panic("unimplemented")
}

// LastReservedIP implements store.Store
func (*Client) LastReservedIP(rangeID string) (net.IP, error) {
	panic("unimplemented")
}

// Lock implements store.Store
func (*Client) Lock() error {
	panic("unimplemented")
}

// Release implements store.Store
func (*Client) Release(ip net.IP) error {
	panic("unimplemented")
}

// ReleaseByID implements store.Store
func (*Client) ReleaseByID(id string) error {
	panic("unimplemented")
}

// Reserve implements store.Store
func (*Client) Reserve(id string, ip net.IP, rangeID string) (bool, error) {
	panic("unimplemented")
}

// Unlock implements store.Store
func (*Client) Unlock() error {
	panic("unimplemented")
}
