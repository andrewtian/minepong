package minepong

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
	"strconv"
)

const (
	protocolVersion = 0x47
)

type Pong struct {
	Version struct {
		Name     string
		Protocol int
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []map[string]string
	} `json:"players"`
	Description interface{} `json:"description"`
	FavIcon     string      `json:"favicon"`
}

func Ping(conn net.Conn, host string) (*Pong, error) {
	if err := SendHandshake(conn, host); err != nil {
		return nil, err
	}

	if err := SendStatusRequest(conn); err != nil {
		return nil, err
	}

	pong, err := ReadPong(conn)
	if err != nil {
		return nil, err
	}

	return pong, nil
}

func makePacket(pl *bytes.Buffer) *bytes.Buffer {
	var buf bytes.Buffer
	// get payload length
	buf.Write(encodeVarint(uint64(len(pl.Bytes()))))

	// write payload
	buf.Write(pl.Bytes())

	return &buf
}

func SendHandshake(conn net.Conn, host string) error {
	pl := &bytes.Buffer{}

	// packet id
	pl.WriteByte(0x00)

	// protocol version
	pl.WriteByte(protocolVersion)

	// server address
	host, port, err := net.SplitHostPort(host)
	if err != nil {
		panic(err)
	}

	pl.Write(encodeVarint(uint64(len(host))))
	pl.WriteString(host)

	// server port
	iPort, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	binary.Write(pl, binary.BigEndian, int16(iPort))

	// next state (status)
	pl.WriteByte(0x01)

	if _, err := makePacket(pl).WriteTo(conn); err != nil {
		return errors.New("cannot write handshake")
	}

	return nil
}

func SendStatusRequest(conn net.Conn) error {
	pl := &bytes.Buffer{}

	// send request zero
	pl.WriteByte(0x00)

	if _, err := makePacket(pl).WriteTo(conn); err != nil {
		return errors.New("cannot write send status request")
	}

	return nil
}

// https://code.google.com/p/goprotobuf/source/browse/proto/encode.go#83
func encodeVarint(x uint64) []byte {
	var buf [10]byte
	var n int
	for n = 0; x > 127; n++ {
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	buf[n] = uint8(x)
	n++
	return buf[0:n]
}

func ReadPong(rd io.Reader) (*Pong, error) {
	r := bufio.NewReader(rd)
	nl, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, errors.New("could not read length")
	}

	pl := make([]byte, nl)
	_, err = io.ReadFull(r, pl)
	if err != nil {
		return nil, errors.New("could not read length given by length header")
	}

	// packet id
	_, n := binary.Uvarint(pl)
	if n <= 0 {
		return nil, errors.New("could not read packet id")
	}

	// string varint
	_, n2 := binary.Uvarint(pl[n:])
	if n2 <= 0 {
		return nil, errors.New("could not read string varint")
	}

	var pong Pong
	if err := json.Unmarshal(pl[n+n2:], &pong); err != nil {
		return nil, errors.New("could not read pong json")
	}

	return &pong, nil
}
