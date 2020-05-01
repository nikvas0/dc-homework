package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"product_upload_inserter/objects"
	"product_upload_inserter/storage"

	"github.com/streadway/amqp"
)

const connectRetries = 10
const sleepBetweenConnectRetriesDuration = 2 * time.Second

func main() {
	err := storage.Init(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
			os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SSLMODE")))

	if err != nil {
		log.Panicf("Error: %v", err)
	}
	defer storage.Close()

	var conn *amqp.Connection
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
		"product_upload", // name
		"fanout",         // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Failed to create exchange: %s", err)
	}

	q, err := ch.QueueDeclare(
		"product_upload", // name
		false,            // durable -- в данном примере не требуется
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Failed to create queue: %s", err)
	}

	err = ch.QueueBind(
		q.Name, // queue name
		"#",
		"product_upload", // exchange
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
		log.Printf("Received product insert task: %s", d.Body)

		var products []objects.Product
		err = json.Unmarshal(d.Body, &products)
		if err != nil {
			log.Fatalln("Got broken JSON.")
		}

		err = storage.InsertProducts(products)
		if err != nil {
			log.Fatalf("Failed to insert products: %v.", err)
		}

		log.Println("Procduts inserted.")

		d.Ack(false)
	}
}
