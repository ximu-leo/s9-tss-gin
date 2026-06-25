package main

import (
	"context"

	"github.com/ethereum/go-ethereum/log"
	tssgin "github.com/ximu-leo/s9-tss-gin"
	"github.com/ximu-leo/s9-tss-gin/common/cliapp"
	"github.com/ximu-leo/s9-tss-gin/config"
	flags2 "github.com/ximu-leo/s9-tss-gin/flags"
)

func runGinHttp(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	log.Info("starting tss gin server")
	cfg := config.NewConfig(ctx)
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
				Action:      cliapp.LifecycleCmd(runGinHttp),
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
