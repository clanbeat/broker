package broker

import (
	"github.com/streadway/amqp"
	"errors"
)

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
