package main

import (
	"log"
	"net"

	"github.com/kjbreil/msdebug"
)

func client() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr())
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	l := msdebug.NewClient(*clientMailslot)
	err = l.Start()
	if err != nil {
		log.Panicln(err)
	}

	for {
		m := <-l.Send

		log.Println("received message", m)

		_, err = conn.Write([]byte(m.Data))
		if err != nil {
			return err
		}

		// receive message
		buf := make([]byte, 512)
		n, err := conn.Read(buf[0:])
		if err != nil {
			// handle error
		}

		m.Data = string(buf[:n])

		l.Receive <- m
	}
}
