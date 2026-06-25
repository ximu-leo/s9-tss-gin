package flags

import "github.com/urfave/cli/v2"

const evnVarPrefix = "GIN"

func prefixEnvVars(name string) []string {
	return []string{evnVarPrefix + "_" + name}
}

var (
	HttpHostFlag = &cli.StringFlag{
		Name:     "http-host",
		Usage:    "The host of the http",
		EnvVars:  prefixEnvVars("HTTP_HOST"),
		Required: true,
	}
	HttpPortFlag = &cli.IntFlag{
		Name:     "http-port",
		Usage:    "The port of the http",
		EnvVars:  prefixEnvVars("HTTP_PORT"),
		Required: true,
	}
)
var requireFlags = []cli.Flag{
	HttpHostFlag,
	HttpPortFlag,
}

var optionalFlags = []cli.Flag{}

func init() {
	Flags = append(requireFlags, optionalFlags...)
}

var Flags []cli.Flag
