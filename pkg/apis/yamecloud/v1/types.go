package v1

import (
	"net"
	"sort"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:resource:path=ippools,scope=Cluster
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IPPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec IPPoolSpec `json:"spec"`
}

type IPList []string
type IPPoolSpec struct {
	Ips map[string]IPList `json:"ips"`
}

func UInt32ToIP(intIP uint32) string {
	var bytes [4]byte
	bytes[0] = byte(intIP & 0xFF)
	bytes[1] = byte((intIP >> 8) & 0xFF)
	bytes[2] = byte((intIP >> 16) & 0xFF)
	bytes[3] = byte((intIP >> 24) & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}

func IPStringToUInt32(ip string) uint32 {
	bits := strings.Split(ip, ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])
	var sum uint32
	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)
	return sum
}

func (spec IPPoolSpec) Reuse() string {
	ips := make(IPList, 0)
	for _, i := range spec.Ips {
		ips = append(ips, i...)
	}
	sort.StringSlice(ips).Sort()
	reuse := func(start, end string) string {
		startIp := IPStringToUInt32(start)
		endIp := IPStringToUInt32(end)
		for startIp <= endIp {
			if !spec.FindIp(UInt32ToIP(startIp)) {
				return UInt32ToIP(startIp)
			}
			startIp++
		}
		return ""
	}
	return UInt32ToIP(IPStringToUInt32(reuse(ips[0], ips[len(ips)-1])) - 1)
}

func (spec IPPoolSpec) Last() string {
	ips := make(IPList, 0)
	for _, i := range spec.Ips {
		ips = append(ips, i...)
	}
	sort.StringSlice(ips).Sort()
	return ips[len(ips)-1]
}

func (spec IPPoolSpec) Release(id string) {
	delete(spec.Ips, id)
}

func (spec IPPoolSpec) FindIp(ip string) bool {
	for _, ipList := range spec.Ips {
		for _, i := range ipList {
			if i == ip {
				return true
			}
		}
	}
	return false
}

func (spec IPPoolSpec) Find(id, ip string) bool {
	if ips, ok := spec.Ips[id]; ok {
		for _, i := range ips {
			if i == ip {
				return true
			}
		}
	}
	return false
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IPPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []IPPool `json:"items"`
}

func init() {
	register(&IPPool{}, &IPPoolList{})
}
