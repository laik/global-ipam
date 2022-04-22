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

func (c *Client) GetByID(id, ip string) (net.IP, error) {
	return nil, nil
}

// LastReservedIP implements store.Store
func (c *Client) LastReservedIP(rangeId string) (net.IP, error) {
	r := &store.LastReservedIPResponse{}
	_, err := c.cli.R().
		SetHeader("Content-Type", "application/json").
		SetFormData(map[string]string{"rangeId": rangeId}).
		SetResult(r).
		Post(fmt.Sprintf("/last-reserved-ip"))
	if err != nil {
		return nil, err
	}
	return net.ParseIP(r.IP), nil
}

// Lock implements store.Store
func (c *Client) Lock() error {
	c.mutex.Lock()
	return nil
}

// Release implements store.Store
func (c *Client) Release(ip net.IP) error {

	_, err := c.cli.R().
		SetHeader("Content-Type", "application/json").
		SetFormData(map[string]string{"ip": ip.String()}).
		Post("/release")
	if err != nil {
		return err
	}
	return nil
}

// ReleaseByID implements store.Store
// N.B. This function eats errors to be tolerant and release as much as possible
func (c *Client) ReleaseByID(id string) error {
	_, err := c.cli.R().
		SetHeader("Content-Type", "application/json").
		SetFormData(map[string]string{"id": id}).
		Post("/release-by-id")
	if err != nil {
		return err
	}
	return nil
}

// Reserve implements store.Store
func (c *Client) Reserve(id string, ip net.IP, rangeId string) (bool, error) {
	r := &store.ReserveResponse{}
	_, err := c.cli.R().
		SetHeader("Content-Type", "application/json").
		SetFormData(map[string]string{
			"id":      id,
			"ip":      ip.String(),
			"rangeId": rangeId,
		}).
		SetResult(r).
		Post("/reserve")
	if err != nil {
		return false, err
	}
	return r.Reserved, nil
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
