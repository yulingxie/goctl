package gen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/sql/model"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/sql/parser"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/sql/template"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/console"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/format"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/stringx"
)

const (
	pwd             = "."
	createTableFlag = `(?m)^(?i)CREATE\s+TABLE` // ignore case
)

type (
	defaultGenerator struct {
		console.Console
		dir  string
		pkg  string
		node string
	}

	code struct {
		importsCode string
		typesCode   string
		newCode     string
		methodsCode string
	}

	Table struct {
		parser.Table
		PrimaryCacheKey        Key
		UniqueCacheKey         []Key
		ContainsUniqueCacheKey bool
	}
)

func NewDefaultGenerator(dir, node string) (*defaultGenerator, error) {
	if dir == "" {
		dir = pwd
	}
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	dir = dirAbs
	pkg := filepath.Base(dirAbs)
	err = util.MkdirIfNotExist(dir)
	if err != nil {
		return nil, err
	}

	generator := &defaultGenerator{node: node, dir: dir, pkg: pkg, Console: console.NewColorConsole()}
	return generator, nil
}

func (g *defaultGenerator) StartFromInformationSchema(tableDatas map[string]*model.TableData, withCache bool) error {
	m := make(map[string]string)
	for _, each := range tableDatas {
		table, err := parser.ConvertDataType(each)
		if err != nil {
			return err
		}

		// 如果db名为yygsubshuangsheng,则按游戏数据库模版生成
		withGameId := table.DbName == "yygsubshuangsheng"
		code, err := g.genModel(*table, withCache, withGameId)
		if err != nil {
			return err
		}

		m[table.Name.Source()] = code
	}

	return g.createFile(m)
}

func (g *defaultGenerator) createFile(modelList map[string]string) error {
	dirAbs, err := filepath.Abs(g.dir)
	if err != nil {
		return err
	}

	g.dir = dirAbs
	g.pkg = filepath.Base(dirAbs)
	err = util.MkdirIfNotExist(dirAbs)
	if err != nil {
		return err
	}

	for tableName, code := range modelList {
		// 生成model文件
		tn := stringx.From(tableName)
		modelFilename, err := format.FileNamingFormat("go_zero", fmt.Sprintf("%s_model", tn.Source()))
		if err != nil {
			return err
		}

		for util.FileExists(filepath.Join(dirAbs, modelFilename+".go")) {
			// 若源文件已存在,则在文件名后加new
			modelFilename += "_new"
		}
		filename := filepath.Join(dirAbs, modelFilename+".go")
		err = ioutil.WriteFile(filename, []byte(code), os.ModePerm)
		if err != nil {
			return err
		}
		g.Info(fmt.Sprintf("生成model文件: %s", g.dir+filename))
	}

	g.Success("Done.")
	return nil
}

func (g *defaultGenerator) genModel(in parser.Table, withCache, withGameId bool) (string, error) {
	if len(in.PrimaryKey.Name.Source()) == 0 {
		return "", fmt.Errorf("table %s: missing primary key", in.Name.Source())
	}

	primaryKey, uniqueKey := genCacheKeys(in)

	importsCode, err := genImports(withCache, withGameId, in.ContainsTime())
	if err != nil {
		return "", err
	}

	table := Table{
		Table:                  in,
		PrimaryCacheKey:        primaryKey,
		UniqueCacheKey:         uniqueKey,
		ContainsUniqueCacheKey: len(uniqueKey) > 0,
	}

	typesCode, err := genTypes(table, withCache)
	if err != nil {
		return "", err
	}

	newCode, err := genNew(g.node, table, withCache, withGameId)
	if err != nil {
		return "", err
	}

	methodsCode, err := genMethods(table, withCache)
	if err != nil {
		return "", err
	}

	code := &code{
		importsCode: importsCode,
		typesCode:   typesCode,
		newCode:     newCode,
		methodsCode: methodsCode,
	}

	output, err := g.executeModel(code, withCache)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}

func (g *defaultGenerator) executeModel(code *code, withCache bool) (*bytes.Buffer, error) {
	file := modelNoCacheTemplateFile
	tpl := template.ModelNoCache
	if withCache {
		file = modelWithCacheTemplateFile
		tpl = template.ModelWithCache
	}

	text, err := util.LoadTemplate(category, file, tpl)
	if err != nil {
		return nil, err
	}
	output, err := util.With("model").Parse(text).GoFmt(true).Execute(map[string]interface{}{
		"pkg":     g.pkg,
		"imports": code.importsCode,
		"types":   code.typesCode,
		"new":     code.newCode,
		"methods": code.methodsCode,
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func wrapWithRawString(v string) string {
	if v == "`" {
		return v
	}
	if !strings.HasPrefix(v, "`") {
		v = "`" + v
	}
	if !strings.HasSuffix(v, "`") {
		v = v + "`"
	} else if len(v) == 1 {
		v = v + "`"
	}
	return v
}
