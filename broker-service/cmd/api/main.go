package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const WEB_PORT = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	rabbitConn, err := connect()
	if err != nil{
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitConn.Close()
	log.Println("Listening to RabbitMQMessages 🚀🚀🚀")

	log.Printf("starting broker service on port %s\n",WEB_PORT)

	app := Config{
		Rabbit: rabbitConn,
	}

	srv := http.Server{
		Addr: fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	if err != nil{
		log.Panic(err)
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