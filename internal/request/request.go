package request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/johannesalke/tcptohttp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	Status      requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	ParsingRequestLine requestState = iota
	ParsingHeaders
	ParsingBody
	Done
)

const bufferSize = 8

/*
HTTP-version  = HTTP-name "/" DIGIT "." DIGIT
HTTP-name     = %s"HTTP"
request-line  = method SP request-target SP HTTP-version

GET /coffee HTTP/1.1
*/

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{Status: ParsingRequestLine, Headers: headers.NewHeaders()}
	buf := make([]byte, bufferSize)
	var readToIndex = 0
	for request.Status != Done {
		if readToIndex == len(buf) {
			nbuf := make([]byte, len(buf)*2)
			copy(nbuf, buf)
			buf = nbuf
		}
		nRead, err := reader.Read(buf[readToIndex:])
		if err == io.EOF { //&& readToIndex == 0

			_, err := request.parse(buf[:readToIndex+nRead])
			if err != nil {
				return nil, err
			}
			//request.Status = Done
			if request.Status != Done {
				return nil, fmt.Errorf("Parsing couldn't finish.")
			}
			break
		} else if err != nil && err != io.EOF {
			return nil, err
		}
		if request.Status == Done {
			return request, nil
		}
		readToIndex += nRead
		for bytes.Contains(buf[:readToIndex], []byte("\r\n")) {

			nParsed, err := request.parse(buf[:readToIndex])
			if err != nil {
				return nil, err
			}
			copy(buf, buf[nParsed:]) //:readToIndex
			readToIndex -= nParsed
			//WARNING WARNING WARNING WARNING WARNING
			if request.Status == Done {
				return request, nil
			}
		}
	}

	return request, nil

}

func (r *Request) parse(data []byte) (int, error) {
	if r.Status == Done {
		return 0, fmt.Errorf("The request has already been fully parsed")
	}
	if !bytes.Contains(data, []byte("\r\n")) && r.Status != ParsingBody {
		return 0, nil
	}
	index := bytes.Index(data, []byte("\r\n"))
	if r.Status == ParsingRequestLine {

		requestLine, err := parseRequestLine(string(data[:index]))
		if err != nil {
			return 0, fmt.Errorf("Error parsing request line: %s", err)
		}
		r.RequestLine = *requestLine
		r.Status = ParsingHeaders
		return index + 2, nil
	} else if r.Status == ParsingHeaders {

		n, finished, err := r.Headers.Parse(data)
		if err != nil {
			return 0, fmt.Errorf("Error parsing header: %s", err)
		}
		if finished == true && (r.Headers.Get("content-length") == "" || r.Headers.Get("content-length") == "0") {
			r.Status = Done
		} else if finished == true && r.Headers.Get("content-length") != "" {
			r.Status = ParsingBody
		}

		return n, nil
	} else if r.Status == ParsingBody {

		r.Body = append(r.Body, data...)
		specLength, err := strconv.Atoi(r.Headers.Get("content-length"))
		if err != nil {
			return 0, err
		}
		if specLength == 0 {
			r.Status = Done
			return 0, nil
		}
		if len(r.Body) > specLength {
			return len(data), fmt.Errorf("Body longer than specified by content-length header")
		} else if len(r.Body) == specLength {
			r.Status = Done
		} else if len(r.Body) < specLength && len(data) == 0 {
			return 0, fmt.Errorf("Body shorter than content-length specified by header")
		}
		//fmt.Printf("Body-Length: %d, Content-Length: %d\n", len(r.Body), specLength)
		return len(data), nil
	}

	//fmt.Println("You shouldn't be here")
	return 0, nil
}

func parseRequestLine(reqLine string) (*RequestLine, error) {
	sections := strings.Split(reqLine, " ")
	if len(sections) != 3 {

		return nil, fmt.Errorf("Malformed Request Line: Missing sections")

	}
	method := sections[0]
	path := sections[1]
	version := strings.Split(sections[2], "/")[1]
	if strings.ToUpper(method) != method {
		return nil, fmt.Errorf("Invalid request method format")
	}

	if version != "1.1" {
		return nil, fmt.Errorf("Http version is not 1.1")
	}
	line := RequestLine{HttpVersion: version, Method: method, RequestTarget: path}
	return &line, nil
}
