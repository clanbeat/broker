package broker

import (
	"github.com/streadway/amqp"
	"github.com/satori/go.uuid"
	"time"
	"errors"
	"strconv"
)

func (ch *Channel) Publish(routingKey string, body []byte) error {
	if ch.amqpChannel == nil {
		return errors.New("rabbitmq connection missing")
	}
	return ch.amqpChannel.Publish(
		ch.Exchange, // exchange
		routingKey,  // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.NewV4().String(),
			Timestamp:    time.Now(),
		},
	)
}

func (ch *Channel) PublishWithUser(routingKey string, userID int64, body []byte) error {
	if ch.amqpChannel == nil {
		return errors.New("rabbitmq connection missing")
	}
	return ch.amqpChannel.Publish(
		ch.Exchange, // exchange
		routingKey,  // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.NewV4().String(),
			Timestamp:    time.Now(),
			UserId: strconv.FormatInt(userID, 10),
		},
	)
}
