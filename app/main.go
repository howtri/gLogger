package main

import (
	"log"
	"net"

	"github.com/howtri/gLogger/app/protocol"
	"github.com/howtri/gLogger/app/server"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	gLogger := server.NewgLogger()
	protocol.RegisterLoggerServer(s, gLogger)

	log.Println("gLogger listening")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Serve failed.")
	}
}
