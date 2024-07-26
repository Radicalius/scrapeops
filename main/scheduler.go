package main

import (
	"reflect"

	"github.com/Radicalius/scrapeops/shared"
	"github.com/robfig/cron"
)

func RunScheduler() {
	q, err := InitQueue()
	if err != nil {
		baseLogger.Fatal("Error initializing queues", "error", err.Error())
	}

	cron := cron.New()
	for _, cronConfig := range Crons {
		cron.AddFunc(cronConfig.Schedule, func() {
			data, err := cronConfig.Item.Serialize()
			if err != nil {
				baseLogger.Fatal("Error serializing message", "error", err.Error())
				return
			}

			err = q.Emit("main", data)
			if err != nil {
				baseLogger.Fatal("Error pushing message to queue from cron job", "error", err.Error())
				return
			}
		})
	}
	cron.Start()

	for {
		id, message, err := q.Peek("main")
		if err != nil {
			baseLogger.Error("Error reading mesage from queue", "error", err.Error())
			continue
		}

		if message == "" {
			continue
		}

		item, err := shared.DeserializeItem([]byte(message))
		if err != nil {
			baseLogger.Error("Error unmarshalling", "error", err.Error())
			continue
		}

		errorReported := false
		for _, provider := range Providers {
			if errorReported {
				break
			}

			if provider.IsRelevant(item) {
				dom, err := HttpFetch(provider.GetUrl(item))
				if err != nil {
					baseLogger.Error("Error fetching page", "url", provider.GetUrl(item), "error", err.Error())
					errorReported = true
					continue
				}

				newItems, err := provider.Apply(dom, item)
				if err != nil {
					baseLogger.Error("Error applying provider", "message", message, "provider", reflect.TypeOf(provider).Name(), "error", err.Error())
					errorReported = true
					continue
				}

				for _, newItem := range newItems {
					data, err := newItem.Serialize()
					if err != nil {
						baseLogger.Fatal("Error serializing message", "error", err.Error())
						return
					}

					err = q.Emit("main", data)
					if err != nil {
						baseLogger.Error("Error emitting new item", "error", err.Error())
						errorReported = true
						continue
					}
				}
			}
		}

		err = q.Delete(id)
		if err != nil {
			baseLogger.Error("Error deleting item from queue", "error", err.Error())
		}
	}
}
