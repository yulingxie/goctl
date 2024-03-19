package api

import (
	"fmt"
	"strings"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/execx"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/gen"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/stringx"
)

type Generator struct {
	serviceName string
	*gen.Generator
	apiInfos []*ApiInfo
}

func NewGenerator(serviceName, dir, token string, projectIds, apiIds []int) (*Generator, error) {
	yapi := NewYapi()
	generator, _ := gen.NewGenerator(dir)

	// 获取项目的所有api id
	for _, projectId := range projectIds {
		ids, err := yapi.GetProjectInfo(token, projectId)
		if err != nil {
			generator.Error(err.Error())
			continue
		}
		apiIds = append(apiIds, ids...)
	}

	// 获取api的具体信息
	apiInfos := []*ApiInfo{}
	for _, id := range apiIds {
		apiInfo, err := yapi.GetApiInfo(token, id)
		if err != nil {
			generator.Error(err.Error())
			continue
		}
		apiInfos = append(apiInfos, apiInfo.ToApiInfo())
	}

	return &Generator{Generator: generator, apiInfos: apiInfos, serviceName: serviceName}, nil
}

func (self *Generator) GenServiceApi() error {
	data := map[string]interface{}{
		"dot":      "`",
		"apiInfos": self.apiInfos,
	}
	if err := self.GenItems([]*gen.Item{
		{
			Name: "logic",
			Type: gen.Dir,
		},
		{
			Name: "models/apim",
			Type: gen.Dir,
		},
		{
			Name:     "logic/api.go",
			Type:     gen.File,
			Template: apiTpl,
			Data:     data,
		},
	}); err != nil {
		return err
	}

	dir := "models/apim/"
	for _, apiInfo := range self.apiInfos {
		data := map[string]interface{}{
			"dot":     "`",
			"apiInfo": apiInfo,
			"reqBody": GetGoCode(apiInfo.ReqBody.Properties),
			"rspBody": GetGoCode(apiInfo.RspBody.Properties.Data.Properties),
		}
		if err := self.GenItems([]*gen.Item{
			{
				Name:     dir + apiInfo.FileName + "_req.go",
				Type:     gen.File,
				Template: reqTpl,
				Data:     data,
			},
		}); err != nil {
			return err
		}

		if len(apiInfo.RspBody.Properties.Data.Properties) > 0 {
			if err := self.GenItems([]*gen.Item{
				{
					Name:     dir + apiInfo.FileName + "_rsp.go",
					Type:     gen.File,
					Template: rspTpl,
					Data:     data,
				},
			}); err != nil {
				return err
			}
		}
	}

	// 生成结构体文件
	structInfos := map[string][]string{}
	for _, apiInfo := range self.apiInfos {
		if structInfos[apiInfo.StructName] == nil {
			structInfos[apiInfo.StructName] = []string{}
		}
		structInfos[apiInfo.StructName] = append(structInfos[apiInfo.StructName], apiInfo.FuncName)
	}
	for structName, funcNames := range structInfos {
		if err := self.GenItems([]*gen.Item{
			{
				Name:     "logic/" + stringx.From(structName).ToSnake() + ".go",
				Type:     gen.File,
				Template: structTpl,
				Data: map[string]interface{}{
					"serviceName": self.serviceName,
					"structName":  structName,
					"funcNames":   funcNames,
				},
			},
		}); err != nil {
			return err
		}
	}

	// 使用gofmt格式化代码
	self.Info("格式化代码")
	_, err := execx.Run("gofmt -l -w "+self.Dir(), "")
	if err != nil {
		return err
	}
	self.Info("生成完毕")
	return nil
}

func getType(vals map[string]interface{}) interface{} {
	typ := vals["type"]
	if t, ok := vals["description"]; ok {
		typ = strings.Split(t.(string), ";")[0]
	}
	return typ
}

func GetGoCode(vals map[string]interface{}) string {
	ret := []string{}
	for name, val := range vals {
		v := val.(map[string]interface{})
		switch v["type"] {
		case "object":
			str := fmt.Sprintf("%s struct{%s} `json:\"%s,omitempty\"`",
				stringx.From(name).ToCamel(),
				GetGoCode(v["properties"].(map[string]interface{})),
				name,
			)
			ret = append(ret, str)
		case "array":
			var str string
			items := v["items"].(map[string]interface{})
			switch items["type"] {
			case "object":
				str = fmt.Sprintf("%s []*struct{%s} `json:\"%s,omitempty\"`",
					stringx.From(name).ToCamel(),
					GetGoCode(items["properties"].(map[string]interface{})),
					name,
				)
			default:
				str = fmt.Sprintf("%s []%s `json:\"%s,omitempty\"`",
					stringx.From(name).ToCamel(),
					getType(items),
					name,
				)
			}
			ret = append(ret, str)
		default:
			tplStr, err := util.LoadTemplate("", "", filedTpl)
			if err != nil {
				return ""
			}
			out, err := util.With("filed").Parse(tplStr).Execute(map[string]interface{}{
				"name":      name,
				"camelName": stringx.From(name).ToCamel(),
				"dot":       "`",
				"type":      getType(v),
				"comment":   v["description"],
			})
			if err != nil {
				return ""
			}
			ret = append(ret, string(out.Bytes()))
		}
	}
	return strings.Join(ret, "\n")
}

