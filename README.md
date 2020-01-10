# pods-autoregister-dns
This project is about a sidecar container to autoregister a Kubernetes POD IP address in DNS cloud provider

The container is a python script.
When it start, it gets the POD ip and then register an A record in the DNS zone using the POD ip.

This container support only Azure Public DNS zone and Azure Private DNS Zone.
I am planning to add more support to different cloud like Google, AWS and more.
Please if you interested, feel free to change the code or add more cloud support using pull request.

The ENV variables mandatory are:
```ENV Variables
DnsProvider: "AzureDNS"
subscription_id: "xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx"
ResourceGroup: "resourcegroupname"
DnsZone: "domainname.com"
TENANT_ID: "xxxxxxxx-xxxx-xxxx-xxxxxxxxxxxxxxxxx"
CLIENT: "xxxxxxxx-xxxx-xxxx-xxxxxxxxxxxxxxxxx"
KEY: "xxxxxxxx-xxxx-xxxx-xxxxxxxxxxxxxxxxx"

DnsProvider can have only two values (case sensitive): AzureDNS or AzurePrivateDNS
subscription_id: is referring to the subscrption that host the DNS zone
ResourceGroup: is referring to the Resource Group that contain the DNS zone
DnsZone: name of the DNS zone
TENANT_ID: is the tenant where we need to log in using the service principal
CLIENT: is the AppId of the Service Principal
KEY: is the Service Principal secret
```
In kubernetes deployment, is good to create a configmap and store the ENV variable like this example:

```Configmap
kind: ConfigMap
apiVersion: v1
metadata:
  name: appconfig
data:
        DnsProvider: "AzurePrivateDNS"
        subscription_id: "xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx"
        ResourceGroup: "resourcegroupname"
        DnsZone: "domainname.local"
        TENANT_ID: "xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx"
        CLIENT: "xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx"
        KEY: "xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx"
```
