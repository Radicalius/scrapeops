package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"plugin"
	"strings"

	"github.com/Radicalius/scrapeops/shared"
)

var Providers []shared.Provider = make([]shared.Provider, 0)
var Schematizers []shared.Schematizer[any] = make([]shared.Schematizer[any], 0)
var Crons []shared.CronConfiguration = make([]shared.CronConfiguration, 0)

func LoadPlugins() {
	pluginDir := os.Getenv("SCRAPEOPS_PLUGIN_DIRECTORY")
	if pluginDir == "" {
		pluginDir = "./plugins"
	}

	files, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		baseLogger.Fatal(fmt.Sprintf("Error opening the directory %s", pluginDir), "error", err.Error())
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			p, err := plugin.Open(pluginDir + "/" + file.Name())
			if err != nil {
				baseLogger.Error("Error opening plugin", "plugin", file.Name(), "error", err.Error())
				continue
			}

			pluginSym, err := p.Lookup("PluginConfiguration")
			if err != nil {
				baseLogger.Error("Error loading PluginConfiguration symbol", "plugin", file.Name(), "error", err.Error())
				continue
			}

			plugin := pluginSym.(**shared.PluginConfiguration)
			if plugin == nil {
				baseLogger.Error("Encountered nil PluginConfiguration symbol", "plugin", file.Name())
				continue
			}

			Providers = append(Providers, (*plugin).Providers...)
			Crons = append(Crons, (*plugin).CronConfigs...)
		}
	}
}
