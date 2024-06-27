package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"plugin"
	"strings"

	scrapeops_plugin "github.com/Radicalius/scrapeops/shared"
)

var Handlers scrapeops_plugin.RawHandlerFuncMap = make(scrapeops_plugin.RawHandlerFuncMap)

func main() {
	pluginDir := os.Getenv("SCRAPEOPS_PLUGIN_DIRECTORY")
	if pluginDir == "" {
		pluginDir = "./plugins"
	}

	files, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		log.Fatalf("Error opening the directory %s: %s", pluginDir, err.Error())
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			p, err := plugin.Open(pluginDir + "/" + file.Name())
			if err != nil {
				fmt.Printf("Error opening plugin %s: %s", file.Name(), err.Error())
				continue
			}

			handlerSym, err := p.Lookup("Handlers")
			if err != nil {
				fmt.Printf("Error loading Handlers symbol in %s: %s", file.Name(), err.Error())
				continue
			}

			handlers := handlerSym.(*scrapeops_plugin.RawHandlerFuncMap)
			if handlers == nil {
				fmt.Printf("Encountered nil Handlers symbol in plugin %s", file.Name())
			}

			for key, f := range *handlers {
				Handlers[key] = f
			}
		}
	}
}
