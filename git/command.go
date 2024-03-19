package git

import (
	"errors"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

func CreateSngService(ctx *cli.Context) error {
	token := strings.TrimSpace(ctx.String("token"))
	name := strings.TrimSpace(ctx.String("name"))
	dec := strings.TrimSpace(ctx.String("dec"))
	if len(name) == 0 {
		return errors.New("未指定git项目名称")
	}

	gitClient := NewGit(token)
	_, err := gitClient.CreateSngService(name, dec)
	return err
}

func CreateSngServiceTest(ctx *cli.Context) error {
	token := strings.TrimSpace(ctx.String("token"))
	name := strings.TrimSpace(ctx.String("name"))
	dec := strings.TrimSpace(ctx.String("dec"))
	if len(name) == 0 {
		return errors.New("未指定git项目名称")
	}

	gitClient := NewGit(token)
	_, err := gitClient.CreateSngServiceTest(name, dec)
	return err
}

func Get(ctx *cli.Context) error {
	token := strings.TrimSpace(ctx.String("token"))
	name := strings.TrimSpace(ctx.String("name"))
	return NewGit(token).GetProject(name)
}

func Open(ctx *cli.Context) error {
	token := strings.TrimSpace(ctx.String("token"))
	name := strings.TrimSpace(ctx.String("name"))
	return NewGit(token).OpenProject(name)
}

func Delete(ctx *cli.Context) error {
	token := strings.TrimSpace(ctx.String("token"))
	id := strings.TrimSpace(ctx.String("id"))
	if len(id) == 0 {
		return errors.New("未指定项目id")
	}
	ids := strings.Split(id, ",")
	projectIds := []int{}
	for _, id := range ids {
		projectId, err := strconv.Atoi(id)
		if err != nil {
			return nil
		}
		projectIds = append(projectIds, projectId)
	}

	return NewGit(token).DeleteProject(projectIds...)
}

func Clone(ctx *cli.Context) error {
	token := strings.TrimSpace(ctx.String("token"))
	names := ctx.StringSlice("name")
	ids := ctx.IntSlice("id")
	NewGit(token).CloneProject(names, ids)
	return nil
}
