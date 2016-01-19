package broker

import (
	"github.com/streadway/amqp"
	"errors"
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
	if len(rabbitmqURI) == 0 {
		return c, errors.New("rabbitmqURI missing")
	}
	conn, err := amqp.Dial(rabbitmqURI)
	if err != nil {
		return c, err
	}
	c.amqpConnection = conn
	return c, nil
}

func (conn *Connection) NewChannel() (*Channel, error) {
	ch := &Channel{}
	if conn.amqpConnection == nil {
		return ch, errors.New("rabbitmq connection missing")
	}
	createdChan, err := conn.amqpConnection.Channel()
	if err != nil {
		return ch, err
	}
	ch.amqpChannel = createdChan

	return ch, nil
}

func (ch *Channel) ExchangeDeclare(exchangeName string) error {
	if ch.amqpChannel == nil {
		return errors.New("rabbitmq connection missing")
	}
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
	if ch.amqpChannel != nil {
		ch.amqpChannel.Close()
	}
}

func (conn *Connection) Close() {
	if conn.amqpConnection != nil {
		conn.amqpConnection.Close()
	}
}
