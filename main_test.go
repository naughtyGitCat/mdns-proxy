/**
 * Created by zhangruizhi on 2022/12/14
 */

package main

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/pion/mdns"
	"golang.org/x/net/ipv4"
)

func TestMDns(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", mdns.DefaultAddress)
	if err != nil {
		panic(err)
	}

	l, err := net.ListenUDP("udp4", addr)
	if err != nil {
		panic(err)
	}

	server, err := mdns.Server(ipv4.NewPacketConn(l), &mdns.Config{})
	if err != nil {
		panic(err)
	}
	answer, src, err := server.Query(context.TODO(), "pve.local")
	fmt.Printf("answer: %s\n", answer.GoString())
	fmt.Printf("src: %s\n", src.(*net.UDPAddr).IP.String())
	fmt.Printf("err: %s\n", err)
}
