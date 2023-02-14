package rpc

import (
	"net"
	"reflect"
)

type RPCData struct {
	Name string
	Args []interface{}
}

type Server struct {
	addr     string
	funcMaps map[string]reflect.Value
}

func NewServer(addr string) *Server {
	return &Server{
		addr:     addr,
		funcMaps: map[string]reflect.Value{},
	}
}

func (s *Server) Register(name string, f any) {
	s.funcMaps[name] = reflect.ValueOf(f)
}

func (s *Server) Run() {
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			panic(err)
		}
		//create a session
		session := NewSession(conn, WithLen(4))
		data, err := session.Read()
		if err != nil {
			panic(err)
		}
		v, err := decode(data)
		if err != nil {
			panic(err)
		}
		f := s.funcMaps[v.Name]
		args := make([]reflect.Value, len(v.Args))
		for _, v := range v.Args {
			args = append(args, reflect.ValueOf(v))
		}
		rsp := f.Call(args)
		rspArgs := make([]interface{}, len(rsp))
		for _, v := range rsp {
			rspArgs = append(rspArgs, v.Interface())
		}
		rspData := RPCData{
			Name: v.Name,
			Args: rspArgs,
		}
		bs, err := encode(rspData)
		if err != nil {
			panic(err)
		}
		err = session.Write(bs)
		if err != nil {
			panic(err)
		}
	}
}
