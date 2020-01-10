# pods-autoregister-dns
This project is about a sidecar container to autoregister a Kubernetes POD IP address in DNS cloud provider

The container is a python script.
When it start, it gets the POD ip and then register an A record in the DNS zone.
```DockerImage
Docker Image: riccardopomato/dnsregister
```

This container support only Azure Public DNS zone and Azure Private DNS Zone.
I am planning to add more support to different cloud like Google, AWS and more.
Please if you interested, feel free to change the code or add more cloud support using pull request.

The mandatory ENV variables are:
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

In the repository there is a test file DnsRegister.yaml useful to deploy a simple pod that automatically register its IP in the Azure DNS zone.
Be aware, the yaml file use the configmap reference, then you should customized the configmap file with the correct ENV variable that point to your subscription and to an existing DNS zone.

```Steps
The steps to test the DNS record creation:
1. Create a Service Principal
2. Create an Azure Private DNS, for example example.local
3. In the Access Control (IAM) of the DNS zone, assign the -Private DNS Zone Contributor- permission to the Service Principal create in the step 1
4. Customize the configmap-privatedns.yaml file with the correct details (subscription_id, ResourceGroup, DnsZone, TENANT_ID, CLIENT, KEY)
5. Create the configmap (kubectl create -f configmap-privatedns.yaml)
6. Create the POD (kubectl create -f DnsRegister.yaml)
7. Check the DNS zone, you should see a record with name DNSserver and with POD IP associated

The steps to test the DNS record deletion
1. Delete the POD (kubectl delete -f DnsRegister.yaml)
2. Check the DNS zone, you should see the record it has been removed

```
