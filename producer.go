package broker

import (
	"encoding/json"
	"errors"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"time"
)

type Body struct {
	Data   json.RawMessage `json:"data"`
	UserID int64           `json:"userId"`
}

func (ch *Channel) Publish(routingKey string, body []byte) error {
	b, err := json.Marshal(&Body{Data: body})
	if err != nil {
		return err
	}
	return ch.publishBody(routingKey, b)
}

func (ch *Channel) PublishWithUser(routingKey string, userID int64, body []byte) error {
	b, err := json.Marshal(&Body{UserID: userID, Data: body})
	if err != nil {
		return err
	}
	return ch.publishBody(routingKey, b)
}

func (ch *Channel) publishBody(routingKey string, body []byte) error {
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
