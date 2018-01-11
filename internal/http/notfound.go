package http

import (
	"net/http"
	"text/template"

	"github.com/romana/rlog"
)

// NotFound writes out the not found xml response to freeswitch to let it know to move on to the file system check
func notFound(w http.ResponseWriter) {
	t, err := template.ParseFiles(notFoundTemplatePath)
	if err != nil {
		rlog.Errorf("could not load notfound template [%s]", err.Error())
		return
	}
	t.ExecuteTemplate(w, notFoundTemplate, nil)
	return
}
