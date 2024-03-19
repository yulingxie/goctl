package gen

import (
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/sql/template"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util"
)

func genImports(withCache, withGameId, timeImport bool) (string, error) {
	file := importsWithNoCacheTemplateFile
	tpl := template.ImportsNoCache
	if withCache {
		file = importsWithCacheTemplateFile
		tpl = template.ImportsWithCache
	}

	text, err := util.LoadTemplate(category, file, tpl)
	if err != nil {
		return "", err
	}

	buffer, err := util.With("import").Parse(text).Execute(map[string]interface{}{
		"time":       timeImport,
		"withGameId": withGameId,
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
