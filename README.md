# multiecho

[![Docker build](https://img.shields.io/docker/cloud/build/nordicdyno/multiecho.svg)][hub]
[![Docker Pulls](https://img.shields.io/docker/pulls/nordicdyno/multiecho.svg)][hub]


Echo server, listen and responses on multiple TCP/UDP ports.

Could be useful for testing is tcp/udp routing (or proxying) works correctly,
i.e. in clouds/kubernetes/docker environments etc...

## Installation

binary:

    go install github.com/nordicdyno/multiecho

docker:

    docker pull nordicdyno/multiecho

## Usage

CLI flags:

    -bind string
        default bind address (used if not provided with port) (default "0.0.0.0")
    -env
        print all environment variables on new connection start
    -t value
        tcp bind ports/addresses in ADDR:PORT format, there ADDR is optional (could be used multiple times)
    -u value
        udp bind ports/addresses in ADDR:PORT format, there ADDR is optional (could be used multiple times)

supported environment variables:

* LISTEN_HOST overwrites -bind flag.
* TCP_LISTEN_PORTS overwrites -t flags.
* UDP_LISTEN_PORTS overwrites -u flags.

examples:

    TCP_LISTEN_PORTS=1544,127.0.0.1:1522 multiecho -u 2444 -u 127.0.0.1:2445

    TCP_LISTEN_PORTS=1544,1522 multiecho

    UDP_LISTEN_PORTS=1544,1522 multiecho -t 2344


## Kubernetes

_TODO: add manifest sample_

[hub]: https://hub.docker.com/r/nordicdyno/multiecho/
