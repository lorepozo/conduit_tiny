package common

const (
	ProviderListenerAddr = ":8000"
	ClientListenerAddr   = ":8001"
)

type Nothing struct{}

func NewNothing() *Nothing {
	return new(Nothing)
}

type RequestProviderReply struct {
	Addr string
}

type RequestProviderError struct {
	Msg string
}

func (e *RequestProviderError) Error() string {
	return e.Msg
}

type ProviderJoinLeaveArgs struct {
	Addr string
}

type Executable struct {
	Interpreter string
	Content     []byte
}

type ExecutionReply struct {
	Output []byte
}
