package template

// Model defines a template for model
var ModelWithCache = `package cachem
{{.imports}}
{{.types}}
{{.new}}
{{.methods}}
`

var ModelNoCache = `package sqlm
{{.imports}}
{{.types}}
{{.new}}
{{.methods}}
`
