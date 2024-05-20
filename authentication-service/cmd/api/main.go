package main

import (
	"context"

	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Emmanuel-MacAnThony/authentication/data"

	"github.com/jackc/pgx/v5"
)

const WEB_PORT = "80"
var counts int64

type Config struct {
	DB *pgx.Conn
	Models data.Models
}

func main() {
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to postgres")
	}

	app := Config{
		DB: conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.routes(),

	}

	err := srv.ListenAndServe()
	if err != nil{
		log.Panic(err)
	}

}

func openDB(dsn string)(*pgx.Conn, error) {
	// db, err := sql.Open("pgx", dsn)
	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	// err = db.Ping()
	
	return db, nil
}

func connectToDB() *pgx.Conn {
	dsn := os.Getenv("DSN")
	fmt.Println(dsn)
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet")
			counts++
		}else {
			log.Println("Connected to database")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds")
		time.Sleep(2*time.Second)
		continue
	}
}