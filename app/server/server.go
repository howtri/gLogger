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
	responseChan := make(chan *protocol.WriteResponse)
	errorChan := make(chan error)
	defer close(responseChan)
	defer close(errorChan)

	go l.ThreadWriteLog(msg, responseChan, errorChan)

	select {
	case response := <-responseChan:
		return response, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (l *gLogger) ThreadWriteLog(msg *protocol.LogMessage, responseChan chan<- *protocol.WriteResponse, errorChan chan<- error) {
	log := msg.GetMessage()

	file, err := os.OpenFile(l.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errorChan <- fmt.Errorf("Failed to open log file: %w", err)
		return
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	_, err = gzipWriter.Write([]byte(log))
	if err != nil {
		errorChan <- fmt.Errorf("Failed to write to gzip writer: %w", err)
		return
	}

	responseChan <- &protocol.WriteResponse{StatusCode: 1}
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
