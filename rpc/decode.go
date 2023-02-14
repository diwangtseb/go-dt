package rpc

import (
	"bytes"
	"encoding/gob"
)

func decode(data []byte) (RPCData, error) {
	buffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffer)
	var rsp RPCData
	err := dec.Decode(&rsp)
	if err != nil {
		panic(err)
	}
	return rsp, nil
}

func encode(data RPCData) ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(data)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes(), nil
}
