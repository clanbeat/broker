package broker

import (
	"github.com/streadway/amqp"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

func (ch *Channel) Publish(routingKey string, body []byte) error {
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
