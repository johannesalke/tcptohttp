package main

import (
	//"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"

	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/johannesalke/tcptohttp/internal/headers"
	"github.com/johannesalke/tcptohttp/internal/request"
	"github.com/johannesalke/tcptohttp/internal/response"
	"github.com/johannesalke/tcptohttp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) *server.HandlerError {
	target := req.RequestLine.RequestTarget
	switch {
	case target == "/yourproblem":
		return &server.HandlerError{
			StatusCode:  response.ClientError,
			ContentType: "text/html",
			Message:     badRequestHtml,
		}
	case target == "/myproblem":
		return &server.HandlerError{
			StatusCode:  response.ServerError,
			ContentType: "text/html",
			Message:     serverErrorHtml,
		}
	case strings.HasPrefix(target, "/httpbin"):
		err := w.WriteStatusLine(response.Success)
		if err != nil {
			fmt.Print(err)
		}
		hdrs := response.GetDefaultHeaders(len([]byte(successHtml)))
		hdrs.Set("Content-Type", "application/json")
		hdrs.Delete("Content-Length")
		hdrs.Set("Transfer-Encoding", "chunked")
		hdrs.Add("Trailer", "X-Content-SHA256")
		hdrs.Add("Trailer", "X-Content-Length")

		err = w.WriteHeaders(hdrs)
		if err != nil {
			fmt.Print(err)
		}
		trgt := strings.TrimPrefix(target, "/httpbin")
		url := "https://httpbin.org" + trgt
		res, err := http.Get(url)
		if err != nil {
			fmt.Print(err)
		}
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Print(err)
		}
		//var chunkBuf = make([]byte, 256)
		/*fileLines := bytes.Split(bodyBytes, []byte("\n"))
		for _, line := range fileLines {
			_, err = w.WriteChunkedBody(line)
			if err != nil {
				fmt.Print(err)
			}
		}*/
		_, err = w.WriteChunkedBody(bodyBytes)
		if err != nil {
			fmt.Print(err)
		}
		_, err = w.WriteChunkedBodyDone(true)
		if err != nil {
			fmt.Print(err)
		}
		checksum := sha256.Sum256(bodyBytes)

		trailers := headers.NewHeaders()

		trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", checksum))
		trailers.Set("X-Content-Length", fmt.Sprintf("%d", w.BodyLength))
		err = w.WriteTrailers(trailers)
		if err != nil {
			fmt.Print(err)
		}
		return nil

	default:
		err := w.WriteStatusLine(response.Success)
		if err != nil {
			fmt.Printf("Error writing status line:%s\n", err)
		}
		hdrs := response.GetDefaultHeaders(len([]byte(successHtml)))
		hdrs.Set("Content-Type", "text/html")
		err = w.WriteHeaders(hdrs)
		if err != nil {
			fmt.Printf("Error writing headers:%s\n", err)
		}
		_, err = w.WriteBody([]byte(successHtml))
		if err != nil {
			fmt.Printf("Error writing body:%s\n", err)
		}
		return nil

	}

}

const badRequestHtml = `
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`

const serverErrorHtml = `
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`

const successHtml = `
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`
