package queues

import (
	"encoding/json"
	"log"
	"time"

	"auth/objects"

	"github.com/streadway/amqp"
)

const connectRetries = 10
const sleepBetweenConnectRetriesDuration = 2 * time.Second

var connection *amqp.Connection
var notifyChannel *amqp.Channel

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

	notifyChannel, err = connection.Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		return err
	}

	err = notifyChannel.ExchangeDeclare(
		"notifications", // name
		"fanout",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Println("Failed to create exchange")
		return err
	}
	return nil
}

func Close() {
	notifyChannel.Close()
	connection.Close()
}

func SheduleConfirmation(user *objects.User, token string) error {
	notification := objects.Notification{Email: user.Email, Text: "go to http://localhost:8082/v1/confirm/" + token}
	body, err := json.Marshal(notification)
	if err != nil {
		log.Println("json marshal failed")
		return err
	}

	log.Printf("Schedule confirmation: id=%d token=%s", user.ID, token)

	err = notifyChannel.Publish(
		"notifications", // exchange
		"",              // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	return nil
}
