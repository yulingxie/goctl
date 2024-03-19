package deploy

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

var (
	InvalidEnvirement    = errors.New("请指定正确的环境 dev/test/qa")
	HostUserNotSpecified = errors.New("主机用户未指定")
	ChiperNotSpecified   = errors.New("密码未指定")
	PkgNotSpecified      = errors.New("安装包未指定")
	InvalidServer        = errors.New("请指定正确的服务器类型")
)

type config struct {
	Env          string
	User         string
	Password     string
	IdentityFile string
}

func (c *config) Validate() error {
	if len(c.Env) == 0 {
		return InvalidEnvirement
	}
	switch c.Env {
	case "dev", "qa", "test":
	default:
		return InvalidEnvirement
	}
	if len(c.User) == 0 {
		return HostUserNotSpecified
	}
	if len(c.Password) == 0 && len(c.IdentityFile) == 0 {
		return ChiperNotSpecified
	}
	return nil
}

func NewConfigFromContext(ctx *cli.Context) (*config, error) {
	cfg := &config{}
	cfg.Env = strings.ToLower(strings.TrimSpace(ctx.String("env")))
	cfg.User = strings.TrimSpace(ctx.String("user"))
	cfg.Password = strings.TrimSpace(ctx.String("password"))
	cfg.IdentityFile = strings.TrimSpace(ctx.String("identity_file"))
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

type generatorConfig struct {
	Name        string
	Env         string
	RemoteFile  string
	LocalFile   string
	URL         string
	BuildNumber string
	Verbose     bool
}

func NewGeneratorConfigFromContext(ctx *cli.Context) (*generatorConfig, error) {
	c := &generatorConfig{}
	c.Env = strings.ToLower(strings.TrimSpace(ctx.String("env")))
	switch c.Env {
	case "test", "qa", "dev":
	default:
		return nil, errors.New("")
	}
	server := strings.TrimSpace(ctx.String("server"))
	switch server {
	case "conn", "game", "sync":
		c.Name = strings.Title(server) + "Server"
	case "db":
		c.Name = strings.ToUpper(server) + "Server"
	case "center":
		c.Name = "GameCenter"
	default:
		return nil, InvalidServer
	}
	c.LocalFile = strings.TrimSpace(ctx.String("file"))
	c.URL = strings.TrimSpace(ctx.String("url"))
	c.Verbose = ctx.Bool("verbose")
	if len(c.LocalFile) == 0 && len(c.URL) == 0 {
		return nil, PkgNotSpecified
	}
	if len(c.LocalFile) != 0 {
		c.RemoteFile = fmt.Sprintf("/var/tmp/%v.tar.gz", c.Name)
	}
	return c, nil
}
