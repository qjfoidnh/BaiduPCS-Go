package pcsutil

import (
	"net"
)

// ListAddresses 列出本地可用的 IP 地址
func ListAddresses() (addresses []string) {
	iFaces, _ := net.Interfaces()
	addresses = make([]string, 0, len(iFaces))
	for k := range iFaces {
		iFaceAddrs, _ := iFaces[k].Addrs()
		for l := range iFaceAddrs {
			switch v := iFaceAddrs[l].(type) {
			case *net.IPNet:
				addresses = append(addresses, v.IP.String())
			case *net.IPAddr:
				addresses = append(addresses, v.IP.String())
			}
		}
	}
	return
}

// ParseHost 解析地址中的host
func ParseHost(address string) string {
	h, _, err := net.SplitHostPort(address)
	if err != nil {
		return address
	}
	return h
}
