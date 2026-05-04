package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Status      int //0 = initialized, 1 = done
	Step        requestSection
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestSection int

const (
	ReqLine requestSection = iota
	Headers
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
	request := &Request{}
	buf := make([]byte, bufferSize)
	var readToIndex = 0
	for request.Status != 1 {
		if readToIndex == len(buf) {
			nbuf := make([]byte, len(buf)*2)
			copy(nbuf, buf)
			buf = nbuf
		}
		nRead, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			request.Status = 1
			break
		} else if err != nil && err != io.EOF {
			return nil, err
		}

		readToIndex += nRead
		nParsed, err := request.parse(buf)
		if err != nil {
			return nil, err
		}
		copy(buf, buf[nParsed:readToIndex])
		readToIndex -= nParsed
		//WARNING WARNING WARNING WARNING WARNING

	}

	return request, nil

}

func (r *Request) parse(data []byte) (int, error) {
	if r.Status == 1 {
		return 0, fmt.Errorf("The request has already been fully parsed")
	}
	if !bytes.Contains(data, []byte("\r\n")) {
		return 0, nil
	}
	index := bytes.Index(data, []byte("\r\n"))
	if r.Step == ReqLine {
		requestLine, err := parseRequestLine(string(data[:index]))
		if err != nil {
			return 0, fmt.Errorf("Error parsing request line: %e", err)
		}
		r.RequestLine = *requestLine
		r.Status = 1
		return index + 2, nil
	}
	fmt.Println("You shouldn't be here")
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

type validMethods int

const (
	GET validMethods = iota
	POST
	PUT
	DELETE
)
