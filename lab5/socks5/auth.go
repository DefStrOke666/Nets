package socks5

import (
	"fmt"
	"io"
)

const (
	noAuthMethod        = uint8(0)
	noAcceptableMethods = uint8(255)
)

func (s *Server) authenticate(conn io.Writer, bufConn io.Reader) error {
	methods, err := readMethods(bufConn)
	if err != nil {
		return fmt.Errorf("failed to get auth methods: %v", err)
	}

	for _, method := range methods {
		if method == noAuthMethod {
			_, err := conn.Write([]byte{socks5Ver, noAuthMethod})
			return err
		}
	}

	_, err = conn.Write([]byte{socks5Ver, noAcceptableMethods})
	if err != nil {
		return err
	}

	return fmt.Errorf("no supported authentication mechanism")
}

func readMethods(r io.Reader) ([]byte, error) {
	header := []byte{0}
	if _, err := r.Read(header); err != nil {
		return nil, err
	}

	numMethods := int(header[0])
	methods := make([]byte, numMethods)
	_, err := io.ReadAtLeast(r, methods, numMethods)
	return methods, err
}
