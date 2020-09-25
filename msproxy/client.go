package main

import "net"

func client() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr())
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}

	// send message
	_, err = conn.Write([]byte("Installer calling server"))
	if err != nil {
		panic(err)
	}

}
