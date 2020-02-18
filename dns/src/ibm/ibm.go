package ibm

import (
    "os"
    "time"
    "bytes"
    "encoding/json"
    "log"
    "io/ioutil"
    "net/http"
    "net/url"
)

var (
	apikey			= os.Getenv("IBM_APIKEY")
	dnssvcs_endpoint	= "https://api.dns-svcs.cloud.ibm.com"
	ibm_dns_instance_id     = os.Getenv("IBM_DnsZone_Instance_ID")
	ibm_dnszone_id		= os.Getenv("IBM_DnsZone_ID")
)

type Status struct {
	token		string
	record_id	string
}


var f interface{}

func Login(s *Status) int {
    data := url.Values{}
    data.Set("grant_type", `urn:ibm:params:oauth:grant-type:apikey`)
    data.Add("apikey", apikey)
    client := http.Client{Timeout: 30 * time.Second}
    req, err := http.NewRequest("POST", "https://iam.cloud.ibm.com/identity/token",  bytes.NewBufferString(data.Encode()))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Accept", "application/json")
    resp, err := client.Do(req)
    log.Printf("IbmCloud Logging in...")
    if err != nil {
	log.Println(err)
	return -1
    }else{
	    defer resp.Body.Close()
	    body, _ := ioutil.ReadAll(resp.Body)
	    if json.Unmarshal(body,&f) == nil {
			m := f.(map[string]interface{})
			if m["access_token"] == nil {
				log.Printf(m["errorMessage"].(string))
				log.Printf("Cannot log in, please check your APIKEY\n")
				return -1
			}else{
				(*s).token = m["access_token"].(string)
				log.Printf("IbmCloud Logged in!")
			}
	    }
	    return 0
    }
}

func Create(s *Status, ip string, hostname string) int {

	message := map[string]interface{}{
		"name": hostname,
		"type":  "A",
		"rdata": map[string]string{
			"ip": ip,
			},
		"ttl": 300,
	}
	bytesRepresentation, err := json.Marshal(message)

//	log.Printf("Token Create: %s\n",t.token)
	uri := dnssvcs_endpoint + "/v1/instances/" + ibm_dns_instance_id + "/dnszones/" + ibm_dnszone_id + "/resource_records"
	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", uri,  bytes.NewBuffer(bytesRepresentation) )
	req.Header.Set("Authorization", s.token)
	resp, err := client.Do(req)
	log.Printf("Creating Record...\n")
	if err != nil {
		log.Println(err)
		return -1
	}else{
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		if json.Unmarshal(body,&f) == nil {
			m := f.(map[string]interface{})
			if m["id"] == nil {
		                log.Println(m)
		                log.Printf("Cannot create record, please check your DNSZONE_ID or DNS_INSTANCE_ID\n")
		                return -1
		        }else{
				(*s).record_id = m["id"].(string)
				bodyprettyjson, _ := json.MarshalIndent(f,"","	")
				log.Printf("body: %s\n",bodyprettyjson)
				log.Printf("record_id %s\n", (*s).record_id)
				log.Printf("Record %s created\n", hostname)
				return 0
			}
		}else {
			return -1
		}
	}
}


func Delete(s *Status) int {
//	log.Printf("Token Create: %s\n",t.token)
//	log.Printf("record_id %s\n", t.record_id)
        uri := dnssvcs_endpoint + "/v1/instances/" + ibm_dns_instance_id + "/dnszones/" + ibm_dnszone_id + "/resource_records/" + s.record_id
        client := http.Client{Timeout: 30 * time.Second}
        req, err := http.NewRequest("DELETE", uri, nil)
        req.Header.Set("Authorization", s.token)
	req.Header.Set("Content-Type", "application/json")
        resp , err := client.Do(req)
        log.Printf("Deleting Record...\n")
        if err != nil {
	        log.Println(err)
	        return -1
	}else{
		defer resp.Body.Close()
	        log.Printf("Record %s deleted\n", s.record_id)
	        return 0
	}

}
