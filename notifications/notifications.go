package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/nikvas0/dc-homework/notifications/email"
	"github.com/streadway/amqp"
)

type Notification struct {
	Email string
	Text  string
}

const connectRetries = 10
const sleepBetweenConnectRetriesDuration = 2 * time.Second

func main() {
	var conn *amqp.Connection
	var err error
	counter := 0
	for {
		conn, err = amqp.Dial(os.Getenv("RABBITMQ"))
		if err != nil {
			counter++
			if counter == connectRetries {
				log.Fatalf("Failed to connect: %s", err)
			}
			log.Printf("Failed to connect to queue: %v. Retrying...", err)
			time.Sleep(sleepBetweenConnectRetriesDuration)
		} else {
			break
		}
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create channel: %s", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"notifications", // name
		"fanout",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf("Failed to create exchange: %s", err)
	}

	q, err := ch.QueueDeclare(
		"emails", // name
		true,     // durable -- в данном примере не требуется
		false,    // delete when unused
		true,     // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("Failed to create queue: %s", err)
	}

	err = ch.QueueBind(
		q.Name, // queue name
		"#",
		"notifications", // exchange
		false,
		nil)

	if err != nil {
		log.Fatalf("Failed to bind queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	for d := range msgs {
		log.Printf("Received a email send task: %s", d.Body)

		notify := Notification{}
		err = json.Unmarshal(d.Body, &notify)
		if err != nil {
			log.Fatalln("SignUp request error: Got broken JSON.")
		}

		email.SendEmail(notify.Email, notify.Text)

		log.Println("Email sent.")

		d.Ack(false)
	}
}
