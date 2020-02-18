package main

import (
	"os/signal"
	"syscall"
        "os"
	"net"
	"dnsrecord"
	"strings"
	"log"
)

var (
	dnsprovider	= os.Getenv("DnsProvider")
)


func main() {
	addrs, _ := net.InterfaceAddrs()
	ipaddress:=strings.Split(addrs[1].String(),"/")[0]
	hostname, _ := os.Hostname()
	log.Printf("DnsProvider: %s\n", dnsprovider)
	log.Printf("IP: %s\n", ipaddress)
	log.Printf("Hostname: %s\n", hostname)
	status := dnsrecord.Status {}
	status.IP = ipaddress
	status.Hostname = hostname
	status.Dnsprovider = dnsprovider
	if dnsrecord.CreateRecord(&status) == 0 {
		sig := make(chan os.Signal, 1)
		signal.Notify(
			sig,
			syscall.SIGTERM,
			syscall.SIGINT,
		)
		defer signal.Stop(sig)
		<-sig
		dnsrecord.DeleteRecord(&status)
	}
}

