package gen

import (
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/sql/template"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/stringx"
)

func genTag(in string, withCache bool) (string, error) {
	if in == "" {
		return in, nil
	}

	text, err := util.LoadTemplate(category, tagTemplateFile, template.Tag)
	if err != nil {
		return "", err
	}

	output, err := util.With("tag").Parse(text).Execute(map[string]interface{}{
		"field":       in,
		"lower_filed": stringx.From(in).Lower(),
		"withCache":   withCache,
	})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
