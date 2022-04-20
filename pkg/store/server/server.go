package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yametech/global-ipam/pkg/store"
)

type Server struct {
	*http.Server
}

func NewServer(ctx context.Context) *Server {
	s := &Server{
		&http.Server{
			Handler: gin.Default(),
		},
	}
	go func() {
		for range ctx.Done() {
			if err := s.Shutdown(ctx); err != nil {
				fmt.Printf("shutdown server error: %v", err)
			}
		}
	}()

	return s
}

func (s *Server) Start() error {
	if err := os.Remove(store.UNIX_SOCK_PATH); err != nil {
		return err
	}

	defer func() {
		if err := os.Remove(store.UNIX_SOCK_PATH); err != nil {
			fmt.Printf("remove unix socket error: %v", err)
		}
	}()
	
	unixListener, err := net.Listen("unix", store.UNIX_SOCK_PATH)
	if err != nil {
		return err
	}

	return s.Serve(unixListener)
}
