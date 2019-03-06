package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
)

type TCPServer struct {
	listener net.Listener
	address  string
	quit     chan bool
	exited   chan bool
	message  string
}

func echoTCPServer(address string) *TCPServer {
	addr, err := net.ResolveTCPAddr("tcp4", address)
	// TODO: return nil, error and decide how to handle it in the calling function
	if err != nil {
		fmt.Println("Failed to resolve address", err.Error())
		os.Exit(1)
	}

	// TODO: return nil, error and decide how to handle it in the calling function
	listener, err := net.Listen("tcp", addr.String())
	if err != nil {
		fmt.Println("Failed to", err.Error())
		os.Exit(1)
	}

	// TODO: do not use this syntax, add the field names
	srv := &TCPServer{
		listener: listener,
		address:  address,
		quit:     make(chan bool),
		exited:   make(chan bool),
		message:  dumpEnv(),
	}
	go srv.serve()
	return srv
}

func (srv *TCPServer) serve() {
	log.Println("TCP listen", srv.address)

	var stopped int32
	var handlers sync.WaitGroup
	go func() {
		<-srv.quit
		atomic.StoreInt32(&stopped, 1)
		fmt.Println("Shutting down...")
		srv.listener.Close()
		fmt.Println("wait connections to stop...")
		handlers.Wait()
		close(srv.exited)
	}()
	for {
		if atomic.LoadInt32(&stopped) == 1 {
			return
		}

		//fmt.Println("Listening for clients")
		// srv.listener.SetDeadline(time.Now().Add(1e9))
		conn, err := srv.listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			if atomic.LoadInt32(&stopped) == 1 {
				return
			}
			// TODO: check is in closing state
			fmt.Println("Failed to accept connection:", err.Error())
			return
		}
		handlers.Add(1)
		go func() {
			srv.handleTCPRequest(conn, 0)
			handlers.Done()
		}()
	}
}

func (srv *TCPServer) handleTCPRequest(conn net.Conn, id int) {
	log.Println("TCP connection open.", conn.RemoteAddr())
	defer func() {
		conn.Close()
		log.Println("TCP connection closed.")
	}()

	conn.Write([]byte(srv.message))

	for {
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			return
		}
		data := buf[:size]

		// log.Println(clientId+" - Received Raw Data:", data)
		// log.Printf(clientId+" - Received Data (converted to string): %s", data)
		conn.Write(data)
	}
}

func (srv *TCPServer) Stop() {
	fmt.Println("...TCP stopping listening on", srv.address)
	// XXX: You cannot use the same channel in two directions.
	//      The order of operations on the channel is undefined.
	close(srv.quit)
	// fmt.Println("wait on exit chan")
	<-srv.exited
	fmt.Println("TCP server stopped successfully on", srv.address)
}
