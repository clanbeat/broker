package broker

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func Log(d amqp.Delivery) string {
	return fmt.Sprintf("[broker][from %s][%s]: %s", d.Exchange, d.RoutingKey, d.MessageId)
}

func logDelivery(d amqp.Delivery) {
	logString := fmt.Sprintf(
		"[broker][from %s][%s]: %s",
		d.Exchange,
		d.RoutingKey,
		d.MessageId)
	log.Println(logString)
}
