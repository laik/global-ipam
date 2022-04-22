package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/yametech/global-ipam/pkg/apis/yamecloud/v1"
	"github.com/yametech/global-ipam/pkg/log"
	"github.com/yametech/global-ipam/pkg/store"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var IPPool = schema.GroupVersionResource{Group: "yamecloud.io", Version: "v1", Resource: "ippools"}

func (s *Server) LastReservedIP(g *gin.Context) {
	defaultResponse := &store.LastReservedIPResponse{}
	// rangeId := g.PostForm("rangeId")
	ipPoolUnstructed, err := s.Interface.Resource(IPPool).Get(g.Request.Context(), store.GLOBAL_IPAM, metav1.GetOptions{})
	if err != nil {
		defaultResponse.Error = err
		g.JSON(http.StatusOK, defaultResponse)
		return
	}

	ippRuntime := &v1.IPPool{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(ipPoolUnstructed.UnstructuredContent(), ippRuntime)
	if err != nil {
		g.JSON(http.StatusOK, defaultResponse)
		return
	}

	defaultResponse.IP = ippRuntime.Spec.Last()

	g.JSON(http.StatusOK, defaultResponse)

}

func (s *Server) ReleaseByID(g *gin.Context) {
	id := g.PostForm("id")
	defaultResponse := &store.ReleaseResponse{}

	ipPoolUnstructed, err := s.Interface.Resource(IPPool).Get(g.Request.Context(), store.GLOBAL_IPAM, metav1.GetOptions{})
	if err != nil {
		defaultResponse.Error = err
		g.JSON(http.StatusOK, defaultResponse)
		return
	}

	ippRuntime := &v1.IPPool{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(ipPoolUnstructed.UnstructuredContent(), ippRuntime)
	if err != nil {
		g.JSON(http.StatusOK, defaultResponse)
		return
	}

	ippRuntime.Spec.Release(id)
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(ippRuntime)
	if err != nil {
		g.JSON(http.StatusOK, defaultResponse)
		return
	}

	ipPoolUnstructed.Object = obj

	_, err = s.Interface.Resource(IPPool).Update(g.Request.Context(), ipPoolUnstructed, metav1.UpdateOptions{})
	if err != nil {
		g.JSON(http.StatusOK, defaultResponse)
		return
	}

	defaultResponse.IsRelease = true

	g.JSON(http.StatusOK, defaultResponse)
}

func (s *Server) Reserve(g *gin.Context) {
	id := g.PostForm("id") // rook-ceph?
	ip := g.PostForm("ip") // 10.0.0.x

	defaultResponse := &store.ReserveResponse{Reserved: false}
	ipPoolUnstructed, err := s.Interface.Resource(IPPool).Get(g.Request.Context(), store.GLOBAL_IPAM, metav1.GetOptions{})
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			ipPoolUnstructed = &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "yamecloud.io/v1",
					"kind":       "IPPool",
					"metadata": map[string]interface{}{
						"name": store.GLOBAL_IPAM,
					},
					"spec": map[string]interface{}{
						"ips": make(map[string]string),
					},
				},
			}
			ipPoolUnstructed, err = s.Interface.Resource(IPPool).Create(g.Request.Context(), ipPoolUnstructed, metav1.CreateOptions{})
			if err != nil {
				defaultResponse.Error = fmt.Errorf("create ippool failed: %v", err)
				log.G(g.Request.Context()).Error(defaultResponse.Error)
				g.JSON(http.StatusOK, defaultResponse)
				return
			}
			goto RESERVE

		case errors.IsAlreadyExists(err):
			goto RESERVE

		}
		g.JSON(http.StatusOK, defaultResponse)
		return
	}

RESERVE:
	ippRuntime := &v1.IPPool{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(ipPoolUnstructed.UnstructuredContent(), ippRuntime)
	if err != nil {
		g.JSON(http.StatusOK, defaultResponse)
		return
	}

	if ippRuntime.Spec.Find(id, ip) {
		g.JSON(http.StatusOK, defaultResponse)
		return
	}

	ippRuntime.Spec.Ips[id] = append(ippRuntime.Spec.Ips[id], ip)

	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(ippRuntime)
	if err != nil {
		g.JSON(http.StatusOK, defaultResponse)
		return
	}

	ipPoolUnstructed.Object = obj

	_, err = s.Interface.Resource(IPPool).Update(g.Request.Context(), ipPoolUnstructed, metav1.UpdateOptions{})
	if err != nil {
		g.JSON(http.StatusOK, defaultResponse)
		return
	}

	defaultResponse.Reserved = true

	g.JSON(http.StatusOK, defaultResponse)
}
