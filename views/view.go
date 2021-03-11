package views

import (
	"text/template"
	"time"
)

// Tmpl is a template.
var Tmpl *template.Template

func date(t *time.Time) string {
	return t.Local().Format("January 2, 2006 15:04:05")
}

func shortDate(t *time.Time) string {
	return t.Local().Format("January 2, 2006")
}

func truncate(s string) string {
	result := s

	if len(s) > 160 {
		result = s[0:160] + "..."
	}

	return result
}

var functions = template.FuncMap{
	"date":      date,
	"shortDate": shortDate,
	"truncate":  truncate,
}

func init() {
	Tmpl = template.Must(template.New("views/*.gohtml").Funcs(functions).ParseGlob("views/*.gohtml"))
}
