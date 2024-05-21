package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const WEB_PORT = "80"

type Config struct {
	mailer Mail
}

func main() {
	app := Config{
		mailer: createMail(),
	}
	log.Println("starting mail service on port", WEB_PORT, "🚀🚀🚀")
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil{
		log.Panic(err)
	}
}

func createMail() Mail{
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain: os.Getenv("MAIL_DOMAIN"),
		Host:  os.Getenv("MAIL_HOST"),
		Port: port,
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
		Encryption: os.Getenv("MAIL_ENCRYPTION"),
		FromName: os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),

	}

	return m
}