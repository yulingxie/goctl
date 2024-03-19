package mongo

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
	nodeName string
	dbName   string
	collName string
}

func NewGenerator(dir, nodeName, dbName, collName string) (*Generator, error) {
	generator, err := gen.NewGenerator(dir)
	if err != nil {
		return nil, err
	}
	return &Generator{Generator: generator, nodeName: nodeName, dbName: dbName, collName: collName}, nil
}

func (self *Generator) GenModel() error {
	collNames := strings.Split(self.collName, ",")
	for _, collName := range collNames {
		filename, err := format.FileNamingFormat("go_zero", collName+"_model")
		typeName := stringx.From(collName).ToCamel()
		if err != nil {
			return err
		}

		text, err := util.LoadTemplate("model", "model.tpl", modelTpl)
		if err != nil {
			return err
		}

		output := filepath.Join(self.Dir(), filename+".go")
		err = util.With("model").Parse(text).GoFmt(true).SaveTo(map[string]interface{}{
			"nodeName":      self.nodeName,
			"dbName":        self.dbName,
			"collName":      collName,
			"typeName":      typeName,
			"lowerTypeName": stringx.From(typeName).Untitle(),
		}, output, false)
		if err != nil {
			return err
		}
		self.Info("生成文件: %v", output)
	}
	return nil
}
