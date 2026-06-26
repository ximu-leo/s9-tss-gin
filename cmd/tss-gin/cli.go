package main

import (
	"context"

	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
	tssgin "github.com/ximu-leo/s9-tss-gin"
	"github.com/ximu-leo/s9-tss-gin/common/cliapp"
	"github.com/ximu-leo/s9-tss-gin/config"
	flags2 "github.com/ximu-leo/s9-tss-gin/flags"
)

// type LifecycleAction func(ctx *cli.Context, close context.CancelCauseFunc) (Lifecycle, error)
// 写一个符合这个模板的“具体函数”
func runGinHttp(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	log.Info("starting tss gin server")
	cfg := config.NewConfig(ctx)
	// 返回自定义的结构体 GinHttpServer，包含主机/端口/路由
	return tssgin.NewGinHttpServer(cfg.Host, cfg.Port)
}

func NewCli() *cli.App {
	flags := flags2.Flags
	return &cli.App{
		Version:              "0.0.1",
		Description:          "An gin services with rpc",
		EnableBashCompletion: true,
		Commands: []*cli.Command{

			{
				Name:        "gin",
				Flags:       flags,
				Description: "Run gin http services",
				// LifecycleCmd 接收的参数类型是 LifecycleAction
				// 我们直接把 runGinHttp 这个函数名传进去！
				Action: cliapp.LifecycleCmd(runGinHttp),
			},

			{
				Name:        "version",
				Description: "Show project version",
				Action: func(ctx *cli.Context) error {
					cli.ShowVersion(ctx)
					return nil
				},
			},
		},
	}
}
