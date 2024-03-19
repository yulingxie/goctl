package gen

import (
	"testing"
)

const testTpl = `
	{{.name}}
`

func TestGenerator_GenItems(t *testing.T) {
	gen, _ := NewGenerator("./test/")
	gen.GenItems([]*Item{
		{"prased_file.go", File, testTpl, map[string]interface{}{"name": "test"}, false},
		{"not_prased_file.go", File, testTpl, map[string]interface{}{"name": "test"}, true},
	})
}
