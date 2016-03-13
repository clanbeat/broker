package broker

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"time"
)

type Delivery struct {
	ContentType     string    // MIME content type
	ContentEncoding string    // MIME content encoding
	CorrelationID   string    // application use - correlation identifier
	MessageID       string    // application use - message identifier
	Timestamp       time.Time // application use - message timestamp
	Type            string    // application use - message type name
	UserID          int64     // application use - creating user - should be authenticated user
	AppID           string    // application use - creating application id

	Redelivered bool
	Exchange    string // basic.publish exhange
	RoutingKey  string // basic.publish routing key

	Body             []byte
	originalDelivery amqp.Delivery
}

func NewDelivery(d amqp.Delivery) *Delivery {
	var userId int64
	var b Body

	data := d.Body
	if err := json.Unmarshal(d.Body, &b); err == nil {
		if b.UserID > 0 {
			userId = b.UserID
		}
		data = b.Data
	}

	return &Delivery{
		originalDelivery: d,
		UserID:           userId,
		ContentType:      d.ContentType,
		ContentEncoding:  d.ContentEncoding,
		CorrelationID:    d.CorrelationId,
		MessageID:        d.MessageId,
		Timestamp:        d.Timestamp,
		Type:             d.Type,
		AppID:            d.AppId,
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
