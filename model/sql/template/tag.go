package template

var Tag = "`gorm:\"column:{{.field}}\" {{if .withCache}}json:\"{{.lower_filed}},omitempty\"{{end}}`"
