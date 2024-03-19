package confluence

import (
	"strings"

	"github.com/urfave/cli"
)

func UpdateSngErrors(ctx *cli.Context) error {
	file := strings.TrimSpace(ctx.String("file"))
	user := strings.TrimSpace(ctx.String("user"))
	pass := strings.TrimSpace(ctx.String("pass"))
	return NewConfluence(user, pass).UpdateSngErrors(file)
}

func UpdateCacErrors(ctx *cli.Context) error {
	file := strings.TrimSpace(ctx.String("file"))
	user := strings.TrimSpace(ctx.String("user"))
	pass := strings.TrimSpace(ctx.String("pass"))
	return NewConfluence(user, pass).UpdateCacErrors(file)
}

func UpdateChangeLog(ctx *cli.Context) error {
	file := strings.TrimSpace(ctx.String("file"))
	user := strings.TrimSpace(ctx.String("user"))
	pass := strings.TrimSpace(ctx.String("pass"))
	serviceName := strings.TrimSpace(ctx.String("service"))
	parentPage := strings.TrimSpace(ctx.String("parent_page"))
	Link := strings.TrimSpace(ctx.String("link"))
	confluenceSpace := strings.TrimSpace(ctx.String("confluence_space"))
	return NewConfluence(user, pass).UpdateChangeLogMarkdown(file, serviceName, parentPage, Link, confluenceSpace)
}

func UpdateSngDevBuild(ctx *cli.Context) error {
	user := strings.TrimSpace(ctx.String("user"))
	pass := strings.TrimSpace(ctx.String("pass"))
	name := strings.TrimSpace(ctx.String("name"))
	version := strings.TrimSpace(ctx.String("version"))
	coverage := strings.TrimSpace(ctx.String("coverage"))
	reportUrl := strings.TrimSpace(ctx.String("report"))
	return NewConfluence(user, pass).UpdateSngDevBuild(name, version, coverage, reportUrl)
}

func UpdateSngTestBuild(ctx *cli.Context) error {
	user := strings.TrimSpace(ctx.String("user"))
	pass := strings.TrimSpace(ctx.String("pass"))
	name := strings.TrimSpace(ctx.String("name"))
	api := strings.TrimSpace(ctx.String("api"))
	conn := strings.TrimSpace(ctx.String("conn"))
	game := strings.TrimSpace(ctx.String("game"))
	benchmark := strings.TrimSpace(ctx.String("benchmark"))
	reportUrl := strings.TrimSpace(ctx.String("report"))
	return NewConfluence(user, pass).UpdateSngTestBuild(name, api, conn, game, benchmark, reportUrl)
}
