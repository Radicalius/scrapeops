package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"plugin"
	"strings"
	"time"

	scrapeops_plugin "github.com/Radicalius/scrapeops/shared"
	"github.com/robfig/cron"
)

var Handlers scrapeops_plugin.RawHandlerFuncMap = make(scrapeops_plugin.RawHandlerFuncMap)

func main() {
	logger := NewLogger(InitLogCollector()).With("environment", os.Getenv("SCRAPEOPS_ENVIRONMENT"))
	metrics := NewMetrics()
	metrics.InitMetricsApis()

	db := NewDatabaseCollection()

	q, err := InitQueue()
	if err != nil {
		logger.Fatal("Error initializing queues", "error", err.Error())
	}

	context := NewContext(q, db, metrics)

	crons := cron.New()
	crons.AddFunc("0 * * * *", func() {
		err := metrics.Flush()
		if err != nil {
			logger.Error("Flushing metrics", "error", err.Error())
		}
	})

	pluginDir := os.Getenv("SCRAPEOPS_PLUGIN_DIRECTORY")
	if pluginDir == "" {
		pluginDir = "./plugins"
	}

	files, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error opening the directory %s", pluginDir), "error", err.Error())
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			p, err := plugin.Open(pluginDir + "/" + file.Name())
			if err != nil {
				logger.Error("Error opening plugin", "plugin", file.Name(), "error", err.Error())
				continue
			}

			pluginSym, err := p.Lookup("PluginConfiguration")
			if err != nil {
				logger.Error("Error loading PluginConfiguration symbol", "plugin", file.Name(), "error", err.Error())
				continue
			}

			plugin := pluginSym.(**scrapeops_plugin.PluginConfiguration)
			if plugin == nil {
				logger.Error("Encountered nil PluginConfiguration symbol", "plugin", file.Name())
				continue
			}

			for key, f := range (*plugin).Handlers {
				Handlers[key] = f
			}

			if (*plugin).DatabaseConfiguration != nil {
				err = db.AddDatabase((*plugin).DatabaseConfiguration.Name, (*plugin).DatabaseConfiguration.Migrations)
				if err != nil {
					logger.Error("Error loading database", "plugin", file.Name(), "error", err.Error())
					continue
				}
			}

			for cronExpr, jobLists := range (*plugin).CronJobs {
				for i, job := range jobLists {
					handlerName := fmt.Sprintf("cron%s-%d", cronExpr, i)
					Handlers[handlerName] = scrapeops_plugin.RawHandlerFunc(func(data []byte, ctx scrapeops_plugin.Context) error {
						return job(context)
					})

					crons.AddFunc(cronExpr, func() {
						context.GetQueue().Emit(handlerName, "")
					})
				}
			}

			for route, apiFunc := range (*plugin).Apis {
				InitApi(route, apiFunc, context, logger)
			}
		}
	}

	Handlers["httpAsync"] = scrapeops_plugin.ToRawHandlerFunc(HttpAsyncHandler)

	crons.Start()
	go http.ListenAndServe(":8080", nil)

	for {
		for handlerName := range Handlers {
			messageId, messageBody, err := q.Peek(handlerName)
			handlerLogger := logger.With("queue", handlerName)
			if err != nil {
				handlerLogger.Error("Error peeking at queue", "error", err.Error())
				continue
			}

			if messageBody == "" {
				continue
			}

			go func(handlerName string, logger_ *Logger) {
				err := Handlers[handlerName]([]byte(messageBody), context.WithLogger(logger_))
				if err != nil {
					logger_.Error("Error processing message", "error", err.Error())
					metrics.IncrementCounter("processor_" + handlerName + "_errors")
					return
				}

				metrics.IncrementCounter("processor_" + handlerName + "_successes")

				err = q.Delete(messageId)
				if err != nil {
					logger_.Error("Error deleting message", "error", err.Error())
				}
			}(handlerName, handlerLogger.With("queueMessage", messageBody))
		}

		time.Sleep(1 * time.Second)
	}
}
