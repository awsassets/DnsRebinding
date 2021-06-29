package core

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"0xdns-rebind/conf"
	"0xdns-rebind/core/random"

	"github.com/miekg/dns"
)

/*
* rebind cache
 */

type RebindCache struct {
	sm sync.Map
}

func (rc *RebindCache) Get(ip string) bool {
	_, ok := rc.sm.Load(ip)
	// 第一次
	if !ok {
		rc.sm.Store(ip, struct{}{})
		return false
	}
	// 第二次
	rc.sm.Delete(ip)
	return true
}

const (
	RebindTTL  = 0
	NsTTL      = 10 * 60
	DefaultTTL = 5 * 60
	XIPTTL     = 24 * 60 * 60
)

var (
	rebindingCache = new(RebindCache)
)

// NewDNSDog dns serve
func NewDNSDog(addr string) (*dns.Server, error) {
	dns.HandleFunc(conf.C.Domain.Main+".", handleRebindRequest)
	return &dns.Server{Addr: addr, Net: "udp"}, nil
}

func handleRebindRequest(w dns.ResponseWriter, req *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(req)
	m.Compress = false

	switch req.Opcode {
	case dns.OpcodeQuery:
		parseQuery(w, req, m)
	}
	if err := w.WriteMsg(m); err != nil {
		log.Println("[Dns] write message error: ", err)
	}
}

func parseQuery(w dns.ResponseWriter, req *dns.Msg, m *dns.Msg) {
	// get ip
	ip := cutStrings(w.RemoteAddr().String(), ":")

	for _, q := range m.Question {
		if q.Qclass != dns.ClassINET {
			dns.HandleFailed(w, req)
			return
		}

		switch q.Qtype {
		case dns.TypeANY:
			fallthrough
		case dns.TypeA:
			if a := giveAnswer(ip, q.Name, dns.TypeA, RebindTTL); a != nil {
				m.Answer = append(m.Answer, a)
			}
		case dns.TypeNS:
			if a := giveAnswer(ip, q.Name, dns.TypeNS, NsTTL); a != nil {
				m.Answer = append(m.Answer, a)
			}
		default:
			dns.HandleFailed(w, req)
		}
	}
}

func giveAnswer(ip, qName string, qType uint16, ttl int) dns.RR {
	if !strings.HasSuffix(qName, fmt.Sprintf("%s.", conf.C.Domain.Main)) {
		return nil
	}
	respIP := conf.C.Domain.IP
	// handler dns-rebinding
	if ttl == RebindTTL {
		if rebindingCache.Get(ip) { // 第二次
			respIP = conf.C.Domain.RebindIP
		}
	}
	rrHeader := dns.RR_Header{
		Name:   qName,
		Rrtype: qType,
		Class:  dns.ClassINET,
		Ttl:    uint32(ttl),
	}
	switch qType {
	case dns.TypeA:
		return &dns.A{Hdr: rrHeader, A: net.ParseIP(respIP)}
	case dns.TypeNS:
		return &dns.NS{Hdr: rrHeader, Ns: conf.C.Domain.NS[random.Int(0, len(conf.C.Domain.NS)-1)]}
	default:
		return nil
	}
}

// cutStrings 移除 `cut` 字符之后的内容
func cutStrings(s, cut string) string {
	if strings.Contains(s, cut) {
		return strings.Split(s, cut)[0]
	}
	return s
}
