package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// main loop running the listener owns everything, including the waitgroup. The accept loop is a goroutine that adds to the waitgroup, and each connection handler is a goroutine that also adds to the waitgroup. The main loop waits for the signal context to be done, then closes the listener and waits for all goroutines to finish.
func server1() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server1 started on %s", ln.Addr().String())

	var wg sync.WaitGroup

	startAcceptLoop(ln, &wg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	// correct order -> close first and then wait
	// Never wg.Wait() before ln.Close(), the program blocks forever, because the accept  loop is still alive and could still be adding new connections to the waitgroup
	ln.Close()
	wg.Wait()

	log.Println("stopped")

}

func startAcceptLoop(ln net.Listener, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			wg.Add(1)
			go handleConnection(conn, wg)
		}
	}()
}
