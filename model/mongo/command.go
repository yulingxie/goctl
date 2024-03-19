package mongo

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

	dbName := strings.TrimSpace(ctx.String("db"))
	if len(dbName) == 0 {
		return errors.New("have no db name")
	}

	collName := strings.TrimSpace(ctx.String("coll"))
	if len(collName) == 0 {
		return errors.New("have no coll name")
	}

	dir := strings.TrimSpace(ctx.String("dir"))
	gen, err := NewGenerator(dir, nodeName, dbName, collName)
	if err != nil {
		return err
	}
	return gen.GenModel()
}
