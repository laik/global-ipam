package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/yametech/global-ipam/pkg/log"
	loglogrus "github.com/yametech/global-ipam/pkg/log/logrus"
	"github.com/yametech/global-ipam/pkg/signals"
	"github.com/yametech/global-ipam/pkg/store/server"
)

func main() {
	fmt.Printf("i am cni server\n")

	stopCh := signals.SetupSignalHandler()
	ctx, cancel := context.WithCancel(context.Background())

	log.L = loglogrus.FromLogrus(logrus.NewEntry(logrus.StandardLogger()))

	go func() {
		<-stopCh
		cancel()
	}()

	if err := server.NewServer(ctx).Start(); err != nil {
		panic(err)
	}

}
