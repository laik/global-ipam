package v1

import (
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

func (spec IPPoolSpec) Last(id string) string {
	if ips, ok := spec.Ips[id]; ok {
		return ips[len(ips)-1]
	}
	return ""
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
