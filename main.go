package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

func main() {
	go func() {
		log.Fatal(dns.ListenAndServe(":8053", "udp", dns.HandlerFunc(handle)))
	}()
	log.Fatal(dns.ListenAndServe(":8053", "tcp", dns.HandlerFunc(handle)))
}

func handle(w dns.ResponseWriter, msg *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(msg)
qs:
	for _, q := range msg.Question {
		switch q.Qtype {
		case dns.TypeA:
			rr := &dns.A{
				Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: q.Qclass, Ttl: 0},
				A:   net.ParseIP("127.0.0.1"),
			}
			m.Answer = append(m.Answer, rr)
		case dns.TypeTXT:
			name := strings.SplitN(q.Name, ".", 2)
			if len(name) != 2 {
				m.Rcode = dns.RcodeBadName
				break qs
			}
			idx := 0
			var err error
			domain := name[1]
			if strings.HasPrefix(name[0], "_spf") {
				numPart := strings.TrimPrefix(name[0], "_spf")
				idx, err = strconv.Atoi(numPart)
				if err != nil {
					m.Rcode = dns.RcodeBadName
					break qs
				}
			} else {
				domain = q.Name
			}

			txt := fmt.Sprintf("v=spf1 include:_spf%d.%s -all", idx+1, domain)
			rr := &dns.TXT{
				Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeTXT, Class: q.Qclass, Ttl: 0},
				Txt: []string{txt},
			}
			m.Answer = append(m.Answer, rr)
		default:
			m.Rcode = dns.RcodeNameError
			break qs
		}
	}
	m.Authoritative = true
	m.Compress = true
	w.WriteMsg(m)
}
