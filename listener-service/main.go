package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/Emmanuel-MacAnThony/listener/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitConn, err := connect()
	if err != nil{
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitConn.Close()
	log.Println("Listening to RabbitMQMessages 🚀🚀🚀")

	consumer, err := events.NewConsumer(rabbitConn)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARN", "log.ERROR"})
	if err != nil {
		log.Println(err)
		
	}



}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff  = 1*time.Second 
	var connection *amqp.Connection

	for {
		c,err := amqp.Dial("amqp://guest:guest@rabbitmq:5672")
		if err != nil{
			fmt.Println("rabbitmq not ready")
			counts++
		}else{
			connection = c
			log.Println("Connected to RabbitMQ 🚀🚀🚀")
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}
		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off:::")
		time.Sleep(backoff)
	}

	return connection, nil

}