package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Application) authentication(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, errors.New("Invalid request payload"), http.StatusBadRequest)
		return
	}
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("Invalid email or password"), http.StatusUnauthorized)
		return
	}
	// Here you would typically verify the password
	isMatch, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !isMatch {
		app.errorJSON(w, errors.New("Invalid email or password"), http.StatusUnauthorized)
		return
	}
	log.Printf("User %s authenticated successfully", user.Email)
	err = app.LogRequest("authenticate", fmt.Sprintf("Logged in user %s", user.Email))
	if err != nil {
		app.errorJSON(w, errors.New("Failed to log authentication request"), http.StatusInternalServerError)
		return
	}
	payload := JsonResponse{
		Error:   false,
		Message: "Authentication successful",
		Data:    user,
	}

	err = app.writeJSON(w, http.StatusAccepted, payload)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

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
