package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	conn *amqp.Connection
}

func (emitter *Emitter) setup() error {
	channel, err := emitter.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	return declareExchange(channel)
}

func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	err = channel.Publish(
		"logs_topic", // exchange
		severity,     // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	if err != nil {
		return err
	}
	return nil
}

func NewEmitter(conn *amqp.Connection) (*Emitter, error) {
	emitter := &Emitter{
		conn: conn,
	}
	err := emitter.setup()
	if err != nil {
		return nil, err
	}
	return emitter, nil
}
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		conn: conn,
	}
	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
