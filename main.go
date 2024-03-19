package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/api"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/chart"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/confluence"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/deploy"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/git"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/jenkins"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/mongo"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/redis"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/sql/command"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/password"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/project/sng"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/protok/monitor"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/upgrade"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/version"
)

var (
	commands = []cli.Command{
		{
			Name:  "upgrade",
			Usage: "更新至插件最新版本",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "branch,b",
					Usage: "代码分支",
				},
			},
			Action: upgrade.Upgrade,
		},
		{
			Name:  "crypt",
			Usage: "通过指定密钥对原文进行加密",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "password",
					Usage: "加密密钥",
				},
				cli.StringFlag{
					Name:  "content",
					Usage: "加密内容",
				},
			},
			Action: password.Crypt,
		},
		{
			Name:  "decrypt",
			Usage: "通过指定密钥对原文进行解密",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "password",
					Usage: "解密密钥",
				},
				cli.StringFlag{
					Name:  "content",
					Usage: "解密内容",
				},
			},
			Action: password.Decrypt,
		},
		{
			Name:  "chart",
			Usage: "生成项目chart文件",
			Subcommands: []cli.Command{
				{
					Name:  "create-sng-service",
					Usage: `生成sng service的chart文件`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name,n",
							Usage: "项目名称",
						},
						cli.StringFlag{
							Name:  "dir,d",
							Usage: "生成目录",
						},
						cli.StringFlag{
							Name:  "fold,f",
							Usage: "生成的文件夹名",
						},
					},
					Action: chart.CreateSngService,
				},
				{
					Name:  "create-sng-gw",
					Usage: `生成sng gateway的chart文件`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name,n",
							Usage: "项目名称",
						},
						cli.StringFlag{
							Name:  "dir,d",
							Usage: "生成目录",
						},
						cli.StringFlag{
							Name:  "fold,f",
							Usage: "生成的文件夹名",
						},
					},
					Action: chart.CreateSngGateway,
				},
			},
		},
		{
			Name:  "api",
			Usage: "根据yapi文档生成项目模板",
			Subcommands: []cli.Command{
				{
					Name:  "create-sng-service",
					Usage: `生成sng项目的api项目模板`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name,n",
							Usage: "项目名称",
						},
						cli.StringFlag{
							Name:  "token,t",
							Usage: "yapi项目token",
						},
						cli.StringFlag{
							Name:  "project",
							Usage: "yapi项目id,可指定多个，其间用`,`分隔",
						},
						cli.StringFlag{
							Name:  "api",
							Usage: "yapi接口id,可指定多个，其间用`,`分隔",
						},
						cli.StringFlag{
							Name:  "dir,d",
							Usage: "生成目录",
						},
					},
					Action: api.CreateSngService,
				},
				{
					Name:  "create-sng-service-test",
					Usage: `生成sng service test的api项目模板`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name,n",
							Usage: "项目名称",
						},
						cli.StringFlag{
							Name:  "token,t",
							Usage: "yapi项目token",
						},
						cli.StringFlag{
							Name:  "project",
							Usage: "yapi项目id",
						},
						cli.StringFlag{
							Name:  "api",
							Usage: "yapi单个接口id",
						},
						cli.StringFlag{
							Name:  "dir,d",
							Usage: "生成目录",
						},
					},
					Action: api.CreateSngServiceTest,
				},
			},
		},
		{
			Name:  "git",
			Usage: "配置gitlab项目",
			Subcommands: []cli.Command{
				{
					Name:  "create-sng-service",
					Usage: `生成配置sng service的gitlab项目配置`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "token,t",
							Usage: "gitlab鉴权token",
						},
						cli.StringFlag{
							Name:  "dec,d",
							Usage: "服务简述",
						},
						cli.StringFlag{
							Name:  "name,n",
							Usage: "服务名",
						},
					},
					Action: git.CreateSngService,
				},
				{
					Name:  "create-sng-service-test",
					Usage: `生成配置sng service test的gitlab项目配置`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "token,t",
							Usage: "gitlab鉴权token",
						},
						cli.StringFlag{
							Name:  "dec,d",
							Usage: "项目简述",
						},
						cli.StringFlag{
							Name:  "name,n",
							Usage: "服务名",
						},
					},
					Action: git.CreateSngServiceTest,
				},
				{
					Name:  "get",
					Usage: `查询项目信息`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "token,t",
							Usage: "gitlab鉴权token",
						},
						cli.StringFlag{
							Name:  "name,n",
							Usage: "项目名",
						},
					},
					Action: git.Get,
				},
				{
					Name:  "open",
					Usage: `使用浏览器打开项目网址`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "token,t",
							Usage: "gitlab鉴权token",
						},
						cli.StringFlag{
							Name:  "name,n",
							Usage: "项目名",
						},
					},
					Action: git.Open,
				},
				{
					Name:  "delete",
					Usage: `删除项目`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "token,t",
							Usage: "gitlab鉴权token",
						},
						cli.StringFlag{
							Name:  "id",
							Usage: "项目id,可以指定多个，用`,`隔开",
						},
					},
					Action: git.Delete,
				},
				{
					Name:  "clone",
					Usage: `克隆项目至本地`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "token,t",
							Usage: "gitlab鉴权token",
						},
						cli.StringSliceFlag{
							Name:  "name,n",
							Usage: "项目名,可以指定多个，用`,`隔开",
						},
						cli.IntSliceFlag{
							Name:  "id",
							Usage: "项目id,可以指定多个，用`,`隔开",
						},
					},
					Action: git.Clone,
				},
			},
		},
		{
			Name:  "jenkins",
			Usage: "配置项目jenkins",
			Subcommands: []cli.Command{
				{
					Name:  "create-sng-service",
					Usage: `创建sng service的jenkins job`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "user,u",
							Usage: "用户名",
						},
						cli.StringFlag{
							Name:  "pass,p",
							Usage: "密码",
						},
						cli.StringFlag{
							Name:  "name,n",
							Usage: "服务名,可指定多个服务，用`,`间隔",
						},
					},
					Action: jenkins.CreateSngService,
				},
				{
					Name:  "update-sng-service",
					Usage: `使用插件中的最新模版更新sng service的jenkins配置`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "user,u",
							Usage: "用户名",
						},
						cli.StringFlag{
							Name:  "pass,p",
							Usage: "密码",
						},
						cli.StringFlag{
							Name:  "name,n",
							Usage: "服务名,可指定多个服务，用`,`间隔",
						},
					},
					Action: jenkins.UpdateSngService,
				},
				{
					Name:  "delete-sng-service",
					Usage: `删除sng service的jenkins配置`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "user,u",
							Usage: "用户名",
						},
						cli.StringFlag{
							Name:  "pass,p",
							Usage: "密码",
						},
						cli.StringFlag{
							Name:  "name,n",
							Usage: "服务名,可指定多个服务，用`,`间隔",
						},
					},
					Action: jenkins.DeleteSngService,
				},
			},
		},
		{
			Name:  "model",
			Usage: "generate model code",
			Subcommands: []cli.Command{
				{
					Name:  "sql",
					Usage: `生成sql model模版代码`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "url, u",
							Usage: `数据库url`,
						},
						cli.StringFlag{
							Name:  "node, n",
							Usage: `数据库节点名，默认为common`,
						},
						cli.StringFlag{
							Name:  "db",
							Usage: `数据库名`,
						},
						cli.StringFlag{
							Name:  "table, t",
							Usage: `表名`,
						},
						cli.StringFlag{
							Name:  "dir, d",
							Usage: "生成目录",
						},
						cli.BoolFlag{
							Name:  "cache, c",
							Usage: "是否是缓存代码",
						},
					},
					Action: command.CreateSqlModel,
				},
				{
					Name:  "mongo",
					Usage: `生成mongo数据库model代码`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir, d",
							Usage: "生成目录",
						},
						cli.StringFlag{
							Name:  "node, n",
							Usage: "数据库节点名，默认为common",
						},
						cli.StringFlag{
							Name:  "db",
							Usage: "数据库表名",
						},
						cli.StringFlag{
							Name:  "coll, c",
							Usage: "集合名，可指定多个，其中用`,`隔开",
						},
					},
					Action: mongo.Model,
				},
				{
					Name:  "redis",
					Usage: `生成redis数据库model代码`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "node, n",
							Usage: "数据库节点名，默认为common",
						},
						cli.StringSliceFlag{
							Name:  "table, t",
							Usage: "表名",
						},
						cli.StringFlag{
							Name:  "dir, d",
							Usage: "生成目录",
						},
					},
					Action: redis.Model,
				},
			},
		},
		{
			Name:  "project",
			Usage: "生成项目模板",
			Subcommands: []cli.Command{
				{
					Name:  "create-sng-service",
					Usage: `生成sng service项目`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "config, c",
							Usage: "配置文件名",
						},
					},
					Action: sng.CreateSngServiceProject,
				},
				{
					Name:  "open",
					Usage: `打开本地项目`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "项目名",
						},
					},
					Action: sng.OpenProject,
				},
			},
		},
		{
			Name:  "confluence",
			Usage: "更新知识库文档",
			Subcommands: []cli.Command{
				{
					Name:  "update-sng-errors",
					Usage: `更新sng服务错误码文档`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "user, u",
							Usage: "账号",
						},
						cli.StringFlag{
							Name:  "pass, p",
							Usage: "密码",
						},
						cli.StringFlag{
							Name:  "file, f",
							Usage: "go源码文件",
						},
					},
					Action: confluence.UpdateSngErrors,
				},
				{
					Name:  "update-cac-errors",
					Usage: `更新cac服务错误码文档`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "user, u",
							Usage: "账号",
						},
						cli.StringFlag{
							Name:  "pass, p",
							Usage: "密码",
						},
						cli.StringFlag{
							Name:  "file, f",
							Usage: "go源码文件",
						},
					},
					Action: confluence.UpdateCacErrors,
				},
				{
					Name:  "update-changelog",
					Usage: `更新CHANGELOG文档`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "user, u",
							Usage: "账号",
						},
						cli.StringFlag{
							Name:  "pass, p",
							Usage: "密码",
						},
						cli.StringFlag{
							Name:  "file, f",
							Usage: "go源码文件",
						},
						cli.StringFlag{
							Name:  "service, s",
							Usage: "go服务名称",
						},
						cli.StringFlag{
							Name:  "parent_page, pp",
							Usage: "指定confluence父页面",
						},
						cli.StringFlag{
							Name:  "link, l",
							Usage: "下载地址",
						},
						cli.StringFlag{
							Name:  "confluence_space, cs",
							Usage: "confluence空间标识",
						},
					},
					Action: confluence.UpdateChangeLog,
				},
				{
					Name:  "update-sng-dev-build",
					Usage: `更新sng dev构建数据`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "user, u",
							Usage: "账号",
						},
						cli.StringFlag{
							Name:  "pass, p",
							Usage: "密码",
						},
						cli.StringFlag{
							Name:  "name, n",
							Usage: "服务名",
						},
						cli.StringFlag{
							Name:  "version, v",
							Usage: "版本",
						},
						cli.StringFlag{
							Name:  "coverage, c",
							Usage: "覆盖率",
						},
						cli.StringFlag{
							Name:  "report, r",
							Usage: "报告地址",
						},
					},
					Action: confluence.UpdateSngDevBuild,
				},
				{
					Name:  "update-sng-test-build",
					Usage: `更新sng test构建数据`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "user, u",
							Usage: "账号",
						},
						cli.StringFlag{
							Name:  "pass, p",
							Usage: "密码",
						},
						cli.StringFlag{
							Name:  "name, n",
							Usage: "服务名",
						},
						cli.StringFlag{
							Name:  "api, a",
							Usage: "是否有api测试",
						},
						cli.StringFlag{
							Name:  "conn, c",
							Usage: "是否有conn测试",
						},
						cli.StringFlag{
							Name:  "game, g",
							Usage: "是否有game测试",
						},
						cli.StringFlag{
							Name:  "benchmark, b",
							Usage: "是否有benchmark测试",
						},
						cli.StringFlag{
							Name:  "report, r",
							Usage: "报告地址",
						},
					},
					Action: confluence.UpdateSngTestBuild,
				},
			},
		},
		{
			Name:  "deploy",
			Usage: "自动部署服务端",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "server, s",
					Usage: "需要部署的服务端 conn/db/game/center/sync",
				},
				cli.StringFlag{
					Name:  "env, e",
					Usage: "目标环境 dev/test/qa",
				},
				cli.StringFlag{
					Name:  "url, u",
					Usage: "远程包路径",
				},
				cli.StringFlag{
					Name:  "file, f",
					Usage: "本地tar.gz包路径",
				},
				cli.StringFlag{
					Name:  "user",
					Usage: "服务器用户名",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "服务器密码",
				},
				cli.StringFlag{
					Name:  "identity_file, i",
					Usage: "服务器密钥",
				},
				cli.StringFlag{
					Name:  "verbose, v",
					Usage: "服务器密钥",
				},
			},
			Action: deploy.Deploy,
		},
		{
			Name:  "protok",
			Usage: "流量监控解析",
			Subcommands: []cli.Command{
				{
					Name:  "monitor",
					Usage: "协议解析",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "server, s",
							Usage:    "带端口号的服务端地址, 例: 1.2.3.4:9000",
							Required: true,
						},
						cli.StringFlag{
							Name:     "protocol, p",
							Usage:    "协议类型, 目前支持: 连接服(conn)/游戏服(game)",
							Required: true,
						},
						cli.BoolFlag{
							Name:  "heartbeat",
							Usage: "是否解析心跳包",
						},
						cli.StringSliceFlag{
							Name:  "filter,f",
							Usage: "需要解析的数据包, 不指定则解析全部。连接服(conn)/游戏服(game)用法: --filter 1001,10 --filter 1002,1",
							Value: &cli.StringSlice{},
						},
					},
					Action: monitor.TrafficMonitor,
				},
			},
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Usage = "a cli tool to generate code"
	app.Version = fmt.Sprintf("%s %s/%s", version.Version, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands
	if err := app.Run(os.Args); err != nil {
		fmt.Println("error:", err)
	}
}
