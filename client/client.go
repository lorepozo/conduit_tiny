package main

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/lukedmor/conduit_tiny/common"
)

func main() {
	serverConn, err := rpc.Dial("tcp", common.ClientListenerAddr)
	if err != nil {
		log.Fatal(err)
	}

	var r common.RequestProviderReply
	err = serverConn.Call("Client.RequestProvider", common.NewNothing(), &r)
	if err != nil {
		log.Fatal(err)
	}
	addr := r.Addr

	serverConn.Close()

	provider, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	args := common.Executable{
		"python",
		[]byte("print('hello, world!')"),
	}
	var reply common.ExecutionReply
	err = provider.Call("Executor.Execute", &args, &reply)
	if err != nil {
		fmt.Printf("%s", reply.Output)
		log.Fatal(err)
	}

	fmt.Printf("%s", reply.Output)
}
