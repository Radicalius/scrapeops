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
	"github.com/robfig/cron"
)

var Handlers scrapeops_plugin.RawHandlerFuncMap = make(scrapeops_plugin.RawHandlerFuncMap)

func main() {

	db := NewDatabaseCollection()

	q, err := InitQueue()
	if err != nil {
		log.Fatalf("Error initializing queue: %s", err.Error())
	}

	context := NewContext(q, db)

	crons := cron.New()

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

			pluginSym, err := p.Lookup("PluginConfiguration")
			if err != nil {
				fmt.Printf("Error loading Handlers symbol in %s: %s", file.Name(), err.Error())
				continue
			}

			plugin := pluginSym.(**scrapeops_plugin.PluginConfiguration)
			if plugin == nil {
				fmt.Printf("Encountered nil Handlers symbol in plugin %s", file.Name())
			}

			for key, f := range (*plugin).Handlers {
				Handlers[key] = f
			}

			if (*plugin).DatabaseConfiguration != nil {
				err = db.AddDatabase((*plugin).DatabaseConfiguration.Name, (*plugin).DatabaseConfiguration.Migrations)
				if err != nil {
					fmt.Printf("Error loading database for plugin %s: %s", file.Name(), err.Error())
				}
			}

			for cronExpr, jobLists := range (*plugin).CronJobs {
				for _, job := range jobLists {
					crons.AddFunc(cronExpr, func() {
						job(context)
					})
				}
			}
		}
	}

	Handlers["httpAsync"] = scrapeops_plugin.ToRawHandlerFunc(HttpAsyncHandler)

	crons.Start()

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

			go func(handlerName string) {
				err := Handlers[handlerName]([]byte(messageBody), context)
				if err != nil {
					fmt.Printf("Error processing message: \n\tqueue: %s\n\tmessage: %s\n\terror: %s\n", handlerName, string(messageBody), err.Error())
					return
				}

				err = q.Delete(messageId)
				if err != nil {
					fmt.Printf("Error deleting message: %s", err.Error())
				}
			}(handlerName)
		}

		time.Sleep(1 * time.Second)
	}
}
