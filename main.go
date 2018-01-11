package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/romana/rlog"
	"goji.io"

	"github.com/voipxswitch/freeswitch-xml-configuration/internal/http"
)

var (
	confPath          = "config.json"
	listenAddressHttp = ":8001"
)

func init() {
	rlog.SetOutput(os.Stdout)
}

func main() {
	rlog.Debug("debug enabled")

	configFilePath := flag.String("config", "", "path to config file")
	flag.Parse()
	if *configFilePath != "" {
		confPath = *configFilePath
	}
	rlog.Infof("loading config from file [%s]", confPath)

	c, err := loadConfigFile(confPath)
	if err != nil {
		rlog.Errorf("could load config [%s]", err.Error())
		os.Exit(1)
	}

	// http settings
	if c.HTTP.ListenHTTP != "" {
		listenAddressHttp = c.HTTP.ListenHTTP
	}

	// start http
	err = http.New(goji.NewMux(), listenAddressHttp, c.FreeSWITCH.ModuleDataDirectory, c.HTTP.TemplatesDir)
	if err != nil {
		rlog.Errorf("could not start http(s) server [%s]", err.Error())
		os.Exit(1)
	}
}

// struct used to unmarshal config.json
type serviceConfig struct {
	HTTP struct {
		TemplatesDir string `json:"templates_directory"`
		ListenHTTP   string `json:"listen"`
	} `json:"http"`
	FreeSWITCH struct {
		ModuleDataDirectory string `json:"module_data_directory"`
	} `json:"freeswitch"`
}

func loadConfigFile(configFile string) (serviceConfig, error) {
	s := serviceConfig{}
	file, err := os.Open(configFile)
	if err != nil {
		return s, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s)
	if err != nil {
		return s, err
	}
	return s, nil
}
