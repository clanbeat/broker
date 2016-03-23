package broker

import (
	"errors"
	"github.com/streadway/amqp"
)

type Consumer struct {
	Info     *ConnectionInfo
	Callback func(d Delivery)
}

func (conn *Connection) Register(c *Consumer) error {
	ch := conn.channel
	ch.consumers = append(ch.consumers, c)
	return ch.startConsumer(c)
}

func (ch *Channel) restartConsumers() error {
	for _, cons := range ch.consumers {
		if err := ch.startConsumer(cons); err != nil {
			return err
		}
	}
	return nil
}

func (ch *Channel) startConsumer(c *Consumer) error {
	info := c.Info
	q, err := ch.BindQueue(info.ExchangeName, info.QueueName, info.RoutingKey)
	if err != nil {
		return err
	}

	msgs, err := ch.ConsumeMessages(q.Name)
	if err != nil {
		return err
	}

	go func() {
		for m := range msgs {
			c.Callback(m)
			m.Ack(false)
		}
	}()
	return nil
}

func (ch *Channel) BindQueue(exchangeName, queueName, routingKey string) (amqp.Queue, error) {
	var err error
	var q amqp.Queue

	if ch.amqpChannel == nil {
		return q, errors.New("rabbitmq connection missing")
	}

	err = ch.ExchangeDeclare(exchangeName)
	if err != nil {
		return q, err
	}

	q, err = ch.amqpChannel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return q, err
	}

	err = ch.amqpChannel.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	return q, err
}

func (ch *Channel) Consume(queueName string) (<-chan amqp.Delivery, error) {
	if ch.amqpChannel == nil {
		return nil, errors.New("rabbitmq connection missing")
	}
	return ch.startConsume(queueName)
}

func (ch *Channel) ConsumeMessages(queueName string) (<-chan Delivery, error) {
	if ch.amqpChannel == nil {
		return nil, errors.New("rabbitmq connection missing")
	}

	messages, err := ch.startConsume(queueName)
	if err != nil {
		return nil, err
	}

	deliveries := make(chan Delivery)
	go func() {
		for m := range messages {
			logDelivery(m)
			deliveries <- *NewDelivery(m)
		}
		close(deliveries)
	}()

	return (<-chan Delivery)(deliveries), nil
}

func (ch *Channel) startConsume(queueName string) (<-chan amqp.Delivery, error) {
	return ch.amqpChannel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)
}
