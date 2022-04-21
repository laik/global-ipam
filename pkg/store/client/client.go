package client

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/yametech/global-ipam/pkg/allocator"
	"github.com/yametech/global-ipam/pkg/store"
)

var _ store.Store = &Client{}

type CniClient struct {
	*resty.Client
}

type Client struct {
	mutex sync.Mutex
	cli   *CniClient
}

func New(name string, IPAMConfig *allocator.IPAMConfig) (store.Store, error) {
	return &Client{
		mutex: sync.Mutex{},
		cli:   NewCniClient(store.UNIX_SOCK_PATH),
	}, nil
}

// Close implements store.Store
func (*Client) Close() error { return nil }

// LastReservedIP implements store.Store
func (c *Client) LastReservedIP(rangeID string) (net.IP, error) {
	resp, err := c.cli.R().Get(fmt.Sprintf("/last-reserved-ip/%s", rangeID))
	if err != nil {
		return nil, err
	}
	_ = resp
	// json.Unmarshal(resp.Body(), &ip)

	return nil, nil
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
// N.B. This function eats errors to be tolerant and release as much as possible
func (*Client) ReleaseByID(id string) error {
	// 	key := s.EtcdKeyPrefix + "/used/"
	// 	resp, err := s.EtcdClient.Get(context.TODO(), key, clientv3.WithPrefix())
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if len(resp.Kvs) > 0 {
	// 		for _, kv := range resp.Kvs {
	// 			if string(kv.Value) == id {
	// 				_, err = s.EtcdClient.Delete(context.TODO(), string(kv.Key))
	// 				if err != nil {
	// 					return err
	// 				}
	// 			}
	// 		}
	// 	}
	// 	return nil
	return nil
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

func NewCniClient(socketAddress string) *CniClient {
	transport := http.Transport{
		Dial: func(_, _ string) (net.Conn, error) {
			return net.Dial("unix", store.UNIX_SOCK_PATH)
		},
	}

	// Create a Resty Client
	client := resty.NewWithClient(&http.Client{Transport: &transport}).
		SetScheme("http").
		SetHostURL("http://dummy")

	return &CniClient{Client: client}
}
