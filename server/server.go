package main

import (
	"log"
	"net"
	"net/rpc"
	"sync"

	"github.com/lukedmor/conduit_tiny/common"
)

type Provider struct {
	cs *conduitServer
}

func (p *Provider) Join(a *common.ProviderJoinLeaveArgs, reply *common.Nothing) error {
	p.cs.providers[a.Addr] = true
	return nil
}

func (p *Provider) Leave(a *common.ProviderJoinLeaveArgs, reply *common.Nothing) error {
	delete(p.cs.providers, a.Addr)
	return nil
}

type Client struct {
	cs *conduitServer
}

func (c *Client) RequestProvider(args *common.Nothing, r *common.RequestProviderReply) error {
	for a := range c.cs.providers { // random iteration
		r.Addr = a
		return nil
	}
	return &common.RequestProviderError{"no providers available"}
}

// internals

type conduitServer struct {
	sync.Mutex

	providers map[string]bool // key is host addr

	fail chan error
}

func newConduitServer() *conduitServer {
	cs := new(conduitServer)
	cs.providers = make(map[string]bool)
	cs.fail = make(chan error)
	return cs
}

func (cs *conduitServer) run() {
	go cs.listenForProviders()
	go cs.listenForClients()
	log.Fatal(<-cs.fail)
}

func (cs *conduitServer) listenForProviders() {
	s := rpc.NewServer()
	s.Register(&Provider{cs})
	l, err := net.Listen("tcp", common.ProviderListenerAddr)
	if err != nil {
		cs.fail <- err
	}
	s.Accept(l)
}

func (cs *conduitServer) listenForClients() {
	s := rpc.NewServer()
	s.Register(&Client{cs})
	l, err := net.Listen("tcp", common.ClientListenerAddr)
	if err != nil {
		cs.fail <- err
	}
	s.Accept(l)
}

func main() {
	newConduitServer().run()
}
