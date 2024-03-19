package gen

import (
	"fmt"

	"github.com/urfave/cli"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/model/sql/template"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util"
)

const (
	category                          = "model"
	fieldTemplateFile                 = "field.tpl"
	importsWithCacheTemplateFile      = "import-with-cache.tpl"
	importsWithNoCacheTemplateFile    = "import-no-cache.tpl"
	modelWithCacheTemplateFile        = "modelWithCache.tpl"
	modelNoCacheTemplateFile          = "modelNoCache.tpl"
	newTemplateFile                   = "new.tpl"
	newWithCacheTemplateFile          = "new-with-cache.tpl"
	newWithGameIdTemplateFile         = "new-with-gameid.tpl"
	newWithCacheAndGameIdTemplateFile = "new-with-cache-and-gameid.tpl"
	tagTemplateFile                   = "tag.tpl"
	typesTemplateFile                 = "types.tpl"
	methodsTemplateFile               = "methods.tpl"
	methodsWithCacheTemplateFile      = "methods-with-cache.tpl"
)

var templates = map[string]string{
	fieldTemplateFile:                 template.Field,
	importsWithCacheTemplateFile:      template.ImportsWithCache,
	importsWithNoCacheTemplateFile:    template.ImportsNoCache,
	modelWithCacheTemplateFile:        template.ModelWithCache,
	modelNoCacheTemplateFile:          template.ModelNoCache,
	newTemplateFile:                   template.New,
	newWithCacheTemplateFile:          template.NewWithCache,
	newWithGameIdTemplateFile:         template.NewWithGameId,
	newWithCacheAndGameIdTemplateFile: template.NewWitchCacheAndGameId,
	tagTemplateFile:                   template.Tag,
	typesTemplateFile:                 template.Types,
	methodsTemplateFile:               template.Methods,
	methodsWithCacheTemplateFile:      template.MethodsWithCache,
}

// Category returns model const value
func Category() string {
	return category
}

// Clean deletes all template files
func Clean() error {
	return util.Clean(category)
}

// GenTemplates creates template files if not exists
func GenTemplates(_ *cli.Context) error {
	return util.InitTemplates(category, templates)
}

// RevertTemplate recovers the delete template files
func RevertTemplate(name string) error {
	content, ok := templates[name]
	if !ok {
		return fmt.Errorf("%s: no such file name", name)
	}

	return util.CreateTemplate(category, name, content)
}

// Update provides template clean and init
func Update() error {
	err := Clean()
	if err != nil {
		return err
	}

	return util.InitTemplates(category, templates)
}
