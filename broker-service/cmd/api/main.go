package main

import (
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"

type Application struct {
	RabbitConn *amqp.Connection
}

func main() {

	conn, err := connectToRabbitMQ()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	app := Application{
		RabbitConn: conn,
	}
	
	log.Printf("Starting broker service on port %s", webPort)

	server := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func connectToRabbitMQ() (*amqp.Connection, error) {
	var count int
	var backoff = 2 * time.Second
	var connection *amqp.Connection
	rabbitMQConnectionString := os.Getenv("RABBITMQ_URL")

	for {
		c, err := amqp.Dial(rabbitMQConnectionString)
		if err != nil {
			log.Printf("RabbitMQ not ready yet ...")
			count++
		} else {
			log.Printf("Connected to RabbitMQ")
			connection = c
			break
		}

		if count > 10 {
			log.Printf("Could not connect to RabbitMQ, exiting ...")
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(count), 2)) * time.Second
		log.Printf("Backing off (rabbitMQ) ...")
		time.Sleep(backoff)
		continue
	}
	return connection, nil
}
