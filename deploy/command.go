package deploy

import (
	"github.com/urfave/cli"
)

func Deploy(ctx *cli.Context) error {
	gcfg, err := NewGeneratorConfigFromContext(ctx)
	if err != nil {
		return err
	}
	cfg, err := NewConfigFromContext(ctx)
	if err != nil {
		return err
	}
	if dpl, err := NewDeployer(cfg, gcfg); err != nil {
		return err
	} else {
		return dpl.Deploy()
	}
}
