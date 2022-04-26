package main

import (
	"io"
	"log"
	"net"
)

type Proxy struct {
	ListenAddr string
	TargetAddr string
	exit       chan struct{}
}

func (p *Proxy) log(args ...interface{}) {
	if ShouldLog {
		log.Println("[DEBUG] ", "["+p.ListenAddr+"] ", args)
	}
}
func (p *Proxy) Start() {
	l, err := net.Listen("tcp", p.ListenAddr)
	if err != nil {
		panic(err)
	}
	go func() {
		p.log("Waiting For Connection")
		for {
			clientConn, err := l.Accept()
			if err != nil {
				p.log(err.Error())
				continue
			}
			p.log("Connection From", clientConn.RemoteAddr().String())
			go p.handleConnection(clientConn)
		}
	}()
	<-p.exit
	l.Close()
}
func (p *Proxy) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()
	serverConn, err := net.Dial("tcp", p.TargetAddr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer serverConn.Close()
	go io.Copy(serverConn, clientConn)
	io.Copy(clientConn, serverConn)
	p.log("Connection Closed", clientConn.RemoteAddr().String())
}
func (p *Proxy) Stop() {
	p.exit <- struct{}{}
}
func NewProxy(listenAddr, targetAddr string) *Proxy {
	return &Proxy{
		ListenAddr: listenAddr,
		TargetAddr: targetAddr,
		exit:       make(chan struct{}),
	}
}
