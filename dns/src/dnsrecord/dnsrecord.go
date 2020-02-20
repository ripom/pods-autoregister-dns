package dnsrecord

import (
	"log"
	"ibm"		// Package for IBM Cloud DNS Resource
	"azure"		// Package for Azure Cloud for DNZzone and PrivateDnszone
)

// Save the Status of the request, this struct include the Status of different Dnsprovider
type Status struct {
	IP		string
	Hostname	string
	Dnsprovider	string
	id		string	//Store IBM Record_ID
}
var ok bool
var token string

// Create Dns Record in the Cloud provider
// This function use a Status struct to store the status
func CreateRecord(s *Status) int {

	// Check with Dnsprovider has been invoked
	switch provider := s.Dnsprovider; provider {
	case "IBMCloudDNS":
		// Invoke ibm package to create a Dns record in IBM CLoud
		// Login to IBM Cloud and get the token
		token, ok = ibm.Login()
		if ok {
			// If Login is successful, Create a DNS record and get the Record_ID
			// The Record_ID (stored in s.id), is important during the Record Delete operation
			s.id, ok = ibm.CreateDnsRecord(token, s.IP, s.Hostname)
			if ok {
				// Record created successfully
				return 0
			}
		}
		return -1
	case "AzurePrivateDNS":
               // Invoke azure package to create a Private Dns record in Azure Cloud
	       // Login to Azure Cloud and get the token 
		token, ok = azure.Login()
                if ok {
			// If Login is successful, Create a Private DNS record
			if azure.CreatePrivateDnsRecord(token, s.IP, s.Hostname) {
				// Record created successfully
				return 0
			}
	        }
		return -1
	case "AzureDNS":
               // Invoke azure package to create a Public Dns record in Azure Cloud
               // Login to Azure Cloud and get the token
               token, ok = azure.Login()
               if ok {
			// If Login is successful, Create a Public DNS record
			if azure.CreateDnsRecord(token, s.IP, s.Hostname) {
				// Record created successfully
				return 0
			}
                }
		return -1
	default:
		// No correct Dnsprovider selected or incorrect value provided
		log.Printf("Error with DnsProvider, you have selected a DnsProvider that not exists!\n")
		return -1
	}
}

func DeleteRecord(s *Status) int {
	switch provider := s.Dnsprovider; provider {
	case "IBMCloudDNS":
		// Invoke ibm package to create a Dns record in IBM CLoud
		// Login to IBM Cloud and get the token
		token, ok = ibm.Login()
		if ok {
			// If Login is successful, Delete the Record_ID (s.id) Dns record
			if !(ibm.DeleteDnsRecord(token, s.id)) {
				return -1
			}
		}
		return 0
	case "AzurePrivateDNS":
		// Invoke azure package to create a Private Dns record in Azure Cloud
		// Login to Azure Cloud and get the token
		token, ok = azure.Login()
		if ok {
			// If Login is successful, Delete the Hostname Private Dns record 
	                if !(azure.DeletePrivateDnsRecord(token, s.Hostname)) {
		                return -1
			}
		}
                return 0
	case "AzureDNS":
		// Invoke azure package to create a Public Dns record in Azure Cloud
		// Login to Azure Cloud and get the token
		token, ok = azure.Login()
		if ok {
			// If Login is successful, Delete the Hostname Public Dns record
			if !(azure.DeleteDnsRecord(token, s.Hostname)) {
				return -1
			}
		}
		return 0
	default:
		// No correct Dnsprovider selected or incorrect value provided
		log.Printf("Error with DnsProvider, you have selected a DnsProvider that not exists!\n")
		return -1
	}
}
