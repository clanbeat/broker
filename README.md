# message_broker

Send and consume persistent messages in json format through topic exchanges

## Installation

Clone the repository to your go workspace

	$ git clone git@github.com:clanbeat/broker.git

## Example

```go
  package main

  import (
    "encoding/json"
  	"github.com/clanbeat/broker"
  )

  func main() {
    //set up broker connection
  	brokerConn, err := broker.New(os.Getenv("RABBITMQ_URI"))
  	if err != nil {
  		log.Fatal(err)
  	}

    //set up channel
  	messageChannel, err := brokerConn.NewChannel()
    if err != nil {
  		log.Fatal(err)
  	}

    //define channel for producing messages
  	if err = messageChannel.ExchangeDeclare(apiExchange); err != nil {
  		log.Fatal(err)
  	}

    message := map[string]int{"from": "someone", "content": "Hello!"}

    if err := messageChannel.Publish("model.event", json.Marshal(message)); err != nil {
      log.Println(err)
    }


  	defer messageChannel.Close()
  	defer brokerConn.Close()
  }
```
