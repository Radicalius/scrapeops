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

func HandleError(route string, input string, err error, w http.ResponseWriter) {
	data, _ := json.Marshal(ErrorResponse{
		Error: err.Error(),
	})

	fmt.Printf("Error handling api call:\n\troute: %s\n\tinput: %s\n\terror: %s\n", route, input, err.Error())

	w.WriteHeader(500)
	w.Write(data)
}

func InitApi(route string, apiFunc sops.RawApiFunc, context sops.Context) {
	http.HandleFunc(fmt.Sprintf("/api/%s", route), func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			HandleError(route, "", err, w)
			return
		}

		out, err := apiFunc(body, context)
		if err != nil {
			HandleError(route, string(body), err, w)
			return
		}

		w.Write(*out)
	})
}
