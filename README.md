# pods-autoregister-dns
This project is about a sidecar container to autoregister a Kubernetes POD IP address in DNS cloud provider

The container is written in GO and it's very small.
When it start, it gets the POD ip and then register an A record in the selected Cloud Provider DNS zone.
```DockerImage
Docker Image: riccardopomato/dnsregister

docker run --env DnsProvider='AzurePrivateDNS' -e subscription_id='xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx' -e ResourceGroup='resourcegroupname' -e DnsZone='example.local' -e TENANT_ID='xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx' -e CLIENT='xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx' -e KEY='xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx' riccardopomato/dnsregister 

docker run --env DnsProvider='AzureDNS' -e subscription_id='xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx' -e ResourceGroup='resourcegroupname' -e DnsZone='example.local' -e TENANT_ID='xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx' -e CLIENT='xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx' -e KEY='xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx' riccardopomato/dnsregister

docker run --env DnsProvider='IBMCloudDNS' -e IBM_DnsZone_ID='xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx' -e IBM_DnsZone_Instance_ID='xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx' -e IBM_APIKEY='xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx' dnsregister
```

This container support only Azure PrivateDNSZone, Azure DnsZone and IBM Cloud DNS Service.
I am planning to add more support to different cloud like Google, AWS and more.
Please, if you're interested, feel free to change the code or add more cloud support using pull request.

The mandatory ENV variables for AzureDNS or AzurePrivateDNS are:
```Azure Env Variables
DnsProvider: "AzureDNS" or "AzurePrivateDNS"
subscription_id: "xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx"
ResourceGroup: "resourcegroupname"
DnsZone: "domainname.com"
TENANT_ID: "xxxxxxxx-xxxx-xxxx-xxxxxxxxxxxxxxxxx"
CLIENT: "xxxxxxxx-xxxx-xxxx-xxxxxxxxxxxxxxxxx"
KEY: "xxxxxxxx-xxxx-xxxx-xxxxxxxxxxxxxxxxx"

DnsProvider can have only three values (case sensitive): AzureDNS or AzurePrivateDNS.
subscription_id: is referring to the Azure subscrption that host the DNS zone
ResourceGroup: is referring to the Azure Resource Group that contain the DNS zone
DnsZone: name of the Azure (public/private) DNS zone
TENANT_ID: is the Azure tenant where we need to log in using the service principal
CLIENT: is the Azure Service Principal AppId
KEY: is the Azure Service Principal secret
```

The mandatory ENV variables for IBMCloudDNS are:
```Ibm Env Variables
DnsProvider: "IBMCloudDNS"
IBM_DnsZone_ID: "xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx"
IBM_DnsZone_Instance_ID: "xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx"
IBM_APIKEY: "xxxxxxxx-xxxx-xxxxx-xxxxxxxxxxxxxxxxx"

DnsProvider can have only three values (case sensitive): IBMCloudDNS.
IBM_DnsZone_ID contains the DnsZone ID, you can collect this ID using IBM CLI.
IBM_DnsZone_Instance_ID contains the DnsZone Instance ID, you can collect this ID using IBM CLI.
IBM_APIKEY contains the APIKEY created on IBM Cloud and that has got the permission to create and delete record in the DNS Zone
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
Be aware, the yaml file use the configmap reference, then you should customize the configmap file with the correct ENV variable that point to your subscription and to an existing DNS zone.

```Steps
The steps to test the DNS record creation using kubernetes:
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
