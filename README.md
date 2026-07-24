# Tcpserver

Simple TCP echo server in Go, written three ways to compare approaches to graceful shutdown.

Each variant listens on `:8080`, accepts connections concurrently, echoes back whatever you send, and shuts down cleanly on `Ctrl+C` / `SIGTERM` — waiting for in-flight connections to finish before exiting.

### The three variants

| File | Approach | Takeaway |
|------|----------|----------|
| `server1.go` | Plain function; `main` owns the listener and `WaitGroup` | Must `ln.Close()` *before* `wg.Wait()`, or the accept loop keeps adding work and you block forever |
| `server2.go` | `Server` struct with `ListenAndServe()` / `Shutdown()`, mirroring `net/http` | `ListenAndServe` blocks, so bind errors have to come back over a channel while `main` watches the signal context |
| `server3.go` | `Server` struct with a non-blocking `Start()` / `Shutdown()` | `Start` returns bind errors synchronously, so `main` reads top to bottom; shutdown takes a `context` with a timeout |

`utils.go` holds the shared connection handler (echo, plus `quit` to disconnect).

### Running

Pick a variant in `main.go`:

```go
func main() {
	// server1()
	// server2()
	server3()
}
```

Then:

```sh
go run .
```


## Trying it out

```sh
nc localhost 8080
```

Type anything to get it echoed back; type `quit` to disconnect.

### Seeing graceful shutdown

Open `nc localhost 8080` in two terminals and echo something in each.

1. Type `quit` in the first — it disconnects, the second still echoes fine.
2. Hit `Ctrl+C` on the server — it stops accepting new connections but the second terminal still echoes.
3. Type `quit` there (or close it) and the server finally exits.

That's the point: shutdown waits for live connections instead of killing them mid-conversation.
