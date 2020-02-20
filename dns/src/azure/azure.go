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
	azure_client_id			= os.Getenv("CLIENT")	// Get Service Principal Client_ID to authenticate against the Azure Cloud
	azure_client_secret		= os.Getenv("KEY")	// Get Service Principal SECRET to authenticate against the Azure Cloud 
	uri				= "https://login.microsoftonline.com/"
	azure_subscription_id		= os.Getenv("subscription_id")	// Get SUBSCRIPTION_ID where the DNS Resource is located 
	azure_tenant_id			= os.Getenv("TENANT_ID")	// Get TENANT_ID where the Service Principal is located
	azure_rg			= os.Getenv("ResourceGroup")	// Get Resource Group Name where the DNS Resource is located
	azure_resource			= "https://management.azure.com/"
	azure_privatednszone_name	= os.Getenv("PrivateDnsZone")	// Get the DNS Resource name
	azure_privatedns_apiversion	= "2018-09-01"
	azure_dns_apiversion		= "2018-05-01"
	azure_privatednszonetypename	= "privateDnsZones"
	azure_dnszonetypename		= "dnsZones"
)


var f interface{}

// Login function with no parameter (the info are collect from environment variables
// The Login function has two return value, one is the Bearer Token and another one is a boolean to indicate if the operation has succedeed
func Login() (string, bool) {

    var token string	//Declare token local variable
    token = ""		//Initialize the variable

    // Create Payload with using Client_ID and Client_Secret
    auth_endpoint := uri + azure_tenant_id + "/oauth2/token"
    data := url.Values{}
    data.Set("grant_type", `client_credentials`)
    data.Set("client_id", azure_client_id)
    data.Set("client_secret", azure_client_secret)
    data.Set("resource", azure_resource)
    client := http.Client{Timeout: 30 * time.Second}
    // Create the request
    req, err := http.NewRequest("POST", auth_endpoint, bytes.NewBufferString(data.Encode()))
    // Perform the request
    resp, err := client.Do(req)
    log.Printf("AzureCloud Logging in...")
    // If there is an error during the request, the error has logged and a false value has return to indicate unsuccessful operation
    if err != nil {
	log.Println(err)
	return token, false	// Return a false value to indicate unsuccessful operation
    }else{
	    defer resp.Body.Close()	// Close the response object using defer key work to indicate this command it has run at the end of the function (deferred)
	    body, _ := ioutil.ReadAll(resp.Body)	// Read the body contained in the response
	    if json.Unmarshal(body,&f) == nil {		// Unmarshal the the json response
			m := f.(map[string]interface{})	// Convert the json response in a map object
			if m["access_token"] == nil {	// If the json response doesn't contain access_key, it means an error is logged
				log.Printf(m["errorMessage"].(string))	// Error message logged
				log.Printf("Cannot log in, please check your CLIET_ID and CLIENT_SECRET\n")
				return token, false	// Return a false value to indicate unsuccessful operation
			}else{	// If the json response contain access_key
				token = m["access_token"].(string)	// Assign the bearer token to the token variable
				log.Printf("AzureCloud Logged in!")	// Return the token and a true value to indicate a successful operation 
			}
	    }
	    return token, true	// Return a false value to indicate unsuccessful operation
    }
}

func CreatePrivateDnsRecord(token string, ip string, hostname string) bool {
	return createRecord(token, ip, hostname,  azure_privatednszonetypename, azure_privatedns_apiversion)
}

func CreateDnsRecord(token string, ip string, hostname string) bool {
	return createRecord(token, ip, hostname,  azure_dnszonetypename, azure_dns_apiversion)
}

func createRecord(token string, ip string, hostname string, zonetype string, apiversion string) bool {

	// Build the payload for the Create request
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
	payload, err := json.Marshal(message)
	// Create the request
	uri := azure_resource + "subscriptions/" + azure_subscription_id + "/resourceGroups/" + azure_rg + "/providers/Microsoft.Network/" + zonetype + "/" + azure_privatednszone_name + "/A/" + hostname + "?api-version=" + apiversion
	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("PUT", uri,  bytes.NewBuffer(payload) )
	// Set the Autherization header passing the token
	req.Header.Set("Authorization", "Bearer " + token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	log.Printf("Creating Record...\n")
	// If there is an error during the request, the error has logged and a false value has return to indicate unsuccessful operation
	if err != nil {
		log.Println(err)
		return false	// Return a false value to indicate unsuccessful operation
	}else{
		defer resp.Body.Close()		// Close the response object using defer key work to indicate this command it has run at the end of the function (deferred)
		body, _ := ioutil.ReadAll(resp.Body)	 // Read the body contained in the response
		body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))	// The server is sending you a UTF-8 text string with a Byte Order Mark (BOM). The BOM identifies that the text is UTF-8 encoded, but it should be removed before decoding.
		if json.Unmarshal(body,&f) == nil {	// Unmarshal the the json response
			m := f.(map[string]interface{})	// Convert the json response in a map object
			if m["id"] == nil {		// If the json response doesn't contain access_key, it means an error is logged
		                log.Println(m)		// Error message logged
		                log.Printf("Cannot create record, please check your DNSZONE_ID or DNS_INSTANCE_ID\n")
		                return false		// Return a false value to indicate unsuccessful operation
		        }else{	// False record has return to indicate unsuccessful operation
				bodyprettyjson, _ := json.MarshalIndent(f,"","	")	// Convertto to a pretty json response
				log.Printf("body: %s\n",bodyprettyjson)
				log.Printf("Record %s created\n", hostname)
				return true	// Return a true value to indicate a successful operation
			}
		}
		log.Println("error su unmarshal")
		log.Println(err)
		return false	// Return a false value to indicate unsuccessful operation
	}
}

func DeletePrivateDnsRecord(token string, hostname string) bool {
        return deleteRecord(token, hostname,  azure_privatednszonetypename, azure_privatedns_apiversion)
}

func DeleteDnsRecord(token string, hostname string) bool {
	return deleteRecord(token, hostname,  azure_dnszonetypename, azure_dns_apiversion)
}


func deleteRecord(token string, hostname string, zonetype string, apiversion string) bool {
	// Create the request
	uri := azure_resource + "subscriptions/" + azure_subscription_id + "/resourceGroups/" + azure_rg + "/providers/Microsoft.Network/" + zonetype + "/" + azure_privatednszone_name + "/A/" + hostname + "?api-version=" + apiversion
        client := http.Client{Timeout: 30 * time.Second}
        req, err := http.NewRequest("DELETE", uri, nil)
	// Set the Autherization header passing the token
        req.Header.Set("Authorization", "Bearer " + token)
	req.Header.Set("Content-Type", "application/json")
        resp , err := client.Do(req)
        log.Printf("Deleting Record...\n")
	// If there is an error during the request, the error has logged and a false value has return to indicate unsuccessful operation
        if err != nil {
	        log.Println(err)
	        return false	// Return a false value to indicate unsuccessful operation
	}else{
		defer resp.Body.Close()	// Close the response object using defer key work to indicate this command it has run at the end of the function (deferred)
	        log.Printf("Record %s deleted\n", hostname)
		return true	// Return a true value to indicate a successful operation
	}

}
