package dnsrecord

import (
	"log"
	"ibm"
	"azure"
)

type Status struct {
	IP		string
	Hostname	string
	Dnsprovider	string
	ibmstatus	ibm.Status
	azurestatus	azure.Status
}

func CreateRecord(s *Status) int {
	switch provider := s.Dnsprovider; provider {
	case "IBMCloudDNS":
		(*s).ibmstatus = ibm.Status{}
		if ibm.Login(&(s.ibmstatus)) == 0 {
			if ibm.Create(&(s.ibmstatus), s.IP, s.Hostname) != 0 {
				return -1
			}else{
				return 0
			}
		}
		return -1
	case "AzurePrivateDNS":
		(*s).azurestatus = azure.Status{}
                if azure.Login(&(s.azurestatus)) == 0 {
	                if azure.Create(&(s.azurestatus), s.IP, s.Hostname) != 0 {
				return -1
			}else{
				return 0
			}
	        }
		return -1
	case "AzureDNS":
               // 
		return -1
	default:
		log.Printf("Error with DnsProvider, you have selected a DnsProvider that not exists!\n")
		return -1
	}
}

func DeleteRecord(s *Status) int {
	switch provider := s.Dnsprovider; provider {
	case "IBMCloudDNS":
		if ibm.Delete(&(s.ibmstatus)) !=0 {
			return -1
		}
		return 0
	case "AzurePrivateDNS":
                if azure.Delete(&(s.azurestatus), s.Hostname) !=0 {
                        return -1
                }
                return 0
	case "AzureDNS":
		//
		return -1
	default:
		log.Printf("Error with DnsProvider, you have selected a DnsProvider that not exists!\n")
		return -1
	}
}
