package jenkins

import (
	"errors"
	"strings"

	"github.com/urfave/cli"
)

func CreateSngService(ctx *cli.Context) error {
	user := strings.TrimSpace(ctx.String("user"))
	pass := strings.TrimSpace(ctx.String("pass"))
	name := strings.TrimSpace(ctx.String("name"))

	if len(name) == 0 {
		return errors.New("未指定服务名")
	}
	names := strings.Split(name, ",")

	return NewJenkins(user, pass).CreateSngService(names...)
}

func UpdateSngService(ctx *cli.Context) error {
	user := strings.TrimSpace(ctx.String("user"))
	pass := strings.TrimSpace(ctx.String("pass"))
	name := strings.TrimSpace(ctx.String("name"))

	if len(name) == 0 {
		return errors.New("未指定服务名")
	}
	names := strings.Split(name, ",")

	return NewJenkins(user, pass).UpdateSngService(names...)
}

func DeleteSngService(ctx *cli.Context) error {
	user := strings.TrimSpace(ctx.String("user"))
	pass := strings.TrimSpace(ctx.String("pass"))
	name := strings.TrimSpace(ctx.String("name"))

	if len(name) == 0 {
		return errors.New("未指定服务名")
	}
	names := strings.Split(name, ",")

	return NewJenkins(user, pass).DeleteSngService(names...)
}
