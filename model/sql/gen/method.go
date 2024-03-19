package gen

import (
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/sql/template"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/stringx"
)

func genMethods(table Table, withCache bool) (string, error) {
	file := methodsTemplateFile
	tpl := template.Methods
	if withCache {
		file = methodsWithCacheTemplateFile
		tpl = template.MethodsWithCache
	}
	text, err := util.LoadTemplate(category, file, tpl)
	if err != nil {
		return "", err
	}
	camel := table.Name.ToCamel()

	output, err := util.With("methods").
		Parse(text).
		Execute(map[string]interface{}{
			"withCache":                 withCache,
			"upperStartCamelObject":     camel,
			"lowerStartCamelObject":     stringx.From(camel).Untitle(),
			"originalPrimaryKey":        wrapWithRawString(table.PrimaryKey.Name.Source()),
			"upperStartCamelPrimaryKey": table.PrimaryKey.Name.ToCamel(),
			"lowerStartCamelPrimaryKey": stringx.From(table.PrimaryKey.Name.ToCamel()).Untitle(),
			"dataType":                  table.PrimaryKey.DataType,
			"cacheKey":                  table.PrimaryCacheKey.KeyExpression,
			"cacheKeyVariable":          table.PrimaryCacheKey.KeyLeft,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
