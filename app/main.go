package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/howtri/gLogger/app/protocol"
	"github.com/howtri/gLogger/app/server"

	"google.golang.org/grpc"
)

func main() {
	pPort := flag.Int("Port", 8080, "Port for the logger to listen on")
	pLogFileName := flag.String("Log File", "gLog.txt", "Log file to write and read log messages to/from")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *pPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	gLogger := server.NewgLogger(pLogFileName)
	protocol.RegisterLoggerServer(s, gLogger)

	log.Println("gLogger listening")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Serve failed.")
	}
}
