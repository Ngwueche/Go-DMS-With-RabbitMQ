package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Application) SendMail(w http.ResponseWriter, r *http.Request) {
	type MailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}
	log.Printf("This is mailer service about to return")
	logToMongo := app.LogRequest("mail", fmt.Sprintf("Mailer service got to this point1"))
	if logToMongo != nil {
		app.errorJSON(w, errors.New("Failed to log authentication request1"), http.StatusInternalServerError)
		return
	}
	var requestPayload MailMessage
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	//--------log
	log.Printf("This is mailer service about to return")
	err = app.LogRequest("mail", fmt.Sprintf("Mailer service got to this point2"))
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}
	err = app.LogRequest("mail", fmt.Sprintf("Mailer service got to this point3"))
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	err = app.Mailer.SendSMTPMessage(msg)
	err = app.LogRequest("mail", fmt.Sprintf("Mailer just sent"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	
	log.Printf("This is mailer service about to return")
	err = app.LogRequest("mail", fmt.Sprintf("Mailer service got to this point3"))
	if err != nil {
		app.errorJSON(w, errors.New("Failed to log authentication request"), http.StatusInternalServerError)
		return
	}
	payload := JsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Application) LogRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, _ := json.Marshal(entry)
	request, err := http.NewRequest("POST", "http://logger-service:8080/write-log", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}
