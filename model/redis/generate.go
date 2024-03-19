package redis

import (
	"path/filepath"
	"strings"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/format"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/gen"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/stringx"
)

type Generator struct {
	*gen.Generator
	nodeName  string
	tableName string
}

func NewGenerator(dir, nodeName, tableName string) (*Generator, error) {
	generator, err := gen.NewGenerator(dir)
	if err != nil {
		return nil, err
	}
	return &Generator{Generator: generator, nodeName: nodeName, tableName: tableName}, nil
}

func (self *Generator) GenModel() error {
	tableNames := strings.Split(self.tableName, ",")
	for _, tableName := range tableNames {
		filename, err := format.FileNamingFormat("go_zero", tableName+"_model")
		if err != nil {
			return err
		}

		var modelTpl string = HmapModelTpl
		if tableName == "Common" {
			modelTpl = CommonModelTpl
		}

		text, err := util.LoadTemplate("model", "model.tpl", modelTpl)
		if err != nil {
			return err
		}

		output := filepath.Join(self.Dir(), filename+".go")
		err = util.With("model").Parse(text).GoFmt(true).SaveTo(map[string]interface{}{
			"nodeName":  self.nodeName,
			"type":      tableName,
			"lowerType": stringx.From(tableName).Untitle(),
		}, output, false)
		if err != nil {
			return err
		}
		self.Info("生成文件: %v", output)
	}
	return nil
}
