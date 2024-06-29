package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"plugin"
	"strings"
	"time"

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

	q, err := InitQueue()
	if err != nil {
		log.Fatalf("Error initializing queue: %s", err.Error())
	}

	for {
		for handlerName := range Handlers {
			messageId, messageBody, err := q.Peek(handlerName)
			if err != nil {
				fmt.Printf("Error peeking at queue for %s; %s", handlerName, err.Error())
				continue
			}

			if messageBody == "" {
				continue
			}

			go func() {
				err := Handlers[handlerName]([]byte(messageBody), nil)
				if err != nil {
					fmt.Printf("Error processing message: \n\tqueue: %s\n\tmessage: %s\n\terror: %s", handlerName, string(messageBody), err.Error())

				}

				err = q.Delete(messageId)
				if err != nil {
					fmt.Printf("Error deleting message: %s", err.Error())
				}
			}()
		}

		time.Sleep(1 * time.Second)
	}
}
