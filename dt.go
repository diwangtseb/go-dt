package dt

import (
	"context"

	"github.com/dtm-labs/client/dtmgrpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type MethodPair struct {
	Action     string
	Compensate string
	ProtoMsg   protoreflect.ProtoMessage
}

type TransactionActor interface {
	ExecuteSaga(ctx context.Context, methodPairs ...MethodPair) error
	ExecuteMsg(ctx context.Context, methodPair ...MethodPair) error
}

type transactionActor struct {
	dtmServerAddr  string
	grpcServerAddr string
	gid            string
}

// MsgExecute implements TransactionActor
func (ta *transactionActor) ExecuteMsg(ctx context.Context, methodPair ...MethodPair) error {
	msg := dtmgrpc.NewMsgGrpc(ta.dtmServerAddr, ta.gid)
	for _, v := range methodPair {
		msg = msg.Add(ta.withServerAction(v.Action), v.ProtoMsg)
	}
	msg.WaitResult = true
	if err := msg.Submit(); err != nil {
		return err
	}
	return nil
}

func NewTransactionActor(dtmAddr, grpcAddr string) TransactionActor {
	return &transactionActor{
		dtmServerAddr:  dtmAddr,
		grpcServerAddr: grpcAddr,
		gid:            dtmgrpc.MustGenGid(dtmAddr),
	}
}

func (ta *transactionActor) withServerAction(action string) string {
	return ta.grpcServerAddr + action
}
func (ta *transactionActor) withServerCompensate(compensate string) string {
	return ta.grpcServerAddr + compensate
}

// Regist implements TransactionRegister
func (ta *transactionActor) ExecuteSaga(ctx context.Context, methodPairs ...MethodPair) error {
	saga := dtmgrpc.NewSagaGrpc(ta.dtmServerAddr, ta.gid)

	for _, methodPair := range methodPairs {
		saga = saga.Add(ta.withServerAction(methodPair.Action), ta.withServerCompensate(methodPair.Compensate), methodPair.ProtoMsg)
	}
	saga.WaitResult = true
	if err := saga.Submit(); err != nil {
		return err
	}
	return nil
}

var _ TransactionActor = (*transactionActor)(nil)
