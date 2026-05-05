package response

import (
	"fmt"
	"io"
	"strconv"

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

type WriterState int

const (
	WritingStatusLine WriterState = iota // 0
	WritingHeaders
	WritingBody
	Done
)

type Writer struct {
	//Conn *net.Conn
	IOWriter    io.Writer
	WriterState WriterState
	BodyLength  int
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	switch statusCode {
	case Success:
		w.IOWriter.Write([]byte("HTTP/1.1 200 OK\r\n"))
		w.WriterState = WritingHeaders
		return nil
	case ClientError:
		w.IOWriter.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		w.WriterState = WritingHeaders
		return nil
	case ServerError:
		w.IOWriter.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		w.WriterState = WritingHeaders
		return nil
	default:
		return fmt.Errorf("Invalid Status code: %s", statusCode)

	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.WriterState != WritingHeaders {
		return fmt.Errorf("Status line must be written before headers")
	}
	for key, value := range headers {
		line := fmt.Sprintf("%s: %s%s", key, value, crlf)
		_, err := w.IOWriter.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	_, err := w.IOWriter.Write([]byte(crlf))
	if err != nil {
		return err
	}
	if (headers.Get("content-length") == "0" || headers.Get("content-length") == "") && headers.Get("Transfer-Encoding") != "chunked" {
		w.WriterState = Done
		return nil
	}
	w.WriterState = WritingBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.WriterState != WritingBody {
		return 0, fmt.Errorf("The body must be written last, and only if the content-length header has the appropriate size.")
	}
	w.BodyLength += len(p)
	w.WriterState = Done
	return w.IOWriter.Write(p)

}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	length := int64(len(p))
	lengthStr := strconv.FormatInt(length, 16)
	chunk := []byte(lengthStr + crlf + string(p) + crlf)
	w.BodyLength += len(p)
	return w.IOWriter.Write(chunk)
}

func (w *Writer) WriteChunkedBodyDone(trailers bool) (int, error) {
	var endline []byte
	if trailers {
		endline = []byte("0" + crlf)
	} else {
		endline = []byte("0" + crlf + crlf)
	}
	w.WriterState = Done
	return w.IOWriter.Write(endline)
}

func (w *Writer) WriteTrailers(trailers headers.Headers) error {
	for key, value := range trailers {
		line := fmt.Sprintf("%s: %s%s", key, value, crlf)
		_, err := w.IOWriter.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	_, err := w.IOWriter.Write([]byte(crlf))
	if err != nil {
		return err
	}
	return nil
}