func getApiTestTpl(apiInfo *ApiInfo) string {
	switch apiInfo.HandlerType {
	case "notify", "queue":
		return apiTestNotifyAndQueueTpl
	default:
		if len(apiInfo.RspBody.Properties.Data.Properties) == 0 {
			return apiTestRpcNoRspTpl
		}
	}
	return apiTestRpcTpl
}

func (self *Generator) GenTestApi() error {
	data := map[string]interface{}{
		"dot":         "`",
		"apiInfos":    self.apiInfos,
		"serviceName": self.serviceName,
	}
	if err := self.GenItems([]*gen.Item{
		{
			Name: "api",
			Type: gen.Dir,
		},
		{
			Name:     "api/init.go",
			Type:     gen.File,
			Template: apiTestInitTpl,
			Data:     data,
		},
		{
			Name: "benchmark",
			Type: gen.Dir,
		},
		{
			Name:     "benchmark/run_all.sh",
			Type:     gen.File,
			Template: benchmarkRunAllTpl,
			Data:     data,
		},
	}); err != nil {
		return err
	}

	for _, apiInfo := range self.apiInfos {
		if err := self.GenItems([]*gen.Item{
			{
				Name:     "api/" + apiInfo.FileName + "_test.go",
				Type:     gen.File,
				Template: getApiTestTpl(apiInfo),
				Data: map[string]interface{}{
					"apiInfo":     apiInfo,
					"testName":    strings.ReplaceAll(self.serviceName, "-service", "-test"),
					"serviceName": self.serviceName,
				},
			},
		}); err != nil {
			return err
		}

		if apiInfo.HandlerType == "rpc" {
			switch apiInfo.Method {
			case "POST", "PUT":
				if err := self.GenItems([]*gen.Item{
					{
						Name:     "benchmark/" + apiInfo.FileName + ".sh",
						Type:     gen.File,
						Template: benchmarkPostTpl,
						Data: map[string]interface{}{
							"apiInfo":     apiInfo,
							"testName":    strings.ReplaceAll(self.serviceName, "-service", "-test"),
							"serviceName": self.serviceName,
						},
					},
					{
						Name: "benchmark/" + apiInfo.FileName + ".txt",
						Type: gen.File,
					},
				}); err != nil {
					return err
				}
			case "GET":
				if err := self.GenItems([]*gen.Item{
					{
						Name:     "benchmark/" + apiInfo.FileName + ".sh",
						Type:     gen.File,
						Template: benchmarkGetTpl,
						Data: map[string]interface{}{
							"apiInfo":     apiInfo,
							"testName":    strings.ReplaceAll(self.serviceName, "-service", "-test"),
							"serviceName": self.serviceName,
						},
					},
				}); err != nil {
					return err
				}
			}
		}
	}

	// benchmark文件夹提升权限
	if _, err := execx.Run("chmod -R 777 ./benchmark", self.Dir()); err != nil {
		return err
	}

	// 使用gofmt格式化代码
	self.Info("格式化代码")
	_, err := execx.Run("gofmt -l -w "+self.Dir(), "")
	if err != nil {
		return err
	}

	self.Info("生成完毕")
	return nil
}

func (self *Generator) GenCpp() error {
	if err := self.GenItems([]*gen.Item{
		{
			Name: "apim",
			Type: gen.Dir,
		},
	}); err != nil {
		return err
	}

	dir := "apim/"
	for _, apiInfo := range self.apiInfos {
		data := map[string]interface{}{
			"apiInfo": apiInfo,
		}
		if err := self.GenItems([]*gen.Item{
			{
				Name:     dir + apiInfo.FileName + "_req.cpp",
				Type:     gen.File,
				Template: cppReqTpl,
				Data:     data,
			},
		}); err != nil {
			return err
		}

		if len(apiInfo.RspBody.Properties.Data.Properties) > 0 {
			if err := self.GenItems([]*gen.Item{
				{
					Name:     dir + apiInfo.FileName + "_rsp.go",
					Type:     gen.File,
					Template: cppRspTpl,
					Data:     data,
				},
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
