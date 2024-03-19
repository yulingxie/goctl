package sng

import (
	"fmt"
	"time"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/execx"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/gen"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/stringx"
)

const (
	dockerfileName = "Dockerfile"
)

type Generator struct {
	*gen.Generator
	cfg *SngServiceProjectConfig
}

func NewGenerator(fileName string) (*Generator, error) {
	generator, err := gen.NewGenerator("./")
	if err != nil {
		return nil, err
	}
	cfg, err := LoadSngServiceProjectConfigFromFile(fileName)
	if err != nil {
		return nil, err
	}
	return &Generator{cfg: cfg, Generator: generator}, nil
}

func (self *Generator) GenSngServiceAndTest() error {
	if err := self.genSngService(); err != nil {
		return err
	}
	// 等待5s，保证gitlab上传完毕
	<-time.After(time.Second * 5)
	return self.genSngServiceTest()
}

func (self *Generator) genSngService() error {
	var (
		info string
		err  error
	)

	serviceDir := self.Dir() + self.cfg.ServiceDir + self.cfg.ServiceName + "-service/"

	if len(self.cfg.GitlabToken) > 0 {
		//配置gitlab项目
		info, err = execx.Run(fmt.Sprintf(
			`goctl git create-sng-service -t %s -n %s-service`,
			self.cfg.GitlabToken, self.cfg.ServiceName,
		), "./")
		if err != nil {
			return err
		}
		self.Info(info)

		// 拉取服务项目至本地
		info, err = execx.Run(fmt.Sprintf(
			`git clone https://gitlab.kaiqitech.com/k7game/server/services/%s-service %s`,
			self.cfg.ServiceName, serviceDir,
		), "./")
		if err != nil {
			return err
		}
		self.Info(info)
	}

	// 配置jenkins
	if len(self.cfg.Jenkins.User) > 0 {
		info, err = execx.Run(fmt.Sprintf(
			`goctl jenkins create-sng-service -u %s -p %s -n %s-service`,
			self.cfg.Jenkins.User, self.cfg.Jenkins.Pass, self.cfg.ServiceName,
		), "./")
		if err != nil {
			return err
		}
		self.Info(info)
	}

	// 创建项目基础模版文件
	data := map[string]interface{}{
		"serviceName": self.cfg.ServiceName + "-service",
		"upperName":   stringx.From(self.cfg.ServiceName).Title(),
	}
	if err := self.GenItems([]*gen.Item{
		{
			Name:     serviceDir + ".vscode",
			Type:     gen.Dir,
			Template: "",
		},
		{
			Name:     serviceDir + ".vscode/launch.json",
			Type:     gen.File,
			Template: vscodeTemplate,
		},
		{
			Name:     serviceDir + ".gitignore",
			Type:     gen.File,
			Template: gitignorTemplate,
			Data:     nil,
		},
		{
			Name:     serviceDir + "CHANGELOG.md",
			Type:     gen.File,
			Template: changeLogTemplate,
			Data:     nil,
		},
		{
			Name:     serviceDir + "Dockerfile",
			Type:     gen.File,
			Template: dockerfileTemplate,
			Data:     data,
		},
		{
			Name:     serviceDir + "docker-entrypoint.sh",
			Type:     gen.File,
			Template: dockerEntrypointShTpl,
		},
		{
			Name:     serviceDir + "go.mod",
			Type:     gen.File,
			Template: goModTemplate,
			Data:     data,
		},
		{
			Name:     serviceDir + "main.go",
			Type:     gen.File,
			Template: mainTemplate,
			Data:     data,
		},
	}); err != nil {
		return err
	}

	// service目录下生成chart
	info, err = execx.Run(fmt.Sprintf(
		`goctl chart create-sng-service -n %s-service -f chart -d %s`,
		self.cfg.ServiceName, serviceDir,
	), "./")
	if err != nil {
		return err
	}
	self.Info(info)

	// service目录下生成api相关文件
	info, err = execx.Run(fmt.Sprintf(
		`goctl api create-sng-service -t %s -project %d -d %s -n %s-service`,
		self.cfg.Yapi.Token, self.cfg.Yapi.Id, serviceDir, self.cfg.ServiceName,
	), "./")
	if err != nil {
		return err
	}
	self.Info(info)

	// 生成sqlm
	sqlmDir := serviceDir + "models/sqlm/"
	for _, sqlmConfig := range self.cfg.Sql {
		for _, table := range sqlmConfig.Table {
			info, err = execx.Run(fmt.Sprintf(
				`goctl model sql -node %s -db %s -table %s -dir %s `,
				sqlmConfig.Node, sqlmConfig.Db, table, sqlmDir,
			), "./")
			if err != nil {
				return err
			}
			self.Info(info)
		}
	}

	// 生成redism
	redismDir := serviceDir + "models/redism/"
	for _, redisConfig := range self.cfg.Redis {
		for _, table := range redisConfig.Table {
			info, err = execx.Run(fmt.Sprintf(
				`goctl model redis -node %s -table %s -dir %s `,
				redisConfig.Node, table, redismDir,
			), "./")
			if err != nil {
				return err
			}
			self.Info(info)
		}
	}

	// 生成mongom
	mongomDir := serviceDir + "models/mongom/"
	for _, mongomConfig := range self.cfg.Mongo {
		for _, coll := range mongomConfig.Coll {
			info, err = execx.Run(fmt.Sprintf(
				`goctl model mongo -node %s -db %s -coll %s -dir %s `,
				mongomConfig.Node, mongomConfig.Db, coll, mongomDir,
			), "./")
			if err != nil {
				return err
			}
			self.Info(info)
		}
	}

	//运行go mod命令以拉取最新引用库
	info, err = execx.Run("go mod tidy -go=1.16 && go mod tidy -go=1.17", serviceDir)
	if err != nil {
		return err
	}
	self.Info(info)

	// 使用gofmt格式化代码
	info, err = execx.Run("gofmt -l -w .", serviceDir)
	if err != nil {
		return err
	}
	self.Info(info)

	//提交并推送至gitlab远程仓库
	if len(self.cfg.GitlabToken) > 0 {
		info, err = execx.Run(`
		rm ./README.MD
		git add -A
		git commit -a -m "初始版本"
		git push origin dev
	`, serviceDir)
		if err != nil {
			return err
		}
		self.Info(info)
	}

	self.Info("服务项目创建完毕")
	return nil
}

