package logdna

import (
	"bytes"
	"net"
	"os"
	"strings"
)

func getHostName() string {
	name, err := os.Hostname()

	if err != nil {
		return err.Error()
	}
	return name
}

func getMacAddr() string {
	addr := ""
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				addr = i.HardwareAddr.String()
				break
			}
		}
	}
	return addr
}
func getIpAddr() string {
	ipaddr := ""
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				addrs, err := i.Addrs()
				if err != nil {
					return err.Error()
				}
				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
					if ip == nil || ip.IsLoopback() {
						continue
					}
					ip = ip.To4()
					if ip == nil {
						continue
					}
					ipaddr = ip.String()
				}
				break
			}
		}
	}
	return ipaddr
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.ToLower(a) == strings.ToLower(e) {
			return true
		}
	}
	return false
}
