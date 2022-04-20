package client

import (
	"net"
	"sync"

	"github.com/yametech/global-ipam/pkg/allocator"
	"github.com/yametech/global-ipam/pkg/store"
)

var _ store.Store = &Client{}

type Client struct {
	mutex sync.Mutex
}

func New(name string, IPAMConfig *allocator.IPAMConfig) (store.Store, error) {
	return &Client{
		mutex: sync.Mutex{},
	}, nil
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
func (c *Client) Lock() error {
	c.mutex.Lock()
	return nil
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
func (c *Client) Unlock() error {
	c.mutex.Unlock()
	return nil
}
