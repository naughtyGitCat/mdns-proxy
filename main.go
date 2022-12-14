/**
 * Created by psyduck on 2022/12/14
 */

package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/akamensky/argparse"
	"github.com/pkg/errors"

	"github.com/miekg/dns"
	"github.com/pion/mdns"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/ipv4"
)

var cmdRootParser = argparse.NewParser("mdns-proxy", "proxy local mdns service to remote")
var debugFlag = cmdRootParser.Flag("", "debug", &argparse.Options{Default: false, Help: "show debug info"})
var versionFlag = cmdRootParser.Flag("", "version", &argparse.Options{})

var runCmd = cmdRootParser.NewCommand("run", "run proxy")

func main() {
	var err = cmdRootParser.Parse(os.Args)
	// 如果解析失败就打印用法
	if err != nil {
		fmt.Println(cmdRootParser.Usage(err))
		return
	}

	// debug输出
	if *debugFlag {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// 版本信息
	if *versionFlag {
		log.Infof("current version: %s\n", Version)
		log.Infof("git tag: %s\n", GitTag)
		log.Infof("build time: %s\n", BuildTime)
		log.Infof("go version: %s\n", GoVersion)
		log.Infof("commit hash: %s\n", CommitHash)
		return
	}
}

func resolveMDnsHostname(hostname string) (string, error) {
	addr, err := net.ResolveUDPAddr("udp", mdns.DefaultAddress)
	if err != nil {
		return "", errors.Wrapf(err, "resolve udp addr %s failed ", mdns.DefaultAddress)
	}

	l, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return "", errors.Wrapf(err, "listen to udp4 %s failed ", addr)
	}

	server, err := mdns.Server(ipv4.NewPacketConn(l), &mdns.Config{})
	if err != nil {
		return "", errors.Wrap(err, "init mdns server failed ")
	}
	answer, src, err := server.Query(context.TODO(), hostname)
	log.Debugf("mdns query answer: %v", answer.GoString())
	if err != nil {
		return "", errors.Wrap(err, "query mdns failed ")
	}
	// fmt.Printf("answer: %s\n", answer.GoString())
	// fmt.Printf("src: %s\n", src.(*net.UDPAddr).IP.String())
	// fmt.Printf("err: %s\n", err)

	return src.(*net.UDPAddr).IP.String(), nil
}

func handler(writer dns.ResponseWriter, req *dns.Msg) {
	log.Debug("now handle dns request, %v", req)
	var resp dns.Msg
	resp.SetReply(req)
	for _, question := range req.Question {
		recordA := dns.A{
			Hdr: dns.RR_Header{
				Name:   question.Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			A: net.ParseIP("127.0.0.1"),
		}
		resp.Answer = append(resp.Answer, &recordA)
	}
	err := writer.WriteMsg(&resp)
	if err != nil {
		return
	}
}
