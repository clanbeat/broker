package broker

import (
	"fmt"
	"github.com/streadway/amqp"
)

func Log(d *amqp.Delivery) string {
	return fmt.Sprintf("[broker][%s from %s][%s]: %s", d.Timestamp, d.Exchange, d.RoutingKey, d.MessageId)
}
