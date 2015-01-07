package main

import (
	"fmt"
	"sync"
	"testing"
)

var servers = map[string]string{
	"desteria": "Play.NirvanaMC.com:25565",
	"gotpvp":   "play.gotpvp.com:25565",
}

func TestPing(t *testing.T) {
	wg := &sync.WaitGroup{}

	for name, host := range servers {
		wg.Add(1)

		go func(name string, host string) {
			fmt.Println("started")

			svr := NewServer(name, host)
			if err := svr.connect(); err != nil {
				panic(err)
			}

			fmt.Println("connected: " + svr.name)

			pong, err := svr.Ping()
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
