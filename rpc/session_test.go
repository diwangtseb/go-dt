package rpc

import (
	"fmt"
	"net"
	"testing"
	"time"
)

const P_LEN = 4

func Test_Session(t *testing.T) {
	ch := make(chan struct{})
	go func() {
		fmt.Println("start")
		lis, err := net.Listen("tcp", ":1234")
		if err != nil {
			panic(err)
		}

		ch <- struct{}{}
		conn, _ := lis.Accept()
		session := NewSession(conn, WithLen(P_LEN))
		fmt.Println(session.sessionOption.len)
		err = session.Write([]byte("hello"))
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		<-ch
		fmt.Println("start 2")
		conn, err := net.Dial("tcp", ":1234")
		if err != nil {
			panic(err)
		}
		session := NewSession(conn, WithLen(P_LEN))
		data, err := session.Read()
		if err != nil {
			panic(err)
		}
		fmt.Printf("read%v", string(data))
	}()
	time.Sleep(time.Second * 1)
}
