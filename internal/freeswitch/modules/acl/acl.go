package acl

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/romana/rlog"
)

const (
	moduleDataFile = "acl.json"
	configTemplate = "configuration/acl/acl.xml"
)

var (
	moduleSettingFile string
	templatePath      string
)

type aclNode struct {
	Action string `json:"action"`
	Type   string `json:"type"`
	Value  string `json:"value"`
}

type aclList struct {
	Name   string    `json:"name"`
	Action string    `json:"action"`
	Nodes  []aclNode `json:"nodes"`
}

type module struct {
	Lists []aclList `json:"acl.conf"`
}

type host map[string]module

func New(m string, t string) error {
	moduleSettingFile = filepath.Join(m, moduleDataFile)
	templatePath = filepath.Join(t, configTemplate)
	rlog.Infof("set acl module settings file [%s]", moduleSettingFile)
	rlog.Infof("set acl template path [%s]", templatePath)
	return nil
}

func Handler(ctx context.Context, hostname string, w http.ResponseWriter) error {
	rlog.Debugf("configuration request for hostname [%s]", hostname)

	h := host{}
	d, err := ioutil.ReadFile(moduleSettingFile)
	if err != nil {
		rlog.Errorf("could not read file [%s]", err.Error())
		return err
	}
	if err = json.Unmarshal(d, &h); err != nil {
		rlog.Errorf("could not unmarshal file [%s]", err.Error())
		return err
	}
	m, ok := h[hostname]
	if !ok {
		rlog.Infof("hostname not found [%s]", hostname)
		return errors.New("hostname not found")
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		rlog.Errorf("could not parse template file [%s]", err.Error())
		return err
	}
	t.Execute(w, m)
	return nil
}
