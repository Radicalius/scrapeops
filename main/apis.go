package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	sops "github.com/Radicalius/scrapeops/shared"
)

type ErrorResponse struct {
	Error string
}

func HandleError(logger *Logger, err error, w http.ResponseWriter) {
	data, _ := json.Marshal(ErrorResponse{
		Error: err.Error(),
	})

	logger.Error("Error handling api call", "error", err.Error())

	w.WriteHeader(500)
	w.Write(data)
}

func InitApi(route string, apiFunc sops.RawApiFunc, context *Context, baseLogger *Logger) {
	http.HandleFunc(fmt.Sprintf("/api/custom/%s", route), func(w http.ResponseWriter, r *http.Request) {
		logger := baseLogger.With("route", route)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			HandleError(logger, err, w)
			return
		}

		logger = logger.With("messageBody", string(body))

		out, err := apiFunc(body, context.WithLogger(logger))
		if err != nil {
			HandleError(logger, err, w)
			return
		}

		w.Write(*out)
	})
}
