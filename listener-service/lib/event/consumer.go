package event

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// returns new consumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}
	return consumer, nil
}

// setup opens a channel and declare the exchange
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

// type used for pushing events to AMQP
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Listens for all new queue publications
func (consumer *Consumer) Listen(topic []string) error {
	cha, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer cha.Close()

	q, err := declareRandomQueue(cha)
	if err != nil {
		return err
	}

	for _, s := range topic {
		err = cha.QueueBind(q.Name, s, getExchangeName(), false, nil)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	messages, err := cha.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			handlePayload(payload)
		}
	}()
	log.Printf("[*] waiting for message [Exchange, Queue][%s,%s].", getExchangeName(), q.Name)
	<-forever
	return nil
}

// takes an action based on the name of the event is the queue
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	request, err := http.NewRequest("POST", "http://logger-service:8080/write-log", bytes.NewBuffer(jsonData))
	if err != nil{
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		return err
	}
	return nil
}
