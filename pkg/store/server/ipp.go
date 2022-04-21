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
	rangeId := g.PostForm("rangeId")
	g.JSON(http.StatusOK, rangeId)

}

func (s *Server) ReleaseByID(g *gin.Context) {
	id := g.PostForm("id")
	g.JSON(http.StatusOK, id)
}

func (s *Server) Reserve(g *gin.Context) {
	id := g.PostForm("id") // rook-ceph?
	// rangeId := g.PostForm("rangeId") // index
	ip := g.PostForm("ip") // 10.0.0.x

	defaultResponse := &store.ReserveResponse{Reserved: false}
	ippUnstructed, err := s.Interface.Resource(IPPool).Get(g.Request.Context(), id, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			ippUnstructed = &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "yamecloud.io/v1",
					"kind":       "IPPool",
					"metadata": map[string]interface{}{
						"name": id,
					},
					"spec": map[string]interface{}{
						"ips": make([]string, 0),
					},
				},
			}
			ippUnstructed, err = s.Interface.Resource(IPPool).Create(g.Request.Context(), ippUnstructed, metav1.CreateOptions{})
			if err != nil {
				defaultResponse.Error = fmt.Errorf("create ippool failed: %v", err)
				log.G(g.Request.Context()).Error(defaultResponse.Error)
				g.JSON(http.StatusOK, defaultResponse)
				return
			}
			goto RESERVE
		}
		g.JSON(http.StatusOK, defaultResponse)
		return
	}

RESERVE:
	ippRuntime := &v1.IPPool{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(ippUnstructed.UnstructuredContent(), ippRuntime)
	if err != nil {
		g.JSON(http.StatusOK, defaultResponse)
		return
	}

	ippRuntime.Spec.Ips = append(ippRuntime.Spec.Ips, ip)

	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(ippRuntime)
	if err != nil {
		g.JSON(http.StatusOK, defaultResponse)
		return
	}
	ippUnstructed.Object = obj

	_, err = s.Interface.Resource(IPPool).Update(g.Request.Context(), ippUnstructed, metav1.UpdateOptions{})
	if err != nil {
		g.JSON(http.StatusOK, defaultResponse)
		return
	}

	defaultResponse.Reserved = true

	g.JSON(http.StatusOK, defaultResponse)
}
