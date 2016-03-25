package broker

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type (
	ConnectionInfo struct {
		ExchangeName string
		QueueName    string
		RoutingKey   string
	}

	Connection struct {
		rabbitURI      string
		amqpConnection *amqp.Connection
		sendError      ErrorTracker
		channel        *Channel
		ExistingRetry  bool
	}

	Channel struct {
		amqpChannel *amqp.Channel
		consumers   []*Consumer
	}

	ErrorTracker func(err error)
)

func New(rabbitmqURI string, sendError ErrorTracker) (*Connection, error) {
	if len(rabbitmqURI) == 0 {
		return nil, errors.New("rabbitmqURI missing")
	}

	c := &Connection{sendError: sendError, rabbitURI: rabbitmqURI}
	if err := c.connect(); err != nil {
		return c, err
	}
	ch, err := c.newChannel()
	if err != nil {
		return c, err
	}
	c.channel = ch
	return c, nil
}

func (conn *Connection) newChannel() (*Channel, error) {
	ch := &Channel{}
	if conn.amqpConnection == nil {
		return ch, errors.New("rabbitmq connection missing")
	}
	if err := ch.connect(conn); err != nil {
		return ch, err
	}
	return ch, nil
}

func (c *Connection) connect() error {
	c.Close()
	conn, err := amqp.Dial(c.rabbitURI)
	if err != nil {
		return err
	}
	c.amqpConnection = conn
	go c.handleFailures(c.amqpConnection.NotifyClose(make(chan *amqp.Error)))
	return nil
}

func (ch *Channel) connect(conn *Connection) error {
	ch.Close()
	createdChan, err := conn.amqpConnection.Channel()
	if err != nil {
		return err
	}
	ch.amqpChannel = createdChan
	go conn.handleFailures(ch.amqpChannel.NotifyClose(make(chan *amqp.Error)))
	return nil
}

func (c *Connection) ExchangeDeclare(exchangeName string) error {
	if c.channel == nil {
		return errors.New("no message channel defined")
	}
	return c.channel.exchangeDeclare(exchangeName)
}

func (ch *Channel) exchangeDeclare(exchangeName string) error {
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
	return nil
}

func (c *Connection) handleFailures(errs chan *amqp.Error) {
	for e := range errs {
		if c.ExistingRetry == true {
			return
		}
		c.ExistingRetry = true
		log.Println(fmt.Sprintf("[broker][connection_error]: %d %s", e.Code, e.Reason, "server side:", e.Server, "recoverable:", e.Recover))
		c.Reset()
	}
}

func (c *Connection) Reset() {
	resetInSeconds := 1 * time.Second

loop:
	for {
		log.Println("[broker][retry_connection] in", resetInSeconds, "seconds")
		if err := c.reConnect(); err != nil {
			log.Println("[broker][retry_error]", err)
		} else {
			break loop
		}
		time.Sleep(resetInSeconds)
		if resetInSeconds >= 5*time.Second {
			c.sendError(errors.New("couldn't connect to rabbitmq in 3 tries"))
			resetInSeconds = 1 * time.Second
		} else {
			resetInSeconds += (2 * time.Second)
		}
	}
	log.Println("[broker][reconnected_successfully]")
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

func (c *Connection) reConnect() error {
	if err := c.connect(); err != nil {
		return err
	}
	if err := c.channel.connect(c); err != nil {
		return err
	}
	if err := c.channel.restartConsumers(); err != nil {
		return err
	}
	c.ExistingRetry = false
	return nil
}
