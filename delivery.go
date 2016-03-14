package broker

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"time"
)

type Delivery struct {
	ContentType     string    // MIME content type
	ContentEncoding string    // MIME content encoding
	CorrelationId   string    // application use - correlation identifier
	MessageId       string    // application use - message identifier
	Timestamp       time.Time // application use - message timestamp
	Type            string    // application use - message type name
	UserId          int64     // application use - creating user - should be authenticated user
	AppId           string    // application use - creating application id

	Redelivered bool
	Exchange    string // basic.publish exhange
	RoutingKey  string // basic.publish routing key

	Body             []byte
	originalDelivery amqp.Delivery
}

func NewDelivery(d amqp.Delivery) *Delivery {
	var userId int64
	data := d.Body

	if d.Type == messageTypeV2 {
		var b Body
		if err := json.Unmarshal(d.Body, &b); err == nil {
			if b.UserID > 0 {
				userId = b.UserID
			}
			data = b.Data
		}
	}

	return &Delivery{
		originalDelivery: d,
		UserId:           userId,
		ContentType:      d.ContentType,
		ContentEncoding:  d.ContentEncoding,
		CorrelationId:    d.CorrelationId,
		MessageId:        d.MessageId,
		Timestamp:        d.Timestamp,
		Type:             d.Type,
		AppId:            d.AppId,
		Redelivered:      d.Redelivered,
		Exchange:         d.Exchange,
		RoutingKey:       d.RoutingKey,
		Body:             data,
	}
}

func (d Delivery) Ack(multiple bool) error {
	return d.originalDelivery.Ack(multiple)
}

func (d Delivery) Nack(multiple, requeue bool) error {
	return d.originalDelivery.Nack(multiple, requeue)
}

func (d Delivery) Reject(requeue bool) error {
	return d.originalDelivery.Reject(requeue)
}
