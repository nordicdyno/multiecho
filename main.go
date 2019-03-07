package main

// TODO:
// * tests and refactoring
// * better logging
// * check is there any connections leaks
//
// MAYBE:
// timeouts

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
)

func main() {
	cfg := configure()
	// fmt.Printf("%#v\n", cfg)
	// return
	if !cfg.showEnv {
		dumpEnv = func() string { return "" }
	}

	var servers []Server
	for _, addr := range cfg.tcpAddrs {
		srv := echoTCPServer(addr)
		servers = append(servers, srv)
	}

	for _, addr := range cfg.udpAddrs {
		srv := echoUDPServer(addr)
		servers = append(servers, srv)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigc
	for _, srv := range servers {
		srv.Stop()
	}
}

type Server interface {
	Stop()
}

func envname(s string) string {
	return strings.SplitN(s, "=", 1)[0]
}

var dumpEnv = func() string {
	envs := os.Environ()
	sort.SliceStable(envs, func(i, j int) bool {
		return envname(envs[i]) < envname(envs[j])
	})
	return strings.Join(
		[]string{
			strings.Repeat("-", 45),
			strings.Join(envs, "\n"),
			strings.Repeat("-", 45),
			"\n",
		},
		"\n",
	)
}

func fatalerror(format string, a ...interface{}) {
	errorf(format, a...)
	os.Exit(1)
}

func errorf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr, "")
}
