package queues

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nikvas0/dc-homework/product_upload/objects"
	"github.com/streadway/amqp"
)

const connectRetries = 10
const sleepBetweenConnectRetriesDuration = 2 * time.Second

var connection *amqp.Connection
var uploadChannel *amqp.Channel

func Init(connectionStr string) error {
	var conn *amqp.Connection
	var err error
	counter := 0
	for {
		conn, err = amqp.Dial(connectionStr)
		if err != nil {
			counter++
			if counter == connectRetries {
				return err
			}
			log.Printf("Failed to connect to queue: %v. Retrying...", err)
			time.Sleep(sleepBetweenConnectRetriesDuration)
		} else {
			break
		}
	}
	log.Println("Connected to queue.")

	connection = conn

	uploadChannel, err = connection.Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		return err
	}

	err = uploadChannel.ExchangeDeclare(
		"product_upload", // name
		"fanout",         // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Println("Failed to create exchange")
		return err
	}
	return nil
}

func Close() {
	uploadChannel.Close()
	connection.Close()
}

func SheduleProductBatch(products []objects.Product) error {
	body, err := json.Marshal(products)
	if err != nil {
		log.Println("json marshal failed")
		return err
	}

	log.Println("Schedule product batch.")

	err = uploadChannel.Publish(
		"product_upload", // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	return nil
}
