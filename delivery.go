package broker

import (
  "time"
  "github.com/streadway/amqp"
  "strconv"
)

type Delivery struct {
   ContentType     string    // MIME content type
   ContentEncoding string    // MIME content encoding
   CorrelationId   string    // application use - correlation identifier
   MessageId       string    // application use - message identifier
   Timestamp       time.Time // application use - message timestamp
   Type            string    // application use - message type name
   UserId          int64    // application use - creating user - should be authenticated user
   AppId           string    // application use - creating application id

   Redelivered bool
   Exchange    string // basic.publish exhange
   RoutingKey  string // basic.publish routing key

   Body []byte
   Err error
}

func NewDelivery(d amqp.Delivery) *Delivery {
  var userId int64
  var err error
  if len(d.UserId) > 0 {
    userId, err = strconv.ParseInt(d.UserId, 10, 64)
  }
  return &Delivery{
    UserId: userId,
    Err: err,
    ContentType: d.ContentType,
    ContentEncoding: d.ContentEncoding,
    CorrelationId: d.CorrelationId,
    MessageId: d.MessageId,
    Timestamp: d.Timestamp,
    Type: d.Type,
    AppId: d.AppId,
    Redelivered: d.Redelivered,
    Exchange: d.Exchange,
    RoutingKey: d.RoutingKey,
    Body: d.Body,
  }
}
