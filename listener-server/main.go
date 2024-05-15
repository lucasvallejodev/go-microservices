package main

import (
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitConn, err := connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var conn *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://rabbitmq:password@localhost:5672")
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v", err)
			counts++
		} else {
			conn = c
			break
		}

		if counts > 5 {
			log.Printf("Failed to connect to RabbitMQ after %d retries\n", counts)
			return nil, err
		}

		backOff = time.Duration(math.Pow(2, float64(counts))) * time.Second
		time.Sleep(backOff)
		continue
	}

	return conn, nil
}
