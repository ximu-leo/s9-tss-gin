package tssgin

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-gonic/gin"
	"github.com/ximu-leo/s9-tss-gin/router"
	"github.com/ximu-leo/s9-tss-gin/services"
)

type GinHttpServer struct {
	Host     string
	Port     int
	Registry *router.Registry
	stopped  atomic.Bool
}

func NewGinHttpServer(httpHost string, httpPort int) (*GinHttpServer, error) {
	serviceManager, err := services.NewManager()
	if err != nil {
		log.Error("new manager fail", "err", err)
		return nil, err
	}
	registry := router.NewRegistry(serviceManager)
	return &GinHttpServer{
		Host:     httpHost,
		Port:     httpPort,
		Registry: registry,
	}, nil
}

func (ms *GinHttpServer) Start(ctx context.Context) error {

	r := gin.Default()

	ms.Registry.Register(r)

	ginAddr := fmt.Sprintf("%s:%d", ms.Host, ms.Port)
	err := r.Run(ginAddr)
	if err != nil {
		log.Error("run fail", "err", err)
		return err
	}
	return nil
}

func (ms *GinHttpServer) Stop(ctx context.Context) error {
	ms.stopped.Store(true)
	return nil
}

func (ms *GinHttpServer) Stopped() bool {
	return ms.stopped.Load()
}
