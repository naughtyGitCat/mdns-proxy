/**
 * Created by zhangruizhi on 2022/12/14
 */

package main

import (
	"github.com/miekg/dns"
	"net"
)

func handler(writer dns.ResponseWriter, req *dns.Msg) {
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
