package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"context"
	"time"

	"github.com/howtri/gLogger/app/protocol"
	"github.com/howtri/gLogger/app/server"

	"google.golang.org/grpc"
)

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	fmt.Printf("Request - Method:%s Duration:%s Error:%v\n", info.FullMethod, duration, err)
	return resp, err
}

func main() {
	pPort := flag.Int("Port", 8080, "Port for the logger to listen on")
	pLogFileName := flag.String("Log File", "gLog.txt", "Log file to write and read log messages to/from")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *pPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingInterceptor,
			// metricsInterceptor,
		),
	)
	gLogger := server.NewgLogger(pLogFileName)
	protocol.RegisterLoggerServer(s, gLogger)

	log.Println("gLogger listening")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Serve failed.")
	}
}
