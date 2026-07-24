package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server3 struct {
	ln net.Listener
	wg sync.WaitGroup
}

func (s *Server3) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("Listen on %s: %w", addr, err)
	}
	s.ln = ln

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			s.wg.Add(1)
			go handleConnection(conn, &s.wg)
		}
	}()
	return nil
}

func (s *Server3) Shutdown(ctx context.Context) error {
	closeErr := s.ln.Close()
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return closeErr
	case <-ctx.Done():
		return ctx.Err()
	}
}

func server3() {
	srv := &Server3{}

	// Reads a linear scirpt, bind failure returns synchronously before anything else happens

	if err := srv.Start(":8080"); err != nil {
		log.Fatal(err)
	}
	log.Printf("Server3 started on %s", srv.ln.Addr().String())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error occurred while shutting down Server3: %v", err)
	}
	log.Println("Server3 stopped")
}
