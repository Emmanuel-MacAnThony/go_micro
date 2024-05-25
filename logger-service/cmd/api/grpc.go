package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Emmanuel-MacAnThony/logger/data"
	"github.com/Emmanuel-MacAnThony/logger/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	res := &logs.LogResponse{Result: "logged"}
	return res, nil
}

func (app *Config) grpcListen(){
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", GRPC_PORT))
	if err != nil{
		log.Fatalf("failed to listen for gRPC: %v", err)

	}

	s := grpc.NewServer()
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.models})

	log.Println("gRPC server started on port %s", GRPC_PORT)

	if err:=s.Serve(lis); err != nil{
		log.Fatalf("failed to listen for gRPC: %v", err)
	}
}