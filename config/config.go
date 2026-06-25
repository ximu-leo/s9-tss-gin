package config

import (
	"github.com/ximu-leo/s9-tss-gin/flags"
)

type Config struct {
	Host string
	Port int
}

func NewConfig(ctx *cli.Context) Config {
	return Config{
		Host: ctx.String(flags.HttpHostFlag.Name),
		Port: ctx.Int(flags.HttpPortFlag.Name),
	}
}
