package utils

import (
	"bytes"
	"errors"
	"strconv"
)

const (
	kMaxMsg  = 4096
	kMaxArgs = 128
)

type Conn struct {
	Rbuf     []byte
	RbufSize int
	Wbuf     []byte
	WbufSize int
	State    int
}

const (
	STATE_REQ = iota
	STATE_RES
	STATE_END
)

func ParseReq(data []byte, length uint32) ([]string, error) {
	if length < 4 {
		return nil, errors.New("invalid length")
	}

	if data[0] != '*' {
		return nil, errors.New("invalid protocol")
	}

	// Read number of arguments
	parts := bytes.Split(data[1:length], []byte("\r\n"))
	if len(parts) < 3 {
		return nil, errors.New("invalid message format")
	}

	numArgs, err := strconv.Atoi(string(parts[0]))
	if err != nil || numArgs < 1 || numArgs > kMaxArgs {
		return nil, errors.New("invalid number of arguments")
	}

	var out []string
	for i := 1; i < len(parts)-1; i += 2 {
		out = append(out, string(parts[i+1]))
	}

	return out, nil
}
