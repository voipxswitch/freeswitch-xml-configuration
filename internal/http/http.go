package http

import (
	"net/http"
	"path/filepath"

	"github.com/romana/rlog"
	"goji.io"
	"goji.io/pat"

	"github.com/voipxswitch/freeswitch-xml-configuration/internal/freeswitch/modules/acl"
	"github.com/voipxswitch/freeswitch-xml-configuration/internal/freeswitch/modules/distributor"
	"github.com/voipxswitch/freeswitch-xml-configuration/internal/freeswitch/modules/sofia"
)

const (
	requestPath       = "/fs/*"
	notFoundTemplate  = "notfound.xml"
)

var (
	h httpHandler

	notFoundTemplatePath string
)

type httpHandler struct {}

func New(root *goji.Mux, httpAddress string, moduleDataDirectoryCfg string, templatePath string) error {
	// setup freeswitch modules
	err := acl.New(moduleDataDirectoryCfg, templatePath)
	if err != nil {
		return err
	}
	rlog.Info("setup acl module")
	err = distributor.New(moduleDataDirectoryCfg, templatePath)
	if err != nil {
		return err
	}
	rlog.Info("setup distributor module")
	err = sofia.New(moduleDataDirectoryCfg, templatePath)
	if err != nil {
		return err
	}
	rlog.Info("setup sofia module")

	// setup not found template
	notFoundTemplatePath = filepath.Join(templatePath, notFoundTemplate)
	rlog.Infof("set not found template path [%s]", notFoundTemplatePath)

	// setup http handler
	v := goji.SubMux()
	root.Handle(pat.New(requestPath), v)
	rlog.Debugf("registered http handler [%s]", requestPath)
	registerMux(v)
	return http.ListenAndServe(httpAddress, root)
}

func registerMux(m *goji.Mux) {
	m.HandleFunc(pat.Post("/configuration"), configuration.Handler)
	rlog.Debug("registered configuration endpoint")
}
