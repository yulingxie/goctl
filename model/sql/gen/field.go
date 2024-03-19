package gen

import (
	"strings"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/sql/parser"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/sql/template"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util"
)

func genFields(fields []*parser.Field, withCache bool) (string, error) {
	var list []string

	for _, field := range fields {
		result, err := genField(field, withCache)
		if err != nil {
			return "", err
		}

		list = append(list, result)
	}

	return strings.Join(list, "\n"), nil
}

func genField(field *parser.Field, withCache bool) (string, error) {
	tag, err := genTag(field.Name.Source(), withCache)
	if err != nil {
		return "", err
	}

	text, err := util.LoadTemplate(category, fieldTemplateFile, template.Field)
	if err != nil {
		return "", err
	}

	output, err := util.With("types").
		Parse(text).
		Execute(map[string]interface{}{
			"name":       field.Name.ToCamel(),
			"type":       field.DataType,
			"tag":        tag,
			"hasComment": field.Comment != "",
			"comment":    field.Comment,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
