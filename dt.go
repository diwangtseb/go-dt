package dt

import (
	"context"
	"log"

	"github.com/dtm-labs/client/dtmgrpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type MethodPair struct {
	Action     string
	Compensate string
	ProtoMsg   protoreflect.ProtoMessage
}

type TransactionActor interface {
	ExecuteSaga(ctx context.Context, methodPairs ...MethodPair)
	ExecuteMsg(ctx context.Context, methodPair ...MethodPair)
}

type transactionActor struct {
	dtmServerAddr  string
	grpcServerAddr string
	gid            string
	logger         log.Logger
}

// MsgExecute implements TransactionActor
func (ta *transactionActor) ExecuteMsg(ctx context.Context, methodPair ...MethodPair) {
	msg := dtmgrpc.NewMsgGrpc(ta.dtmServerAddr, ta.gid)
	for _, v := range methodPair {
		msg = msg.Add(ta.withServerAction(v.Action), v.ProtoMsg)
	}
	msg.WaitResult = true
	if err := msg.Submit(); err != nil {
		ta.logger.Fatal(err)
	}
}

func NewTransactionActor(dtmAddr, grpcAddr string) TransactionActor {
	return &transactionActor{
		dtmServerAddr:  dtmAddr,
		grpcServerAddr: grpcAddr,
		gid:            dtmgrpc.MustGenGid(dtmAddr),
		logger:         log.Logger{},
	}
}

func (ta *transactionActor) withServerAction(action string) string {
	return ta.grpcServerAddr + action
}
func (ta *transactionActor) withServerCompensate(compensate string) string {
	return ta.grpcServerAddr + compensate
}

// Regist implements TransactionRegister
func (ta *transactionActor) ExecuteSaga(ctx context.Context, methodPairs ...MethodPair) {
	saga := dtmgrpc.NewSagaGrpc(ta.dtmServerAddr, ta.gid)

	for _, methodPair := range methodPairs {
		saga = saga.Add(ta.withServerAction(methodPair.Action), ta.withServerCompensate(methodPair.Compensate), methodPair.ProtoMsg)
	}
	saga.WaitResult = true
	if err := saga.Submit(); err != nil {
		ta.logger.Fatal(err)
	}
}

var _ TransactionActor = (*transactionActor)(nil)
