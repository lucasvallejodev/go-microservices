package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth, omoitempty"`
	Log    LogPayload  `json:"log, omoitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Broker service is running",
	}

	_ = app.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmition(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authentica(w, requestPayload.Auth)
	case "log":
		app.logEntry(w, requestPayload.Log)
	default:
		app.errorJSON(w, errors.New("unknown action"), http.StatusBadRequest)
	}
}

func (app *Config) authentica(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://auth-server/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"), response.StatusCode)
		return
	}
	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("error calling auth server"), response.StatusCode)
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, errors.New(jsonFromService.Message), response.StatusCode)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = jsonFromService.Message
	payload.Data = jsonFromService.Data

	app.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) logEntry(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-server/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("error calling logger server"), response.StatusCode)
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, errors.New(jsonFromService.Message), response.StatusCode)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = jsonFromService.Message

	app.WriteJSON(w, http.StatusOK, payload)
}
