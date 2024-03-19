package api

import (
	"strings"

	"github.com/urfave/cli"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/stringx"
)

func CreateSngService(ctx *cli.Context) error {
	token := strings.TrimSpace(ctx.String("token"))
	dir := strings.TrimSpace(ctx.String("dir"))
	serviceName := strings.TrimSpace(ctx.String("name"))
	projectIds := stringx.IntSlice(ctx.String("project"))
	apiIds := stringx.IntSlice(ctx.String("api"))
	gen, err := NewGenerator(serviceName, dir, token, projectIds, apiIds)
	if err != nil {
		return err
	}

	if err := gen.GenServiceApi(); err != nil {
		return err
	}

	return nil
}

func CreateSngServiceTest(ctx *cli.Context) error {
	token := strings.TrimSpace(ctx.String("token"))
	dir := strings.TrimSpace(ctx.String("dir"))
	serviceName := strings.TrimSpace(ctx.String("name"))
	projectIds := stringx.IntSlice(ctx.String("project"))
	apiIds := stringx.IntSlice(ctx.String("api"))
	gen, err := NewGenerator(serviceName, dir, token, projectIds, apiIds)
	if err != nil {
		return err
	}

	if err := gen.GenTestApi(); err != nil {
		return err
	}

	return nil
}

func Cpp(ctx *cli.Context) error {
	token := strings.TrimSpace(ctx.String("token"))
	dir := strings.TrimSpace(ctx.String("dir"))
	serviceName := strings.TrimSpace(ctx.String("service"))
	projectIds := stringx.IntSlice(ctx.String("project"))
	apiIds := stringx.IntSlice(ctx.String("api"))
	gen, err := NewGenerator(serviceName, dir, token, projectIds, apiIds)
	if err != nil {
		return err
	}

	if err := gen.GenCpp(); err != nil {
		return err
	}

	return nil
}
