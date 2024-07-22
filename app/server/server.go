package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/howtri/gLogger/app/protocol"
)

type gLogger struct {
	protocol.UnimplementedLoggerServer
	writeLock sync.Mutex
}

func NewgLogger() *gLogger {
	return &gLogger{}
}

func (l *gLogger) Write(ctx context.Context, msg *protocol.LogMessage) (*protocol.WriteResponse, error) {
	// TODO: Wrap in goroutine as an anonymous function
	// TODO: Additional logic
	fmt.Printf("Log message to write is %s!\n", msg.GetMessage())
	return &protocol.WriteResponse{StatusCode: 500}, nil
}

func (l *gLogger) Read(ctx context.Context, msg *protocol.ReadRequest) (*protocol.ReadResponse, error) {
	// TODO: Wrap in goroutine as an anonymous function
	// TODO: Additional logic
	fmt.Printf("Read request is %s!\n", msg.GetService())

	message := protocol.LogMessage{Message: "Hi there you beautiful message"}
	var messages []*protocol.LogMessage
	messages = append(messages, &message)
	return &protocol.ReadResponse{StatusCode: 500, Logs: messages}, nil
}
