package redis

import (
	"errors"
	"strings"

	"github.com/urfave/cli"
)

func Model(ctx *cli.Context) error {
	nodeName := strings.TrimSpace(ctx.String("node"))
	if len(nodeName) == 0 {
		nodeName = "common"
	}

	tableName := strings.TrimSpace(ctx.String("table"))
	if len(tableName) == 0 {
		return errors.New("未指定table name")
	}

	dir := strings.TrimSpace(ctx.String("dir"))

	gen, err := NewGenerator(dir, nodeName, tableName)
	if err != nil {
		return err
	}

	return gen.GenModel()
}
