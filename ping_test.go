package minepong

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

var testServers = map[string]string{
	"desteria": "Play.NirvanaMC.com:25565",
	"gotpvp":   "play.gotpvp.com:25565",
}

type server struct {
	name string
	host string

	conn net.Conn
}

func newServer(name string, host string) *server {
	return &server{
		name: name,
		host: host,
	}
}

func (s *server) connect() error {
	var err error
	s.conn, err = net.Dial("tcp", s.host)
	if err != nil {
		return err
	}

	return nil
}

func (s *server) disconnect() error {
	return s.conn.Close()
}

func TestPing(t *testing.T) {
	wg := &sync.WaitGroup{}

	for name, host := range testServers {
		wg.Add(1)

		go func(name string, host string) {
			fmt.Println("started")

			svr := newServer(name, host)
			if err := svr.connect(); err != nil {
				panic(err)
			}

			fmt.Println("connected: " + svr.name)

			pong, err := Ping(svr.conn, svr.host)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(svr.name, pong.Players.Online, pong.Players.Max)

			wg.Done()
		}(name, host)
	}

	wg.Wait()
}
