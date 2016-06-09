package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"path"

	"github.com/lukedmor/conduit_tiny/common"
)

type Executor struct{}

func (xr *Executor) Execute(x *common.Executable, r *common.ExecutionReply) error {
	dir := os.TempDir()
	fileName := path.Join(dir, "content")
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Errorf("couldn't create file in temp dir: %s", err)
		return err
	}

	_, err = file.Write(x.Content)
	if err != nil {
		fmt.Errorf("couldn't write content to temp file: %s", err)
		return err
	}
	file.Close()

	// should validate x.Interpreter first, this is hacky
	// and must fit paradigm like `python file.py`
	cmd := exec.Command(x.Interpreter, fileName)
	r.Output, err = cmd.CombinedOutput()
	return err
}

type provider struct {
	server *rpc.Client
	addr   string
}

func newProvider() *provider {
	p := new(provider)
	serverConn, err := rpc.Dial("tcp", common.ProviderListenerAddr)
	if err != nil {
		log.Fatalf("couldn't connect to conduit: %s", err)
	}
	p.server = serverConn
	p.addr = ":8002"
	return p
}

func (p *provider) terminate() {
	args := common.ProviderJoinLeaveArgs{p.addr}
	err := p.server.Call("Provider.Leave", &args, nil)
	if err != nil {
		fmt.Errorf("error leaving conduit: %s", err)
	}
}

func (p *provider) join() {
	args := common.ProviderJoinLeaveArgs{p.addr}
	err := p.server.Call("Provider.Join", &args, nil)
	if err != nil {
		log.Fatalf("couldn't join conduit: %s", err)
	}
}

func (p *provider) listen() {
	rpc.Register(new(Executor))
	l, err := net.Listen("tcp", p.addr)
	if err != nil {
		log.Fatal(err)
	}
	go p.join()
	rpc.Accept(l)
}

func main() {
	p := newProvider()
	defer p.terminate()
	p.listen()
}
