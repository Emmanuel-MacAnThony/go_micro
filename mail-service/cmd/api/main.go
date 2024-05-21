package main

import (
	"fmt"
	"log"
	"net/http"
)

const WEB_PORT = "80"

type Config struct {
}

func main() {
	app := Config{}
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