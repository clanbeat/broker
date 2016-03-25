# message_broker

Send and consume persistent messages in json format through topic exchanges

## Installation

Clone the repository to your go workspace

	$ git clone git@github.com:clanbeat/broker.git

## Producer example

Sending data to a topic exchange.


```go
  package main

  import (
    "encoding/json"
  	"github.com/clanbeat/broker"
  )

  func main() {
    //set up broker connection
		messageChannel, err = broker.New(os.Getenv("CLOUDAMQP_URL"), errorTrackerFunc())
		if err != nil {
			errorTracker.Error(err)
		}

		//declare exchange once if you start pushing messages
		if err := messageChannel.ExchangeDeclare("myExchange"); err != nil {
			errorTracker.Error(err)
		}
		defer messageChannel.Close()


    message := map[string]string{"from": "someone", "content": "Hello!"}

		producer := &broker.Producer{
			Info: &broker.ConnectionInfo{
				ExchangeName: "myExchange",
				RoutingKey:   "model.event.created",
			},
			Data: json.Marshal(message),
		}

    if err := messageChannel.Publish(producer); err != nil {
      log.Println(err)
    }
  }
```

## Consumer example

Consuming data from an exchange

```go
  package main

  import (
    "encoding/json"
  	"github.com/clanbeat/broker"
  )

  func main() {
    //set up broker connection
		messageChannel, err = broker.New(os.Getenv("CLOUDAMQP_URL"), errorTrackerFunc())
		if err != nil {
			errorTracker.Error(err)
		}

		//declare exchange once if you start pushing messages
		if err := messageChannel.ExchangeDeclare("myExchange"); err != nil {
			errorTracker.Error(err)
		}
		defer messageChannel.Close()
		registerBrokerHandler()
  }

	func registerBrokerHandler() {
		consumer := &broker.Consumer{
			Info: &broker.ConnectionInfo{
				ExchangeName: "myExchange",
				QueueName:    "receiving_model_events",
				RoutingKey:   "model.event.#",
			},
			Callback: messageHandler(),
		}
		brokerConn.Register(consumer)
	}


	func messageHandler() func(msg broker.Delivery) {
		return func(msg broker.Delivery) {
			log.Println(msg)
		}
	}

```
