package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

const webPort = "8080"

type Application struct {
	Mailer Mail
}

func main() {
	app := Application{
		Mailer: createMail(),
	}
	log.Printf("Starting mailer service on port %s", webPort)

	server := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		FromName:    os.Getenv("MAIL_FROMNAME"),
		FromAddress: os.Getenv("MAIL_FROMADDRESS"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
	}
	return m
}
