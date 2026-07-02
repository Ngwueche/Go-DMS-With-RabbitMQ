package main

import (
	"listener/lib/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	//try connect to rabbitmq
	log.Printf("Starting the listener service")

	conn, err := connectToRabbitMQ()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	//start listening to messages #2
	log.Println("Listening and consuming rabbitMQ messages")
	//create consumers
	consumer, err := event.NewConsumer(conn)
	if err != nil {
		panic(err)
	}
	//watch the queue and consume events
	err = consumer.Listen([]string{"user.created", "user.updated", "user.deleted","log.INFO", "log.ERROR", "log.WARNING"})
	if err != nil {
		log.Println(err)
		return
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

func connectToRabbitMQ2() (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	maxRetries := 10

	for i := 1; i <= maxRetries; i++ {
		conn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
		if err == nil {
			log.Printf("COnnected to RabbitMQ")
			return conn, nil
		}
		log.Printf("RabbitMQ not ready yet (attempt %d/%): %v", i, maxRetries, err)

		backOff := time.Duration(i*i) * time.Second
		log.Printf("Backing off for %v seconds ...", backOff)

		time.Sleep(backOff)
	}
	return nil, err
}
