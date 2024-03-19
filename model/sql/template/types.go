package template

// Types defines a template for types in model
var Types = `
type (
	{{.upperStartCamelObject}}Model struct {
		{{if .withCache}}cache  cache.Cache{{end}}
		{{if .withCache}}cacheKey  string{{end}}
		client *sqlx.Client
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
`
