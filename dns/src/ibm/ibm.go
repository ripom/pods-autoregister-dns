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
	apikey			= os.Getenv("IBM_APIKEY")			// Get IBM APIKEY to authenticate against the IBM Cloud
	dnssvcs_endpoint	= "https://api.dns-svcs.cloud.ibm.com"
	ibm_dns_instance_id     = os.Getenv("IBM_DnsZone_Instance_ID")		// Get the IBM Dns instance id of the DNS resource
	ibm_dnszone_id		= os.Getenv("IBM_DnsZone_ID")			// Get the IBM Dns zone id of the DNS resource
)

var f interface{}

// Login function with no parameter (the info are collect from environment variables
// The Login function has two return value, one is the Bearer Token and another one is a boolean to indicate if the operation has succedeed
func Login() (string, bool) {

    var token string	//Declare token local variable
    token =""		//Initialize the variable

    // Create Payload with two value grant_type and apikey
    data := url.Values{}
    data.Set("grant_type", `urn:ibm:params:oauth:grant-type:apikey`)
    data.Add("apikey", apikey)
    client := http.Client{Timeout: 30 * time.Second}
    // Create the request
    req, err := http.NewRequest("POST", "https://iam.cloud.ibm.com/identity/token",  bytes.NewBufferString(data.Encode()))
    // Set the request header	
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Accept", "application/json")
    // Perform the request
    resp, err := client.Do(req)
    log.Printf("IbmCloud Logging in...")
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
				log.Printf("Cannot log in, please check your APIKEY\n")
				return token, false	// Return a false value to indicate unsuccessful operation
			}else{	// If the json response contain access_key
				token = m["access_token"].(string)	// Assign the bearer token to the token variable
				log.Printf("IbmCloud Logged in!")
				return token, true  // Return the token and a true value to indicate a successful operation
			}
	    }
	    return token, false	// Return a false value to indicate unsuccessful operation
    }
}

func CreateDnsRecord(token string, ip string, hostname string) (string, bool) {

	var id string
	id = ""

	// Build the payload for the Create request
	message := map[string]interface{}{
		"name": hostname,
		"type":  "A",
		"rdata": map[string]string{
			"ip": ip,
			},
		"ttl": 300,
	}
	payload, err := json.Marshal(message)
	// Create the request
	uri := dnssvcs_endpoint + "/v1/instances/" + ibm_dns_instance_id + "/dnszones/" + ibm_dnszone_id + "/resource_records"
	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", uri,  bytes.NewBuffer(payload) )
	// Set the Autherization header passing the token
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	log.Printf("Creating Record...\n")
	// If there is an error during the request, the error has logged and a false value has return to indicate unsuccessful operation
	if err != nil {
		log.Println(err)
		return id, false	// Return a false value to indicate unsuccessful operation
	}else{
		defer resp.Body.Close()		// Close the response object using defer key work to indicate this command it has run at the end of the function (deferred)
		body, _ := ioutil.ReadAll(resp.Body)	// Read the body contained in the response
		if json.Unmarshal(body,&f) == nil {	// Unmarshal the the json response
			m := f.(map[string]interface{})	// Convert the json response in a map object
			if m["id"] == nil {		// If the json response doesn't contain access_key, it means an error is logged
		                log.Println(m)		// Error message logged
		                log.Printf("Cannot create record, please check your DNSZONE_ID or DNS_INSTANCE_ID\n")
		                return id, false	// Return a false value to indicate unsuccessful operation
		        }else{	// False record has return to indicate unsuccessful operation
				id := m["id"].(string)	// Assign the Record_ID to the id variable
				bodyprettyjson, _ := json.MarshalIndent(f,"","	")	// Convertto to a pretty json response
				log.Printf("body: %s\n",bodyprettyjson)			// Log the json response
				log.Printf("record_id %s\n", id)			// Log the Record_ID
				log.Printf("Record %s created\n", hostname)
				return id, true		// Return the token and a true value to indicate a successful operation
			}
		}
		return id, false        // Return a false value to indicate unsuccessful operation
	}
}

func DeleteDnsRecord(token string, id string) bool {
	// Create the request
        uri := dnssvcs_endpoint + "/v1/instances/" + ibm_dns_instance_id + "/dnszones/" + ibm_dnszone_id + "/resource_records/" + id
        client := http.Client{Timeout: 30 * time.Second}
        req, err := http.NewRequest("DELETE", uri, nil)
	// Set the Autherization header passing the token
        req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
        resp , err := client.Do(req)
        log.Printf("Deleting Record...\n")
        // If there is an error during the request, the error has logged and a false value has return to indicate unsuccessful operation
	if err != nil {
	        log.Println(err)
	        return false	// Return a false value to indicate unsuccessful operation
	}else{
		defer resp.Body.Close()		// Close the response object using defer key work to indicate this command it has run at the end of the function (deferred)  
	        log.Printf("Record %s deleted\n", id)
	        return true	// Return a true value to indicate a successful operation   
	}

}
