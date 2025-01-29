package utils

import (
	"fmt"
	"net"
)

/**
  @author: victor2022
  @since: 2025/1/29
*/

// GetLocalIps 获取本地ip地址
func GetLocalIps() []net.IP {
	addrs, err := net.InterfaceAddrs()
	result := make([]net.IP, 0)
	if err != nil {
		result = append(result, GetLocalIpByConn())
	} else {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					result = append(result, ipnet.IP)
				}
			}
		}
	}
	return result
}

func GetLocalIpByConn() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("无法获取 IP:", err)
		return net.IP{
			0, 0, 0, 0,
		}
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
