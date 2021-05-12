package etcd

import (
	"github.com/containernetworking/cni/pkg/types"
	"github.com/yametech/global-ipam/pkg/allocator"
	"net"
	"testing"
)

func Test_Lock(t *testing.T) {
	_, ipNet, err := net.ParseCIDR("10.16.0.0/16")

	_IPAMConfig := &allocator.IPAMConfig{
		Range: &allocator.Range{
			RangeStart: net.IPv4(10, 16, 0, 1),
			RangeEnd:   net.IPv4(10, 16, 0, 254),
			Subnet:     types.IPNet(*ipNet),
		},
		EtcdConfig: &allocator.EtcdConfig{
			EtcdURL: "http://10.200.100.200:42379",
		},
	}

	_, _ = err, _IPAMConfig

	_, err = New("", _IPAMConfig)
	if err != nil {
		t.Error(err)
	}

	//err = s1.Lock()
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//s2, err := New("", _IPAMConfig)
	//err = s2.Lock()
	//if err != nil {
	//	t.Error(err)
	//}

}
