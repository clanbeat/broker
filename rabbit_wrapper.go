package broker

import (
	"github.com/streadway/amqp"
)

type (
	Connection struct {
		amqpConnection *amqp.Connection
	}

	Channel struct {
		amqpChannel *amqp.Channel
		Exchange string
	}
)

func New(rabbitmqURI string) (*Connection, error) {
	c := &Connection{}

	conn, err := amqp.Dial(rabbitmqURI)
	if err != nil {
		return c, err
	}
	c.amqpConnection = conn
	return c, nil
}

func (conn *Connection) NewChannel() (*Channel, error) {
	ch := &Channel{}
	createdChan, err := conn.amqpConnection.Channel()
	if err != nil {
		return ch, err
	}
	ch.amqpChannel = createdChan

	return ch, nil
}

func (ch *Channel) ExchangeDeclare(exchangeName string) error {
	err := ch.amqpChannel.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		return err
	}
	ch.Exchange = exchangeName
	return nil
}

func (ch *Channel) Close() {
	ch.amqpChannel.Close()
}

func (conn *Connection) Close() {
	conn.amqpConnection.Close()
}
