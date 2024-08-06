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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"net/http"
)

// Define Prometheus metrics
var (
	grpcRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method"},
	)
	grpcRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Duration of gRPC requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
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

func metricsInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	grpcRequests.WithLabelValues(info.FullMethod).Inc()
	timer := prometheus.NewTimer(grpcRequestDuration.WithLabelValues(info.FullMethod))
	defer timer.ObserveDuration()

	return handler(ctx, req)
}

func main() {
	pPort := flag.Int("Port", 8081, "Port for the logger to listen on")
	pLogFileName := flag.String("Log File", "gLog.txt", "Log file to write and read log messages to/from")
	flag.Parse()

	// Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe("0.0.0.0:2112", nil)
	}()

	// gRPC endpoint
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *pPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingInterceptor,
			metricsInterceptor,
		),
	)
	gLogger := server.NewgLogger(pLogFileName)
	protocol.RegisterLoggerServer(s, gLogger)

	log.Println("gLogger listening")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Serve failed.")
	}
}
