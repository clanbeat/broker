package broker

import (
	"errors"
	"github.com/streadway/amqp"
	"log"
)

type (
	ErrorTracker func(err error)
	Connection   struct {
		amqpConnection *amqp.Connection
		sendError      ErrorTracker
	}

	Channel struct {
		amqpChannel *amqp.Channel
		Exchange    string
	}
)

func New(rabbitmqURI string, sendError ErrorTracker) (*Connection, error) {
	c := &Connection{}
	if len(rabbitmqURI) == 0 {
		return c, errors.New("rabbitmqURI missing")
	}
	conn, err := amqp.Dial(rabbitmqURI)
	if err != nil {
		return c, err
	}
	c.amqpConnection = conn
	c.sendError = sendError

	errs := c.amqpConnection.NotifyClose(make(chan *amqp.Error))
	go handleFailures(errs, c.sendError)

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

	errs := ch.amqpChannel.NotifyClose(make(chan *amqp.Error))
	go handleFailures(errs, conn.sendError)
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

func handleFailures(errs chan *amqp.Error, sendError ErrorTracker) {
	for e := range errs {
		log.Printf("connection error: %d %s", e.Code, e.Reason)
		sendError(errors.New(e.Error()))
	}
}
