package server

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/howtri/gLogger/app/protocol"
)

type gLogger struct {
	protocol.UnimplementedLoggerServer
	writeLock sync.Mutex
	fileName  string
}

func NewgLogger(pFileName *string) *gLogger {
	return &gLogger{fileName: *pFileName}
}

func (l *gLogger) Write(ctx context.Context, msg *protocol.LogMessage) (*protocol.WriteResponse, error) {
	// TODO: Wrap in goroutine as an anonymous functiona
	log := msg.GetMessage()

	file, err := os.OpenFile(l.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	_, err = gzipWriter.Write([]byte(log))
	if err != nil {
		fmt.Println("Error writing to gzip writer:", err)
		return nil, err
	}
	return &protocol.WriteResponse{StatusCode: 1}, nil
}

func (l *gLogger) Read(ctx context.Context, msg *protocol.ReadRequest) (*protocol.ReadResponse, error) {
	// TODO: Wrap in goroutine as an anonymous function
	fmt.Printf("Read request is %s!\n", msg.GetService())
	fmt.Println("Log we are reading from is %s", l.fileName)

	file, err := os.Open(l.fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	scanner := bufio.NewScanner(gzipReader)

	var messages []*protocol.LogMessage
	for scanner.Scan() {
		message := protocol.LogMessage{Message: scanner.Text()}
		messages = append(messages, &message)
	}
	return &protocol.ReadResponse{StatusCode: 1, Logs: messages}, nil
}
