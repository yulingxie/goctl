package upgrade

import (
	"github.com/urfave/cli"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/console"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/execx"
)

func Upgrade(ctx *cli.Context) error {
	branch := ctx.String("branch")
	if len(branch) == 0 {
		// 默认使用master分支
		branch = "master"
	}

	_, err := execx.Run("go get -u gitlab.kaiqitech.com/k7game/server/tools/goctl@"+branch, "")
	if err != nil {
		return err
	}
	console.NewColorConsole().Info("goctl %s 更新成功", branch)
	return nil
}
