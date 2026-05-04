package response

import (
	"fmt"
	"io"

	"github.com/johannesalke/tcptohttp/internal/headers"
)

type StatusCode int

const (
	Success     StatusCode = iota // 200
	ClientError                   // 400
	ServerError                   // 500
)

const crlf = "\r\n"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case Success:
		w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return nil
	case ClientError:
		w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return nil
	case ServerError:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return nil
	default:
		return fmt.Errorf("Invalid Status code: %s", statusCode)

	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		line := fmt.Sprintf("%s: %s%s", key, value, crlf)
		_, err := w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte(crlf))
	if err != nil {
		return err
	}
	return nil
}
