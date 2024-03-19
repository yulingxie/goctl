package gen

import (
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/sql/template"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util"
)

func genTypes(table Table, withCache bool) (string, error) {
	fields := table.Fields
	fieldsString, err := genFields(fields, withCache)
	if err != nil {
		return "", err
	}

	text, err := util.LoadTemplate(category, typesTemplateFile, template.Types)
	if err != nil {
		return "", err
	}

	output, err := util.With("types").
		Parse(text).
		Execute(map[string]interface{}{
			"withCache":             withCache,
			"upperStartCamelObject": table.Name.ToCamel(),
			"fields":                fieldsString,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
