package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	addr := conn.RemoteAddr().String()

	log.Printf("[%s] connected", addr)

	// Tell the client we're alive
	conn.Write([]byte("Welcome! \n Type anything and i will echo you back \n Type quit to disocnnect. \n\n"))

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("[%s] disconnected", addr)
			return
		}

		msg := string(buf[:n])
		log.Printf("[%s] received: %q", addr, msg)

		if msg == "quit\n" {
			conn.Write([]byte("Goodbyte \n"))
			log.Printf("[%s] disconnected", addr)
			return
		}

		reply := fmt.Sprintf("Echo: %s", msg)
		conn.Write([]byte(reply))
	}
}
