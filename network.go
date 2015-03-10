package gosync

import (
    "net"
    "log"
    "strings"
)

func GetLocalIp() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        log.Fatalf("Could not find network interfaces on this host: %s", err.Error())
    }
    var ips []string
    for _, a := range addrs {
        if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                ips = append(ips, ipnet.IP.String())
            }
        }
    }
    return strings.Join(ips, ",")
}