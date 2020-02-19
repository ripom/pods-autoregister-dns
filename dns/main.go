package main

import (
	"os/signal"
	"syscall"
        "os"
	"net"
	"dnsrecord" // package that include the selection of the Dnsprovider context and the Login, Create and Delete function to the dedicated Dnsprovider
	"strings"
	"log"
)

var (
	//Read the DnsProvider environment parameter
	dnsprovider	= os.Getenv("DnsProvider")
)


func main() {
	// Get the POD's IP address with subnet (example of IP address 172.16.0.4/16)
	addrs, _ := net.InterfaceAddrs()

	// Remove the subnet and get the IP (example 172.16.0.4)
	ipaddress:=strings.Split(addrs[1].String(),"/")[0]

	//Get the POD's hostname
	hostname, _ := os.Hostname()

	// Log some info
	log.Printf("DnsProvider: %s\n", dnsprovider)
	log.Printf("IP: %s\n", ipaddress)
	log.Printf("Hostname: %s\n", hostname)

	// Create the status, storing IP, Hostname and the Dnsprovider context
	status := dnsrecord.Status {}
	status.IP = ipaddress
	status.Hostname = hostname
	status.Dnsprovider = dnsprovider

	// Create the record in the Dns Provider resource
	if dnsrecord.CreateRecord(&status) == 0 {

		// Waiting for the exit signal
		sig := make(chan os.Signal, 1)
		signal.Notify(
			sig,
			syscall.SIGTERM,
			syscall.SIGINT,
		)
		defer signal.Stop(sig)
		<-sig

		// Remove the Dns record on exit
		dnsrecord.DeleteRecord(&status)
	}
}

