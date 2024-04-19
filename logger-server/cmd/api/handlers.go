package main

import (
	"log-server/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload

	err := app.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	entry := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.Models.LogEntry.Insert(entry)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: "Log written",
	}

	app.WriteJSON(w, http.StatusOK, response)
}
