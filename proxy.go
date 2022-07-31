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
	Stats      Stats
}

func (p *Proxy) log(args ...interface{}) {
	Log(LOG_DEBUG, "[PROXY]", p.ListenAddr, args)
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
	if addr, ok := clientConn.RemoteAddr().(*net.TCPAddr); ok {
		ip := addr.IP.String()
		_, exists := userDB.Get(ip)
		p.Stats.TotalConnections++
		switch conf.Users.Mode {
		case "whitelist":
			if !exists {
				clientConn.Close()
				p.Stats.RejectedConnections++
				p.log("[WHITELIST]", ip, "doesn't exists in user db")
				return
			}
			p.Stats.AcceptedConnections++
		case "blacklist":
			if exists {
				clientConn.Close()
				p.Stats.RejectedConnections++
				p.log("[BLACKLIST]", ip, "exists in user db")
				return
			}
			p.Stats.AcceptedConnections++
		}
	}
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
		Stats: Stats{
			TotalConnections:    0,
			AcceptedConnections: 0,
			RejectedConnections: 0,
		},
	}
}
