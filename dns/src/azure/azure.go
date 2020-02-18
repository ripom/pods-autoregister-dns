package azure

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
	azure_client_id			= os.Getenv("CLIENT")
	azure_client_secret		= os.Getenv("KEY")
	uri				= "https://login.microsoftonline.com/"
	azure_subscription_id		= os.Getenv("subscription_id")
	azure_tenant_id			= os.Getenv("TENANT_ID")
	azure_rg			= os.Getenv("ResourceGroup")
	azure_resource			= "https://management.azure.com/"
	azure_privatednszone_name	= os.Getenv("PrivateDnsZone")
)

type Status struct {
	token		string
}


var f interface{}

func Login(s *Status) int {
    auth_endpoint := uri + azure_tenant_id + "/oauth2/token"
    data := url.Values{}
    data.Set("grant_type", `client_credentials`)
    data.Set("client_id", azure_client_id)
    data.Set("client_secret", azure_client_secret)
    data.Set("resource", azure_resource)
    client := http.Client{Timeout: 30 * time.Second}
    req, err := http.NewRequest("POST", auth_endpoint, bytes.NewBufferString(data.Encode()))
//    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//    req.Header.Set("Accept", "application/json")
    resp, err := client.Do(req)
    log.Printf("AzureCloud Logging in...")
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
				log.Printf("Cannot log in, please check your CLIET_ID and CLIENT_SECRET\n")
				return -1
			}else{
//				bodyprettyjson, _ := json.MarshalIndent(f,"","  ")
 //                               log.Printf("body: %s\n",bodyprettyjson)
				(*s).token = m["access_token"].(string)
				log.Printf("AzureCloud Logged in!")
			}
	    }
	    return 0
    }
}

func Create(s *Status, ip string, hostname string) int {

        message := map[string]interface{}{
		"properties": map[string]interface{}{
//			"metadata": map[string]string {
//				"key1": hostname,
//			},
	                "ttl": 300,
			"aRecords": []map[string]string{map[string]string{
				"ipv4Address": ip,
			}},

		},
	}
	bytesRepresentation, err := json.Marshal(message)
//	prettyjson, _ := json.MarshalIndent(message,"","  ")
//	log.Printf("Message: %s\n", prettyjson)
//	log.Printf("Token Create: %s\n",t.token)
	uri := azure_resource + "subscriptions/" + azure_subscription_id + "/resourceGroups/" + azure_rg + "/providers/Microsoft.Network/privateDnsZones/" + azure_privatednszone_name + "/A/" + hostname + "?api-version=2018-09-01"
	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("PUT", uri,  bytes.NewBuffer(bytesRepresentation) )
	req.Header.Set("Authorization", "Bearer " + s.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	log.Printf("Creating Record...\n")
	if err != nil {
		log.Println(err)
		return -1
	}else{
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
		if json.Unmarshal(body,&f) == nil {
			m := f.(map[string]interface{})
			if m["id"] == nil {
		                log.Println(m)
		                log.Printf("Cannot create record, please check your DNSZONE_ID or DNS_INSTANCE_ID\n")
		                return -1
		        }else{
//				(*s).record_id = m["id"].(string)
				bodyprettyjson, _ := json.MarshalIndent(f,"","	")
				log.Printf("body: %s\n",bodyprettyjson)
//				log.Printf("record_id %s\n", (*s).record_id)
				log.Printf("Record %s created\n", hostname)
				return 0
			}
		}else {
			log.Println("error su unmarshal")
			log.Println(err)
			return -1
		}
	}
}


func Delete(s *Status, hostname string) int {
//	log.Printf("Token Create: %s\n",t.token)
//	log.Printf("record_id %s\n", t.record_id)
	uri := azure_resource + "subscriptions/" + azure_subscription_id + "/resourceGroups/" + azure_rg + "/providers/Microsoft.Network/privateDnsZones/" + azure_privatednszone_name + "/A/" + hostname + "?api-version=2018-09-01"
        client := http.Client{Timeout: 30 * time.Second}
        req, err := http.NewRequest("DELETE", uri, nil)
        req.Header.Set("Authorization", "Bearer " + s.token)
	req.Header.Set("Content-Type", "application/json")
        resp , err := client.Do(req)
        log.Printf("Deleting Record...\n")
        if err != nil {
	        log.Println(err)
	        return -1
	}else{
		defer resp.Body.Close()
	        log.Printf("Record %s deleted\n", hostname)
	        return 0
	}

}
