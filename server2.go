package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// clone of http.ListenAndServer in standard library

type Server2 struct {
	ln net.Listener
	wg sync.WaitGroup
}

func (s *Server2) ListenAndServe() error {
	ln, err := net.Listen("tcp", ":8080")
	log.Printf("Server2 started on %s", ln.Addr().String())

	if err != nil {
		return err
	}
	s.ln = ln

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		s.wg.Add(1)
		go handleConnection(conn, &s.wg)
	}
}

func (s *Server2) Shutdown(ctx context.Context) error {
	s.ln.Close()
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func server2() {
	srv := &Server2{}

	// main does't call listenServe directly because it blocks and main also needs to watch the signal context at same time
	// So bind error has to travel back over a channel instead of being a normal return value

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		// bind failure (eg port already in use) will be reported here
		if err != nil && !errors.Is(err, net.ErrClosed) {
			log.Fatal(err)
		}
	case <-ctx.Done():
		// signal received, shutdown gracefully
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	}
}
