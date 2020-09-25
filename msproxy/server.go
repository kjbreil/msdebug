package main

import (
	"log"
	"net"

	"github.com/kjbreil/msdebug"
)

func server() error {
	// Start the TCP server
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr())
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	// start the receiving mailbox callback
	s := msdebug.NewClient(*callbackMailslot)
	err = s.Start()
	if err != nil {
		log.Panicln(err)
	}

	// nested for loop to handle dropped connection
	for {
		// listen for an incoming connection
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		for {
			// receive message
			buf := make([]byte, 512)
			n, err := conn.Read(buf[0:])
			if err != nil {
				// error reading indicates client disconnected, close the connection and start loop over again
				conn.Close()
				log.Println("client connection closed")
				break
			}

			m := msdebug.Message{
				Data:     string(buf[:n]),
				Callback: *callbackMailslot,
			}
			m.Send(*serverMailslot)
			resp := <-s.Send
			log.Println("response", resp)
		}
	}
}
