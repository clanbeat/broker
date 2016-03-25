package broker

import (
	"encoding/json"
	"errors"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"time"
)

const messageTypeV1 = "v1"
const messageTypeV2 = "v2"

type Body struct {
	Data   json.RawMessage `json:"data"`
	UserID int64           `json:"userId"`
}

type Producer struct {
	Info   *ConnectionInfo
	Data   []byte
	UserID int64
}

func (c *Connection) Publish(p *Producer) error {
	if c.channel == nil {
		return errors.New("no message channel defined")
	}
	if p.Info == nil {
		return errors.New("exchange info missing")
	}
	return c.channel.publishBody(*p.Info, &Body{UserID: p.UserID, Data: p.Data})
}

func (ch *Channel) publishBody(info ConnectionInfo, b *Body) error {
	if ch.amqpChannel == nil {
		return errors.New("rabbitmq connection missing")
	}

	messageType := messageTypeV1
	if b.UserID > 0 {
		messageType = messageTypeV2
	}

	body, err := json.Marshal(&b)
	if err != nil {
		return err
	}

	return ch.amqpChannel.Publish(
		info.ExchangeName, // exchange
		info.RoutingKey,   // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.NewV4().String(),
			Timestamp:    time.Now(),
			Type:         messageType,
		},
	)
}
