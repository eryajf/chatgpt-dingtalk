package ops

import (
	"crypto/tls"
	"net"
	"regexp"
	"strings"
	"time"
)

// 域名信息
type DomainMsg struct {
	CreateDate string `json:"create_date"`
	ExpiryDate string `json:"expiry_date"`
	Registrar  string `json:"registrar"`
}

// GetDomainMsg 获取域名信息
func GetDomainMsg(domain string) (dm DomainMsg, err error) {
	var conn net.Conn
	conn, err = net.Dial("tcp", "whois.verisign-grs.com:43")
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(domain + "\r\n"))
	if err != nil {
		return
	}
	buf := make([]byte, 1024)
	var num int
	num, err = conn.Read(buf)
	if err != nil {
		return
	}
	response := string(buf[:num])
	re := regexp.MustCompile(`Creation Date: (.*)\n.*Expiry Date: (.*)\n.*Registrar: (.*)`)
	match := re.FindStringSubmatch(response)
	if len(match) > 3 {
		dm.CreateDate = strings.TrimSpace(strings.Split(match[1], "Creation Date:")[0])
		dm.ExpiryDate = strings.TrimSpace(strings.Split(match[2], "Expiry Date:")[0])
		dm.Registrar = strings.TrimSpace(strings.Split(match[3], "Registrar:")[0])
	}
	return
}

// GetDomainCertMsg 获取域名证书信息
func GetDomainCertMsg(domain string) (cm tls.ConnectionState, err error) {
	var conn net.Conn
	conn, err = net.DialTimeout("tcp", domain+":443", time.Second*10)
	if err != nil {
		return
	}
	defer conn.Close()
	tlsConn := tls.Client(conn, &tls.Config{
		ServerName: domain,
	})
	defer tlsConn.Close()
	err = tlsConn.Handshake()
	if err != nil {
		return
	}
	cm = tlsConn.ConnectionState()
	return
}
