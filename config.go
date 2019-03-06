package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	defaultHost = "0.0.0.0"
	// defaultTCPPort = "3333"
	// defaultUDPPort = "3334"

	defaultHostEnvName = "LISTEN_HOST"
	tcpPortsEnvName    = "TCP_LISTEN_PORTS"
	udpPortsEnvName    = "UDP_LISTEN_PORTS"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type config struct {
	tcpAddrs []string
	udpAddrs []string
	showEnv  bool
	verbose  bool
}

// var Usage = func() {
// 	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
//
// 	flag.PrintDefaults()
// }

func configure() config {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("  supported environment variables:\n")
		fmt.Printf("    %v overwrites -bind flag.\n", defaultHostEnvName)
		fmt.Printf("    %v overwrites -t flags.\n", tcpPortsEnvName)
		fmt.Printf("    %v overwrites -u flags.\n", udpPortsEnvName)
		fmt.Println()
		fmt.Println("  examples:")
		fmt.Println("")
		fmt.Println("    TCP_LISTEN_PORTS=1544,127.0.0.1:1522 multiecho -u 2444 -u 127.0.0.1:2445")
		fmt.Println()
		fmt.Println("    TCP_LISTEN_PORTS=1544,1522 multiecho")
		fmt.Println()
	}

	var (
		tcpAddr arrayFlags
		udpAddr arrayFlags
	)
	flagHost := flag.String("bind", "0.0.0.0",
		"default bind address (used if not provided with port)")
	flag.Var(&tcpAddr, "t",
		"tcp bind ports/addresses in ADDR:PORT format, there ADDR is optional (could be used multiple times)")
	flag.Var(&udpAddr, "u",
		"udp bind ports/addresses in ADDR:PORT format, there ADDR is optional (could be used multiple times)")

	// verbose := flag.Bool("verbose", false, "verbose logs")
	showEnv := flag.Bool("env", false, "print all environment variables to all new connection")
	flag.Parse()

	host := os.Getenv(defaultHostEnvName)
	if len(host) == 0 {
		host = *flagHost
	}

	cfg := config{
		tcpAddrs: addrs(tcpPortsEnvName, host, tcpAddr),
		udpAddrs: addrs(udpPortsEnvName, host, udpAddr),
		showEnv:  *showEnv,
	}
	if len(cfg.tcpAddrs) == 0 && len(cfg.udpAddrs) == 0 {
		fmt.Fprintf(os.Stderr, "Error: at least one port should be used\n\n")
		flag.Usage()
		os.Exit(1)
	}
	return cfg
}

func addrs(envname, defaulthost string, defaultaddrs []string) []string {
	var out []string

	// apply env vars
	outList := os.Getenv(envname)
	if len(outList) != 0 {
		for _, port := range strings.Split(outList, ",") {
			out = append(out, port)
		}
	} else {
		for _, addr := range defaultaddrs {
			out = append(out, addr)
		}
	}

	// make ip part explicit
	for i, addr := range out {
		parts := strings.SplitN(addr, ":", 2)
		if len(parts) < 2 {
			out[i] = defaulthost + ":" + parts[0]
		} else if len(parts[0]) == 0 {
			out[i] = defaulthost + ":" + parts[1]
		}
	}

	return out
}
