package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
)

type UDPServer struct {
	conn    *net.UDPConn
	address string
	quit    chan bool
	exited  chan bool
	message string
}

func echoUDPServer(address string) *UDPServer {
	addr, err := net.ResolveUDPAddr("udp", address)
	// TODO: return nil, error and decide how to handle it in the calling function
	if err != nil {
		fmt.Println("Failed to resolve address", err.Error())
		os.Exit(1)
	}

	// TODO: return nil, error and decide how to handle it in the calling function
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Failed to", err.Error())
		os.Exit(1)
	}

	// TODO: do not use this syntax, add the field names
	srv := &UDPServer{
		conn:    conn,
		address: address,
		quit:    make(chan bool),
		exited:  make(chan bool),
		message: dumpEnv(),
	}
	// TODO: no need to export Serve as it is only called internally
	go srv.Serve()
	return srv
}

func (srv *UDPServer) Serve() {
	log.Println("UDP listen", srv.address)
	var handlers sync.WaitGroup
	var stopped int32
	go func() {
		<-srv.quit
		atomic.StoreInt32(&stopped, 1)
		srv.conn.Close()
		// srv.listener.Close()
		// fmt.Println("wait connections to stop...")
		handlers.Wait()
		close(srv.exited)
	}()
	// TODO: add timeout (just for fun)
	var seenMu sync.Mutex
	seen := map[string]bool{}
	for {
		if atomic.LoadInt32(&stopped) == 1 {
			return
		}

		//fmt.Println("Listening for clients")
		// srv.listener.SetDeadline(time.Now().Add(1e9))
		buf := make([]byte, 2048)
		n, addr, err := srv.conn.ReadFromUDP(buf)
		// addr.String
		if err != nil {
			if atomic.LoadInt32(&stopped) == 1 {
				return
			}
			fmt.Printf("Failed to accept connection on %v: %v\n", srv.address, err.Error())
		}
		handlers.Add(1)
		go func(buf []byte, addr *net.UDPAddr) {
			seenMu.Lock()
			showmsg := seen[addr.String()]
			seenMu.Unlock()
			if !showmsg {
				srv.write(addr, []byte(srv.message))
			}
			seenMu.Lock()
			seen[addr.String()] = true
			seenMu.Unlock()

			srv.write(addr, buf)

			// srv.handleTCPRequest(conn, 0)
			handlers.Done()
		}(buf[:n], addr)
	}
}

func (srv *UDPServer) write(addr *net.UDPAddr, payload []byte) {
	var buf []byte
	tail := payload
	for len(tail) > 0 {
		if len(tail) < 512 {
			buf = tail[:len(tail)]
			tail = []byte{}
		} else {
			buf = tail[:512]
			tail = tail[512:]
		}

		fmt.Println("write to", addr.String())
		fmt.Println(string(buf))
		_, err := srv.conn.WriteTo(buf, addr)
		if err != nil {
			fmt.Printf("net.WriteTo() to %v failed: %s\n", addr.String(), err)
		}
	}
}

func (srv *UDPServer) Stop() {
	fmt.Println("Stop listening on", srv.address)
	// XXX: You cannot use the same channel in two directions.
	//      The order of operations on the channel is undefined.
	close(srv.quit)
	fmt.Println("wait on exit chan")
	<-srv.exited
	fmt.Println("Stopped successfully")
}
