package sdk

import (
	"fmt"
	"strings"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/execx"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/gen"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/stringx"
	"gitlab.kaiqitech.com/nitro/nitro/v3/instrument/logger"
)

type Generator struct {
	*gen.Generator
}

func NewGenerator(dir string) *Generator {
	generator, _ := gen.NewGenerator(dir)
	return &Generator{
		Generator: generator,
	}
}

func (self *Generator) genSngSdkGo(serviceName string) error {
	apiFilePath := fmt.Sprintf("%v%v/logic/api.go", self.Dir(), serviceName)
	apimDir := fmt.Sprintf("%v%v/models/apim/", self.Dir(), serviceName)
	foldName := fmt.Sprintf("%v/", strings.ReplaceAll(serviceName, "-", "_"))
	packageName := strings.ReplaceAll(serviceName, "-", "")

	genItems := []*gen.Item{
		{
			Name: foldName,
			Type: gen.Dir,
		},
	}

	endpoints := getApiEndpoints(apiFilePath)
	for _, endpoint := range endpoints {
		handleName := strings.Split(endpoint.Name, ".")[0]
		funcName := strings.Split(endpoint.Name, ".")[1]
		genFileNamePrefix := stringx.From(handleName).ToSnake()
		genFilename := stringx.From(funcName).ToSnake()
		reqSourceFiles := []string{
			apimDir + genFilename + "_req.go",
			apimDir + genFileNamePrefix + "_" + genFilename + "_req.go",
		}
		rspSourceFiles := []string{
			apimDir + genFilename + "_rsp.go",
			apimDir + genFileNamePrefix + "_" + genFilename + "_rsp.go",
		}

		genItems = append(genItems, &gen.Item{
			Name:     foldName + genFilename + ".go",
			Type:     gen.File,
			Template: sngSdkTpl,
			Data: map[string]interface{}{
				"service":      serviceName,
				"endpointName": endpoint.Name,
				"pkg":          packageName,
				"import":       getImportContent(reqSourceFiles...) + "\n" + getImportContent(rspSourceFiles...),
				"req":          getReqOrRspContent(reqSourceFiles...),
				"rsp":          getReqOrRspContent(rspSourceFiles...),
				"handle":       handleName,
				"func":         funcName,
			},
		})
	}

	return self.GenItems(genItems)
}

func (self *Generator) GenAllSngServiceSdk() error {
	allServices := []string{
		// "ai-service",
		// "stage-service",
		// "match-service",
		// "casualgame-service",
		// "box-service",
		// "team-service",
		// "gametree-service",
		// "money-service",
		// "item-service",
		// "game-service",
		// "account-service",
		"login-service",
		// "user-service",
		// "quest-service",
		// "antiaddiction-service",
	}
	for _, service := range allServices {
		logger.Infof("gen sdk %s", service)
		// 拉取服务项目至本地
		info, err := execx.Run(fmt.Sprintf(
			`git clone https://gitlab.kaiqitech.com/k7game/server/services/%s`, service,
		), self.Dir())
		if err != nil {
			self.Error(err.Error())
			continue
		}
		self.Info(info)
		// 生成sdk
		if err = self.genSngSdkGo(service); err != nil {
			self.Error(err.Error())
			continue
		}

		// 删除项目
		_, err = execx.Run("rm -Rf "+service, self.Dir())
		if err != nil {
			self.Error(err.Error())
			continue
		}
	}
	// 格式化
	_, err := execx.Run("gofmt -l -w ./", self.Dir())
	if err != nil {
		self.Error(err.Error())
	}
	return nil
}
