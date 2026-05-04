package headers

import (
	"bytes"
	"fmt"

	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	/*if bytes.Equal(data[0:2], []byte("\r\n")) { //Triggers when a line is empty, which signals the end of the headers section
		return 2, true, nil
	}*/
	if !bytes.Contains(data, []byte("\r\n")) { //Triggers when there is not enough material for parsing a full line.
		return 0, false, nil
	}
	index := bytes.Index(data, []byte("\r\n"))
	if index == 0 {
		return 2, true, nil
	}

	name, content, err := parseHeader(string(data[:index]))
	if err != nil {
		return 0, false, err
	}
	for _, c := range name {
		if !validateChar(byte(c)) {
			return 0, false, fmt.Errorf("Error: Field-name contains invalid characters")
		}
	}
	currentContents := h[strings.ToLower(name)]
	if currentContents == "" {
		h[strings.ToLower(name)] = content
	} else {
		h[strings.ToLower(name)] = currentContents + ", " + content
	}

	//fmt.Println("You shouldn't be here")
	return index + 2, false, nil
}

func parseHeader(data string) (name string, value string, err error) {
	sections := strings.SplitN(data, ":", 2)
	//fmt.Println("Log: ", sections[0])
	if strings.Contains(sections[0], " ") {
		return "", "", fmt.Errorf("Malformed header: Whitespace in field-name")
	}
	name = sections[0]
	value = strings.TrimSpace(sections[1])
	return name, value, nil

}

var specialChars = []string{"!", "#", "$", "%", "&", "'", "*", "+", "-", ".", "^", "_", "`", "|", "~"}

func validateChar(c byte) (valid bool) {

	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'z') || (c >= '0' && c <= 9) {
		return true
	}

	for _, sc := range specialChars {
		if sc == string(c) {
			return true
		}
	}
	return false

}

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h Headers) Set(key, value string) {
	h[strings.ToLower(key)] = value
}
