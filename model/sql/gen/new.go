package gen

import (
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/sql/template"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/stringx"
)

func genNew(node string, table Table, withCache, withGameId bool) (string, error) {
	file := newTemplateFile
	tpl := template.New
	if withCache && !withGameId {
		file = newWithCacheTemplateFile
		tpl = template.NewWithCache
	} else if !withCache && withGameId {
		file = newWithGameIdTemplateFile
		tpl = template.NewWithGameId
	} else if withCache && withGameId {
		file = newWithCacheAndGameIdTemplateFile
		tpl = template.NewWitchCacheAndGameId
	}

	text, err := util.LoadTemplate(category, file, tpl)
	if err != nil {
		return "", err
	}

	output, err := util.With("new").
		Parse(text).
		Execute(map[string]interface{}{
			"nodeName":                  node,
			"dbName":                    table.DbName,
			"table":                     table.Name.Source(),
			"withCache":                 withCache,
			"lowerStartCamelPrimaryKey": stringx.From(table.PrimaryKey.Name.ToCamel()).Untitle(),
			"upperStartCamelObject":     table.Name.ToCamel(),
			"lowerCamelTableName":       stringx.From(table.Name.ToCamel()).Untitle(),
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
