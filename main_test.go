/**
 * Created by zhangruizhi on 2022/12/14
 */

package main

import (
	"context"
	"fmt"
	"github.com/pion/mdns"
	"golang.org/x/net/ipv4"
	"net"
	"testing"
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
	fmt.Println(answer)
	fmt.Println(src)
	fmt.Println(err)
}
