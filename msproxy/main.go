package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	bindAdd          = flag.String("ip", "0.0.0.0", "ip to connect to as client or ip to lisnen to as server")
	bindPort         = flag.String("port", "6666", "port to connect to as client or port to lisnen to as server")
	isServer         = flag.Bool("server", false, "run as a server")
	serverMailslot   = flag.String("ms", `\\.\mailslot\win900`, "server mailslot location")
	callbackMailslot = flag.String("cb", `\\.\mailslot\win900cb`, "server callback mailslot location")
	isClient         = flag.Bool("client", false, "run as a client")
	clientMailslot   = flag.String("cms", `\\.\mailslot\win900pr`, "the proxy mailslot")
)

func main() {
	// parse the command line flags
	flag.Parse()

	switch {
	case *isServer:
		err := server()
		if err != nil {
			log.Panicln(err)
		}
	case *isClient:
		err := client()
		if err != nil {
			log.Panicln(err)
		}
	default:
		log.Fatalf("neither server or client command line passed")
	}

}

func serverAddr() string {
	return fmt.Sprintf("%s:%s", *bindAdd, *bindPort)
}
