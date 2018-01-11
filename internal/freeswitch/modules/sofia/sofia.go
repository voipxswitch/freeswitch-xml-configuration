package sofia

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
	moduleDataFile = "sofia.json"
	configTemplate = "configuration/sofia/sofia.xml"
)

var (
	moduleSettingFile string
	templatePath      string
)

type settings struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type gateways struct {
	Name     string     `json:"name"`
	Settings []settings `json:"settings"`
}

type profiles struct {
	Name     string     `json:"name"`
	Gateways []gateways `json:"gateways"`
	Settings []settings `json:"settings"`
}

type sofia struct {
	Globals  []settings `json:"globals"`
	Profiles []profiles `json:"profiles"`
}

type module struct {
	Sofia sofia `json:"sofia.conf"`
}

type host map[string]module

func New(m string, t string) error {
	moduleSettingFile = filepath.Join(m, moduleDataFile)
	templatePath = filepath.Join(t, configTemplate)
	rlog.Infof("set module settings file [%s]", moduleSettingFile)
	rlog.Infof("set template path [%s]", templatePath)
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
	t.Execute(w, m.Sofia)
	return nil
}
