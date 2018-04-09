####Trident install and configuration####
#Create the Trident config file
mkdir -p /etc/netappdvp

cat << EOF > /etc/netappdvp/config.json
{
    "version": 1,
    "storageDriverName": "solidfire-san",
    "Endpoint": "https://admin:solidfire@10.193.136.240/json-rpc/9.0",
    "SVIP": "10.193.137.240:3260",
    "TenantName": "docker",
    "InitiatorIFace": "default",
    "Types": [
        {
            "Type": "docker-default",
            "Qos": {
                "minIOPS": 1000,
                "maxIOPS": 2000,
                "burstIOPS": 4000
            }
        },
        {
            "Type": "docker-app",
            "Qos": {
                "minIOPS": 4000,
                "maxIOPS": 6000,
                "burstIOPS": 8000
            }
        },
        {
            "Type": "docker-db",
            "Qos": {
                "minIOPS": 6000,
                "maxIOPS": 8000,
                "burstIOPS": 10000
            }
        }
    ]
}
EOF

#Install the Triedent plugin
docker plugin install --grant-all-permissions --alias netapp netapp/trident-plugin:18.01 config=config.json

#Create the Docker volume for the Graphite database
docker volume create -d netapp --name graphite-db -o type=docker-db -o size=50G
