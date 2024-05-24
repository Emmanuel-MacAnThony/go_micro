package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/Emmanuel-MacAnThony/logger/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	WEB_PORT  = "80"
	RPC_PORT  = "5001"
	MONGO_URL = "mongodb://mongo_logger:27017"
	GRPC_PORT = "5001"
)

type Config struct {
	models data.Models
}

var client *mongo.Client

func main() {

	mongoClient, err :=connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func(){
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{data.New(client)}

	go rpc.Register(new(RPCServer))
	go app.rpcListen()

	log.Println("starting http server 🚀🚀🚀🚀")
	srv := http.Server{
		Addr: fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.routes(),

	}
	err = srv.ListenAndServe()
	
	if err != nil {
		log.Panic(err)
	}
	
}



func (app *Config) rpcListen() error{
	log.Println("Starting RPC server on port::", RPC_PORT)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", RPC_PORT))
	if err != nil{
		return err
	}
	defer listen.Close()
	for {
		rpcConn, err := listen.Accept()
		if err != nil{
			continue
		}

		go rpc.ServeConn(rpcConn)
	}


}

func connectToMongo() (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(MONGO_URL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil{
		log.Println("Error connecting:", err)
		return nil, err
	}
	log.Println("Connected to DB")

	return c, nil

}