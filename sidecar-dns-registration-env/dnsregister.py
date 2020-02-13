# DNSregister is a container that used to register a DNS A record in a DNS cloud provider

  
# Importing library 
import requests, json
import logging
import socket 
import time
import signal
from os import environ
from azure.mgmt.dns import DnsManagementClient
from azure.mgmt.privatedns import PrivateDnsManagementClient
from azure.common.credentials import ServicePrincipalCredentials

# Load the environment variable
subscription_id = environ.get('subscription_id')
dnsProvider = environ.get('DnsProvider')
rg = environ.get('ResourceGroup')
dnszone = environ.get('DnsZone')
DNSZONE_ID=environ.get('IBM_DnsZone_ID')
INSTANCE_ID=environ.get('IBM_DnsZone_Instance_ID')

# Set the logging detail level
logging.basicConfig(level=logging.INFO)
logging.info('DNS provider: %s', dnsProvider)

class GracefulKiller:
    kill_now = False
    signals = {
        signal.SIGINT: 'SIGINT',
        signal.SIGTERM: 'SIGTERM'
    }
    def __init__(self):
        signal.signal(signal.SIGINT, self.exit_gracefully)
        signal.signal(signal.SIGTERM, self.exit_gracefully)
    def exit_gracefully(self, signum, frame):
        (hostname,ip)=get_Host_name_IP()
        logging.info("Deregister the A %s record...", hostname)
#   Remove the A record based to the DNS cloud provider
        if dnsProvider=='AzureDNS':
            logging.info("AzureDNS removing record...")
            cred=connect2Azure()
            DeleteAzureDNS_record(cred, subscription_id, rg, dnszone, hostname)
            logging.info("Record DNS removed")
        elif dnsProvider=='AzurePrivateDNS':
            logging.info("AzurePrivateDNS removing record...")
            cred=connect2Azure()
            DeleteAzurePrivateDNS_record(cred, subscription_id, rg, dnszone, hostname)
            logging.info("Record DNS removed")
        self.kill_now = True

def connect2Azure():
    # Load the environment variable to determinate the security context

    # Tenant ID for your Azure subscription
    TENANT_ID = environ.get('TENANT_ID') 
    # Your service principal App ID
    CLIENT = environ.get('CLIENT') 
    # Your service principal password
    KEY = environ.get('KEY') 
    try:
        credentials = ServicePrincipalCredentials(
                client_id = CLIENT,
                    secret = KEY,
                        tenant = TENANT_ID)
        logging.info("Connection to Azure estabilished")
    except Exception as e:
        logging.error("Exception occurred during Azure connection!", exc_info=True)
    return (credentials)

def connect2IBM():
    KEY = environ.get('IBM_APIKEY')
    data = {'grant_type':'urn:ibm:params:oauth:grant-type:apikey', 'apikey': KEY}
    payload = json.dumps(data)
    headers={'Content-Type': 'application/x-www-form-urlencoded', 'Accept': 'application/json'}
    try:
        response = requests.post("https://iam.cloud.ibm.com/identity/token", headers=headers, data=data)
        TOKEN=json.loads(response.text)["access_token"]
        logging.info("Connection to IBM Cloud estabilished")
    except Exception as e:
        logging.error("Exception occurred during IBM Cloud connection!", exc_info=True)
    return (TOKEN)

def CreateAzureDNS_record(credentials, subscription_id, rg, dnszone, hostname, ip):
    try:
        dns_client = DnsManagementClient(
                credentials,
                    subscription_id
                    )
        record_set = dns_client.record_sets.create_or_update(
            rg,
            dnszone,
            hostname,
            'A',
            {
                "ttl": 300,
                "arecords": [
                {
                    "ipv4_address": ip
                }]
            }
        )
        logging.info("AzureDNS record %s created", hostname)
    except Exception as e:
          logging.error("Exception occurred during AzureDNS record creation", exc_info=True)

def CreateAzurePrivateDNS_record(credentials, subscription_id, rg, dnszone, hostname, ip):
    try:
        dns_client = PrivateDnsManagementClient(
            credentials,
            subscription_id
            )
        record_set = dns_client.record_sets.create_or_update(
            rg,
            dnszone,
            'A',
            hostname,
            {
                "ttl": 300,
                "arecords": [
                {
                    "ipv4_address": ip
                }]
            }
        )
        logging.info("AzurePrivateDNS record %s created", hostname)
    except Exception as e:
        logging.error("Exception occurred during AzurePrivateDNS record creation", exc_info=True)

def DeleteAzurePrivateDNS_record(credentials, subscription_id, rg, dnszone, host):
    try:
        dns_client = PrivateDnsManagementClient(credentials, subscription_id)
        record_set = dns_client.record_sets.delete(rg, dnszone,'A', host)
        logging.info("AzurePrivateDNS record %s removed", host)
    except Exception as e:
        logging.error("Exception occurred during AzurePrivateDNS record deletion", exc_info=True)

def DeleteAzureDNS_record(credentials, subscription_id, rg, dnszone, host):
    try:
        dns_client = DnsManagementClient(credentials, subscription_id)
        record_set = dns_client.record_sets.delete(rg, dnszone, host, 'A')
        logging.info("AzureDNS record %s removed", host)
    except Exception as e:
        logging.error("Exception occurred during AzurePrivateDNS record deletion", exc_info=True)

# Function to get hostname and IP address 
def get_Host_name_IP(): 
    try: 
        host_name = socket.gethostname() 
        host_ip = socket.gethostbyname(host_name)
        logging.info("Getting pod info...")
        logging.info("Hostname : %s", host_name)
        logging.info("IP : %s", host_ip)
        return (host_name,host_ip)
    except Exception as e: 
        logging.error("Unable to get Hostname and IP", exc_info=True)

# Check which DNS provider has been selected
(hostname,ip)=get_Host_name_IP() 
if dnsProvider=='AzureDNS':
    cred=connect2Azure()
    CreateAzureDNS_record(cred, subscription_id, rg, dnszone, hostname, ip)
elif dnsProvider=='AzurePrivateDNS':
    cred=connect2Azure()
    CreateAzurePrivateDNS_record(cred, subscription_id, rg, dnszone, hostname, ip)
elif dnsProvider=='IBMCloudDNS':
    cred=connect2IBM()
    #CreateIBMCloudDNS_record(cred, subscription_id, rg, dnszone, hostname, ip)

if __name__ == '__main__':
    killer = GracefulKiller()
    logging.info("Running...")
    while not killer.kill_now:
        time.sleep(1)

