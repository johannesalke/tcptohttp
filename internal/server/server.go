package server

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"sync/atomic"

	"github.com/johannesalke/tcptohttp/internal/request"
	"github.com/johannesalke/tcptohttp/internal/response"
)

type Server struct {
	Listener net.Listener
	Active   atomic.Bool
	Handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		Listener: listener,
		Active:   atomic.Bool{},
		Handler:  handler,
	}

	s.Active.Store(true)
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.Active.Store(false)
	if err := s.Listener.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Server) listen() {
	for s.Active.Load() {
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.Active.Load() == false {
				return
			}
			fmt.Printf("Error accepting connection: %s", err)
		}

		fmt.Println("Calling handle!")
		go s.handle(conn)

	}

}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)

	if err != nil {
		fmt.Println("Error getting request from connection:", err)
	}
	var buf bytes.Buffer
	var w = &response.Writer{IOWriter: conn, WriterState: response.WritingStatusLine, BodyLength: 0}
	handlerError := s.Handler(w, req)
	if handlerError != nil {

		err := writeError(conn, *handlerError)
		if err != nil {
			fmt.Print(err)
			return
		}
		return
	}

	headers := response.GetDefaultHeaders(buf.Len())
	err = response.WriteStatusLine(conn, response.Success)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	err = response.WriteHeaders(conn, headers)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	_, err = buf.WriteTo(conn)
	if err != nil && err != io.EOF {
		fmt.Printf("Error: %s\n", err)
		return
	}

}

func writeError(w io.Writer, handlerError HandlerError) error {
	err := response.WriteStatusLine(w, handlerError.StatusCode)
	if err != nil {
		return fmt.Errorf("Error writing error to connection: %s\n", err)
	}
	defHeaders := response.GetDefaultHeaders(len(handlerError.Message))
	if handlerError.ContentType != "" {
		defHeaders.Set("Content-Type", handlerError.ContentType)
	}
	err = response.WriteHeaders(w, defHeaders)
	if err != nil {
		return fmt.Errorf("Error writing error headers to connection: %s\n", err)
	}
	w.Write([]byte(handlerError.Message))
	return nil
}

type Handler func(w *response.Writer, req *request.Request) *HandlerError

// The official http library type for reference:
// type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type HandlerError struct {
	StatusCode  response.StatusCode
	ContentType string
	Message     string
}
