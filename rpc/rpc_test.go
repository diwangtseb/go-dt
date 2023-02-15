package rpc

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestRpc(t *testing.T) {
	s := NewServer(":1234")
	s.Register("getUid", getUid)
	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
		s.Run()
	}()
	go func() {
		<-ch
		time.Sleep(time.Second * 1)
		conn, err := net.Dial("tcp", ":1234")
		if err != nil {
			fmt.Println(err)
		}
		client := NewClient(conn)
		var gu func() (string, error)
		client.callRpc("getUid", &gu)
		r, err := gu()
		if err != nil {
			panic(err)
		}
		fmt.Println(r)
	}()
	time.Sleep(time.Second * 5)
}

func getUid() (string, error) {
	return "1", nil
}
