package rpc

import "testing"

func TestRpc(t *testing.T) {
	s := NewServer(":1234")
	s.Register("getUid", getUid)
	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
	}()
}

func getUid() (string, error) {
	return "1", nil
}
