package main

import (
	"github.com/miekg/dns"
	"log"
	"net"
)

var cdnConfig = map[string]string{
	"web.cdn.test.c.lanpangzi.": "127.0.0.1",
}

// 处理到来的请求
func handler(writer dns.ResponseWriter, req *dns.Msg) {
	var resp dns.Msg
	resp.SetReply(req) // 创建应答
	for _, question := range req.Question {
		ip := cdnConfig[question.Name]
		if len(ip) == 0 {
			return
		}
		recordA := dns.A{
			Hdr: dns.RR_Header{
				Name:   question.Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    100,
			},
			A: net.ParseIP(ip).To4(), // 全部解析为127.0.0.1
		}
		resp.Answer = append(resp.Answer, &recordA) // 写入应答
	}
	err := writer.WriteMsg(&resp) // 回写信息
	if err != nil {
		return
	}
}

func main() {
	dns.HandleFunc(".", handler)                   // 绑定函数
	err := dns.ListenAndServe(":1053", "udp", nil) // 启动
	if err != nil {
		log.Println(err)
	}
}
