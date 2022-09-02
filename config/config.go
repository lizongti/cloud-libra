package config

import (
	"os"
	"strings"

	"github.com/tidwall/sjson"
	"github.com/urfave/cli"
)

type Config struct {
	Env      []byte
	CmdFlags []byte
	AppFlags []byte
	Stdin    []byte
	Dir      []byte
	Etcd     []byte
	Redis    []byte
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) initEnv() {
	var err error
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		c.Env, err = sjson.SetBytes(c.Env, pair[0], pair[1])
		if err != nil {
			panic(err)
		}
	}
}

func (c *Config) initCmdFlags(ctx *cli.Context) {
	var err error
	for _, name := range ctx.FlagNames() {
		if s := ctx.String(name); s == "" {
			continue
		}
		c.CmdFlags, err = sjson.SetBytes(c.CmdFlags, name, ctx.String(name))
		if err != nil {
			panic(err)
		}
	}
}

func (c *Config) initAppFlags(ctx *cli.Context) {
	var err error
	for _, name := range ctx.GlobalFlagNames() {
		if ctx.GlobalString(name) == "" {
			continue
		}
		c.AppFlags, err = sjson.SetBytes(c.AppFlags, name, ctx.GlobalString(name))
		if err != nil {
			panic(err)
		}
	}
}