func (self *Generator) genSngServiceTest() error {
	var (
		info string
		err  error
	)

	serviceTestDir := self.Dir() + self.cfg.ServiceTestDir + self.cfg.ServiceName + "-test/"
	if len(self.cfg.GitlabToken) > 0 {
		//配置gitlab项目
		info, err = execx.Run(fmt.Sprintf(
			`goctl git create-sng-service-test -t %s -n %s-test`,
			self.cfg.GitlabToken, self.cfg.ServiceName,
		), "./")
		if err != nil {
			return err
		}
		self.Info(info)

		// 拉取服务项目至本地
		info, err = execx.Run(fmt.Sprintf(
			`git clone https://gitlab.kaiqitech.com/k7game/server/test/%s-test %s`,
			self.cfg.ServiceName, serviceTestDir,
		), "./")
		if err != nil {
			return err
		}
		self.Info(info)
	}

	// 创建项目基础模版文件
	data := map[string]interface{}{
		"serviceName": self.cfg.ServiceName + "-service",
		"testName":    self.cfg.ServiceName + "-test",
		"upperName":   stringx.From(self.cfg.ServiceName).Title(),
	}
	if err := self.GenItems([]*gen.Item{
		{
			Name:     serviceTestDir + ".gitignore",
			Type:     gen.File,
			Template: gitignorTemplate,
			Data:     nil,
		},
		{
			Name:     serviceTestDir + "go.mod",
			Type:     gen.File,
			Template: sngApiTestGoMod,
			Data:     data,
		},
	}); err != nil {
		return err
	}

	// test目录下生成api测试文件
	info, err = execx.Run(fmt.Sprintf(
		`goctl api create-sng-service-test -t %s -project %d -d %s -n %s-service`,
		self.cfg.Yapi.Token, self.cfg.Yapi.Id, serviceTestDir, self.cfg.ServiceName,
	), "./")
	if err != nil {
		return err
	}
	self.Info(info)

	//运行go mod命令以拉取最新引用库
	info, err = execx.Run("go mod tidy -go=1.16 && go mod tidy -go=1.17", serviceTestDir)
	if err != nil {
		return err
	}
	self.Info(info)

	// 使用gofmt格式化代码
	info, err = execx.Run("gofmt -l -w .", serviceTestDir)
	if err != nil {
		return err
	}
	self.Info(info)

	//提交并推送至gitlab远程仓库
	if len(self.cfg.GitlabToken) > 0 {
		info, err = execx.Run(`
		rm ./README.MD
		git add -A
		git commit -a -m "初始版本"
		git push origin test
	`, serviceTestDir)
		if err != nil {
			return err
		}
		self.Info(info)
	}

	self.Info("测试项目创建完毕")
	return nil
}
