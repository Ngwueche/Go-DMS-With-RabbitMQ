package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Application) broker(w http.ResponseWriter, r *http.Request) {
	payload := JsonResponse{
		Error:   false,
		Message: "Broker service is up and running",
	}

	err := app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *Application) handleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logRequest(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)

	default:
		app.errorJSON(w, errors.New("Invalid action"))
	}
}

func (app *Application) authenticate(w http.ResponseWriter, authPayload AuthPayload) {

	jsonData, err := json.Marshal(authPayload)
	if err != nil {
		app.errorJSON(w, errors.New("Error marshaling authentication data"), http.StatusInternalServerError)
		return
	}
	request, err := http.NewRequest("POST", "http://authentication-service:8080/authenticate", bytes.NewBuffer(jsonData))
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

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("Invalid credentials"), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusAccepted {
		log.Printf("Unexpected status code from authentication service: %s", err)
		app.errorJSON(w, errors.New("Error calling authentication service"))
		return
	}
	var serverResponse JsonResponse
	err = json.NewDecoder(response.Body).Decode(&serverResponse)
	if err != nil {
		app.errorJSON(w, errors.New("Error decoding authentication response"), http.StatusInternalServerError)
		return
	}
	if serverResponse.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}
	var payload JsonResponse
	payload.Error = false
	payload.Message = "Authentication successful"
	payload.Data = serverResponse.Data

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, errors.New("Error writing authentication response"), http.StatusInternalServerError)
		return
	}
}

func (app *Application) logRequest(w http.ResponseWriter, logPayload LogPayload) {

	jsonData, err := json.Marshal(logPayload)
	if err != nil {
		app.errorJSON(w, errors.New("Error marshaling log data"), http.StatusInternalServerError)
		return
	}
	request, err := http.NewRequest("POST", "http://logger-service:8080/write-log", bytes.NewBuffer(jsonData))
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

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}
	var serverResponse JsonResponse
	err = json.NewDecoder(response.Body).Decode(&serverResponse)
	if err != nil {
		app.errorJSON(w, errors.New("Error decoding log response"), http.StatusInternalServerError)
		return
	}
	if serverResponse.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}
	var payload JsonResponse
	payload.Error = false
	payload.Message = "Log entry created successfully"

	err = app.writeJSON(w, http.StatusAccepted, payload)
	if err != nil {
		http.Error(w, "Internal Server Error	", http.StatusInternalServerError)
		return
	}
}
func (app *Application) sendMail(w http.ResponseWriter, sendPayload MailPayload) {

	jsonData, err := json.Marshal(sendPayload)
	if err != nil {
		app.errorJSON(w, errors.New("Error marshaling log data"), http.StatusInternalServerError)
		return
	}
	request, err := http.NewRequest("POST", "http://mailer-service:8080/send-mail", bytes.NewBuffer(jsonData))
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

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}
	var serverResponse JsonResponse
	err = json.NewDecoder(response.Body).Decode(&serverResponse)
	if err != nil {
		app.errorJSON(w, errors.New("Error calling mailer service"), http.StatusInternalServerError)
		return
	}
	if serverResponse.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}
	var payload JsonResponse
	payload.Error = false
	payload.Message = "Email sent successfully"

	err = app.writeJSON(w, http.StatusAccepted, payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
