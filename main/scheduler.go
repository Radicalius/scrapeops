package main

import (
	"encoding/json"
	"reflect"

	"github.com/Radicalius/scrapeops/shared"
)

func RunScheduler() {
	q, err := InitQueue()
	if err != nil {
		baseLogger.Fatal("Error initializing queues", "error", err.Error())
	}

	for {
		id, message, err := q.Peek("main")
		if err != nil {
			baseLogger.Error("Error reading mesage from queue", "error", err.Error())
			continue
		}

		if message == "" {
			continue
		}

		var item shared.Item
		err = json.Unmarshal([]byte(message), &item)
		if err != nil {
			baseLogger.Error("Error unmarshalling", "error", err.Error())
		}

		errorReported := false
		for _, provider := range Providers {
			if errorReported {
				break
			}

			if provider.IsRelevant(&item) {
				dom, err := HttpFetch(provider.GetUrl(&item))
				if err != nil {
					baseLogger.Error("Error fetching page", "url", provider.GetUrl(&item), "error", err.Error())
					errorReported = true
					continue
				}

				newItems, err := provider.Apply(dom, &item)
				if err != nil {
					baseLogger.Error("Error applying provider", "message", message, "provider", reflect.TypeOf(provider).Name(), "error", err.Error())
					errorReported = true
					continue
				}

				for _, newItem := range newItems {
					err := q.Emit("main", newItem)
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
