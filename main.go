/**
 * Created by psyduck on 2022/12/14
 */

package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

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
var runPortInt = runCmd.Int("p", "port", &argparse.Options{Required: false, Help: "expose dns service port", Default: 58})
var runIPStr = runCmd.String("", "ip", &argparse.Options{Required: false, Help: "expose dns service ip", Default: "0.0.0.0"})

var mDnsConn *mdns.Conn

func main() {
	var err = cmdRootParser.Parse(os.Args)
	if err != nil {
		fmt.Println(cmdRootParser.Usage(err))
		return
	}

	// log output
	if *debugFlag {
		log.SetReportCaller(true)
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

	if runCmd.Happened() {
		err := initMDnsConn()
		if err != nil {
			log.Panic(err)
		}
		dns.HandleFunc(".", handler)
		exposedURL := fmt.Sprintf("%s:%d", *runIPStr, *runPortInt)
		err = dns.ListenAndServe(exposedURL, "udp", nil)
		if err != nil {
			log.Panic(err)
		}
	}
}

func initMDnsConn() error {
	addr, err := net.ResolveUDPAddr("udp", mdns.DefaultAddress)
	if err != nil {
		return errors.Wrapf(err, "resolve udp addr %s failed ", mdns.DefaultAddress)
	}

	l, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return errors.Wrapf(err, "listen to udp4 %s failed ", addr)
	}

	mDnsConn, err = mdns.Server(ipv4.NewPacketConn(l), &mdns.Config{})
	if err != nil {
		return errors.Wrap(err, "init mdns server failed ")
	}
	return nil
}

// resolveMDnsHostname
// https://github.com/pion/mdns/blob/master/examples/query/main.go
func resolveMDnsHostname(hostname string) (string, error) {
	hostname = strings.TrimSuffix(hostname, ".")
	log.Debugf("now resolve hostname %s via mdns", hostname)
	answer, src, err := mDnsConn.Query(context.TODO(), hostname)
	log.Debugf("mdns query answer: %v", answer.GoString())
	if err != nil {
		return "", errors.Wrap(err, "query mdns failed ")
	}
	// fmt.Printf("answer: %s\n", answer.GoString())
	// fmt.Printf("src: %s\n", src.(*net.UDPAddr).IP.String())
	// fmt.Printf("err: %s\n", err)

	// https://stackoverflow.com/questions/50428176/how-to-get-ip-and-port-from-net-addr-when-it-could-be-a-net-udpaddr-or-net-tcpad
	return src.(*net.UDPAddr).IP.String(), nil
}

func handler(writer dns.ResponseWriter, req *dns.Msg) {
	log.Debugf("now handle dns request: \n %v\n", req)
	var resp dns.Msg
	resp.SetReply(req)
	for _, question := range req.Question {
		log.Debugf("now handle dns question, %+v", question)
		hostIP, err := resolveMDnsHostname(question.Name)
		if err != nil {
			log.Errorf("resolve mdns hostname failed, %s", err)
			return
		}
		log.Debugf("resolved mdns hostname %s to %s", question.Name, hostIP)
		recordA := dns.A{
			Hdr: dns.RR_Header{
				Name:   question.Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			A: net.ParseIP(hostIP),
		}
		resp.Answer = append(resp.Answer, &recordA)
	}
	err := writer.WriteMsg(&resp)
	if err != nil {
		log.Errorf("write response failed, %s", err)
		return
	}
}
