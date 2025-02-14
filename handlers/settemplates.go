package handlers

import "html/template"

var Tmpl *template.Template

func SetTemplates(t *template.Template) {
	Tmpl = t
}
