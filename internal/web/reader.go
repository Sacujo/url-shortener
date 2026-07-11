package web

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func ReadRequest(reader *bufio.Reader) (Request, error) {
	req, err := readRequestLine(reader)
	if err != nil {
		return Request{}, err
	}

	headers, err := readHeaders(reader)
	if err != nil {
		return Request{}, err
	}

	req.Headers = headers

	body, err := readBody(reader, headers)
	if err != nil {
		return Request{}, err
	}

	req.Body = body

	return req, nil
}

func readRequestLine(reader *bufio.Reader) (Request, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return Request{}, err
	}
	line = strings.TrimSpace(line)
	return parseRequestLine(line)
}

func readHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return nil, fmt.Errorf("invalid header: %q", line)
		}

		headers[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return headers, nil
}

func readBody(reader *bufio.Reader, headers map[string]string) ([]byte, error) {
	if contentLengthStr, ok := headers["Content-Length"]; ok {
		var contentLength int
		_, err := fmt.Sscanf(contentLengthStr, "%d", &contentLength)
		if err != nil {
			return nil, fmt.Errorf("invalid Content-Length: %q", contentLengthStr)
		}

		body := make([]byte, contentLength)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}
	return nil, nil
}

func parseRequestLine(line string) (Request, error) {
	var request Request
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return request, fmt.Errorf("invalid request line: %s", line)
	}

	request.Method = parts[0]
	request.Path = parts[1]
	request.Version = parts[2]

	return request, nil
}
