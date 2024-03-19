package sng

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/console"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/execx"
	"gitlab.kaiqitech.com/nitro/nitro/v3/util/pathx"
)

func CreateSngServiceProject(context *cli.Context) error {
	fileName := context.String("config")
	if len(fileName) == 0 {
		return errors.New("未指定配置文件")
	}

	gen, err := NewGenerator(fileName)
	if err != nil {
		return err
	}
	if err := gen.GenSngServiceAndTest(); err != nil {
		return err
	}
	return nil
}

func OpenProject(ctx *cli.Context) error {
	projectName := ctx.String("name")
	if len(projectName) == 0 {
		return errors.New("未指定项目名")
	}
	return openProject(projectName)
}

func openProject(projectName string) error {
	projectNames := strings.Split(projectName, ",")

	gitlabDir := os.Getenv("GITLAB")
	if len(gitlabDir) == 0 {
		return errors.New("未配置环境变量$GITLAB")
	}

	githubDir := os.Getenv("GITHUB")
	if len(gitlabDir) == 0 {
		return errors.New("未配置环境变量$GITHUB")
	}

	for _, projectName := range projectNames {
		projectDir := ""
		// 先查找gitlab项目
		filepath.WalkDir(gitlabDir, func(path string, dirInfo fs.DirEntry, err error) error {
			if dirInfo.IsDir() && strings.ToLower(dirInfo.Name()) == strings.ToLower(projectName) {
				if pathx.FileExists(filepath.Join(path, "go.mod")) ||
					pathx.FileExists(filepath.Join(path, "ChangeLog")) ||
					pathx.FileExists(filepath.Join(path, "README.md")) {
					projectDir = path
					return errors.New("找到项目")
				}
			}
			return nil
		})
		if len(projectDir) > 0 {
			execx.Run("code "+projectDir, "./")
			continue
		}
		// 再查找github项目
		filepath.WalkDir(githubDir, func(path string, dirInfo fs.DirEntry, err error) error {
			if dirInfo.IsDir() && strings.ToLower(dirInfo.Name()) == strings.ToLower(projectName) {
				if pathx.FileExists(filepath.Join(path, "go.mod")) ||
					pathx.FileExists(filepath.Join(path, "ChangeLog")) ||
					pathx.FileExists(filepath.Join(path, "README.md")) {
					projectDir = path
					return errors.New("找到项目")
				}
			}
			return nil
		})
		if len(projectDir) > 0 {
			execx.Run("code "+projectDir, "./")
			continue
		}
		console.NewColorConsole().Error("未找到项目: %s", projectName)
	}

	return nil
}
