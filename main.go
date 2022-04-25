package main

import (
	"io"
	"log"
	"net"
)

func main() {
	l , err := net.Listen("tcp","0.0.0.0:3050")
	if err != nil {
		panic(err)
	}
	log.Println("Waiting For Connection")
	for {
		clientConn, err := l.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Println("Connection From" , clientConn.RemoteAddr().String())
		go func() {
			defer clientConn.Close()
			serverConn , err := net.Dial("tcp","172.16.0.10:80")
			if err != nil {
				log.Println(err.Error())
				return
			}
			defer serverConn.Close()
			go io.Copy(serverConn,clientConn)
			io.Copy(clientConn,serverConn)
			log.Println("Proxy Finished",clientConn.RemoteAddr().String())
		}()
	}
}