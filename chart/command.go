package chart

import (
	"errors"
	"strings"

	"github.com/urfave/cli"
)

func CreateSngService(ctx *cli.Context) error {
	name := strings.TrimSpace(ctx.String("name"))
	dir := strings.TrimSpace(ctx.String("dir"))
	foldName := strings.TrimSpace(ctx.String("fold"))
	if len(name) == 0 {
		return errors.New("generate sng service chart error: no name")
	}
	gen, err := NewGenerator(name, dir, foldName)
	if err != nil {
		return err
	}
	return gen.GenSngServiceCharts()
}

func CreateSngGateway(ctx *cli.Context) error {
	name := strings.TrimSpace(ctx.String("name"))
	dir := strings.TrimSpace(ctx.String("dir"))
	foldName := strings.TrimSpace(ctx.String("fold"))
	if len(name) == 0 {
		return errors.New("generate sng gateway chart error: no name")
	}
	gen, err := NewGenerator(name, dir, foldName)
	if err != nil {
		return err
	}
	return gen.GenSngGatewayCharts()
}
