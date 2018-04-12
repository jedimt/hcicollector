#Install HCI Collector####

for i in {16..21} {21..16} {16..21} {21..16} {16..21} {21..16} ; do echo -en "\e[38;5;${i}m#\e[0m" ; done ; echo

#Set some bash colors
Red='\033[0;31m'          # Red
Green='\033[0;32m'        # Green
Yellow='\033[0;33m'       # Yellow
Blue='\033[0;34m'         # Blue
Purple='\033[0;35m'       # Purple
Cyan='\033[0;36m'         # Cyan
White='\033[0;37m'        #White

####Trident install and configuration####
#Read in variables
echo -e ${Red} "Enter the SolidFire management virtual IP (MVIP): "
read SFMVIP
echo -e ${Red} "Enter the SolidFire storage virtual IP (SVIP): "
read SFSVIP
echo -e ${Red} "Enter the SolidFire username (case sensitive): "
read SFUSER
echo -e ${Red} "Enter the Solidfire password: "
read -s SFPASSWORD
echo -e ${Blue} "Enter the tenant account to use for Trident: "
read TACCOUNT
echo -e ${Blue} "Enter the volume name to create for Graphite"
read GRAPHITEVOL
echo -e ${Yellow} "Enter the password to use for the Grafana admin account: "
read -s GPASSWORD
echo -e ${Green} "Enter the vCenter username: "
read VCENTERUSER
echo -e ${Green} "Enter the vCenter password: "
read -s VCENTERPASSWORD
echo -e ${Green} "Enter the vCenter hostname. Ex. vcsa: "
read VCENTERHOSTNAME
echo -e ${Green} "Enter the vCenter domain. Ex. rtp.openenglab.netapp.com: "
read VCENTERDOMAIN
echo -e ${Purple} "Enter the IP address of this Docker host: "
read DOCKERIP
echo -e ${White} "Beginning Install"

#Create the Trident config file
mkdir -p /etc/netappdvp

cat << EOF > /etc/netappdvp/config.json
{
    "version": 1,
    "storageDriverName": "solidfire-san",
    "Endpoint": "https://$SFUSER:$SFPASSWORD@$SFMVIP/json-rpc/9.0",
    "SVIP": "$SFSVIP:3260",
    "TenantName": "$TACCOUNT",
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

echo "Installing Trident and creating the graphite-db volume"
#Install the Triedent plugin
docker plugin install --grant-all-permissions --alias netapp netapp/trident-plugin:18.01 config=config.json

#Create the Docker volume for the Graphite database
docker volume create -d netapp --name $GRAPHITEVOL -o type=docker-db -o size=50G

#Dccker compose configuration
echo "Creating the docker-compose.yml file"
cat << EOF > /opt/github/sfcollector/docker-compose.yml
version: "2"
services:
  graphite:
    build: ./graphiteconfig
    container_name: graphite-v.6
    restart: always
    ports:
        - "8080:80"
        - "8125:8125/udp"
        - "8126:8126"
        - "2003:2003"
        - "2004:2004"
    volumes: #Trident or local volumes for persistent storage
        - $GRAPHITEVOL:/opt/graphite/storage/whisper
    networks:
        - net_sfcollector

  grafana:
    build: ./grafana
    container_name: grafana-v.6
    restart: always
    ports:
        - "80:3000"
    networks:
        - net_sfcollector
    environment:
        #Set password for Grafana web interface
        - GF_SECURITY_ADMIN_PASSWORD=$GPASSWORD
        #Optional SMTP configuration for alert queries
        #- GF_SMTP_ENABLED=true
        #- GF_SMTP_HOST=smtp.gmail.com:465
        #- GF_SMTP_USER=<email address>
        #- GF_SMTP_PASSWORD=<email password>
        #- GF_SMTP_SKIP_VERIFY=true

  sfcollector-alpine:
    build: ./collector-alpine
    container_name: sfcollector-v.6
    restart: always
    networks:
        - net_sfcollector

  vsphere-collector:
    build: ./vsphere-graphite
    container_name: vsphere-graphite-v.6
    restart: always
    networks:
        - net_sfcollector
    depends_on:
        - graphite

networks:
  net_sfcollector:
    driver: bridge

volumes:
  $GRAPHITEVOL:
    external: true
EOF

#Wrapper script for the SolidFire collector
echo "Creating the SolidFire collector wrapper.sh script"
cat << EOF > /opt/github/sfcollector/collector-alpine/wrapper.sh
#!/usr/bin/env bash
while true
do
/usr/bin/python /solidfire_graphite_collector.py -s $SFMVIP -u $SFUSER -p $SFPASSWORD -g graphite &
sleep 60
done
EOF

#Make the file executiable
echo -e ${Cyan} "Marking wrapper.sh as executable"
chmod a+x /opt/github/sfcollector/collector-alpine/wrapper.sh


#Create the storage-schemas.conf file for Graphite
echo "Creating the storage-schemas.conf file"
cat << EOF > /opt/github/sfcollector/graphiteconfig/storage-schemas.conf
[stats]
pattern = ^stats\.*$
retentions = 5s:1d,1m:7d

[netapp]
pattern = ^netapp\.*
retentions = 1m:7d,5m:28d,10m:1y

[vsphere]
pattern = ^vsphere\.*
retentions = 1m:7d,5m:28d,10m:1y

[carbon]
pattern = ^carbon\.*$
retentions = 5s:1d,1m:7d

[statsd_internal_counts]
pattern = ^stats_counts\.statsd.*$
retentions = 5s:1d,1m:7d

[statsd_internal]
pattern = ^statsd\..*$
retentions = 5s:1d,1m:7d

[statsd]
pattern = ^stats_counts\..*$
retentions = 5s:1d,1m:28d,1h:2y

[statsd_gauges_internal]
pattern = ^stats\.gauges\.statsd\..*$
retentions = 5s:1d,1m:7d

[statsd_gauges]
pattern = ^stats\.gauges\..*$
retentions = 5s:1d,1m:28d,1h:2y

[catchall]
pattern = ^.*
retentions = 1m:5d,1m:28d
EOF

#Create the vsphere-graphite.json file for the vSphere-Graphite collector
echo "Creating the vsphere-graphite.json file"
cat << EOF > /opt/github/sfcollector/vsphere-graphite/vsphere-graphite.json
{
  "Domain": ".$VCENTERDOMAIN",
  "Interval": 60,
  "FlushSize": 100,
  "VCenters": [
    { "Username": "$VCENTERUSER", "Password": "$VCENTERPASSWORD", "Hostname": "$VCENTERHOSTNAME" }
  ],
  "Backend": {
    "Type": "graphite",
    "Hostname": "graphite",
    "Port": 2003
  },
  "Metrics": [
    {
      "ObjectType": [ "VirtualMachine", "HostSystem" ],
      "Definition": [
        { "Metric": "cpu.usage.average", "Instances": "" },
        { "Metric": "cpu.usage.maximum", "Instances": "" },
        { "Metric": "cpu.usagemhz.average", "Instances": "" },
        { "Metric": "cpu.usagemhz.maximum", "Instances": "" },
        { "Metric": "cpu.totalCapacity.average", "Instances": "" },
        { "Metric": "cpu.ready.summation", "Instances": "" },
        { "Metric": "mem.usage.average", "Instances": "" },
        { "Metric": "mem.usage.maximum", "Instances": "" },
        { "Metric": "mem.consumed.average", "Instances": "" },
        { "Metric": "mem.consumed.maximum", "Instances": "" },
        { "Metric": "mem.active.average", "Instances": "" },
        { "Metric": "mem.active.maximum", "Instances": "" },
        { "Metric": "mem.vmmemctl.average", "Instances": "" },
        { "Metric": "mem.vmmemctl.maximum", "Instances": "" },
        { "Metric": "disk.commandsAveraged.average", "Instances": "*" },
        { "Metric": "mem.totalCapacity.average", "Instances": "" }
      ]
    },
    {
      "ObjectType": [ "VirtualMachine" ],
      "Definition": [
        { "Metric": "virtualDisk.totalWriteLatency.average", "Instances": "*" },
        { "Metric": "virtualDisk.totalReadLatency.average", "Instances": "*" },
        { "Metric": "virtualDisk.numberReadAveraged.average", "Instances": "*" },
        { "Metric": "virtualDisk.numberWriteAveraged.average", "Instances": "*" },
        { "Metric": "cpu.ready.summation", "Instance": ""}
      ]
    },
    {
      "ObjectType": [ "HostSystem" ],
      "Definition": [
        { "Metric": "disk.maxTotalLatency.latest", "Instances": "" },
        { "Metric": "disk.numberReadAveraged.average", "Instances": "*" },
        { "Metric": "disk.numberWriteAveraged.average", "Instances": "*" },
        { "Metric": "disk.deviceLatency.average", "Instances": "*" },
        { "Metric": "disk.deviceReadLatency.average", "Instances": "*" },
        { "Metric": "disk.deviceWriteLatency.average", "Instances": "*" },
        { "Metric": "disk.kernelLatency.average", "Instances": "*" },
        { "Metric": "disk.queueLatency.average", "Instances": "*" },
        { "Metric": "datastore.datastoreIops.average", "Instances": "*" },
        { "Metric": "datastore.datastoreMaxQueueDepth.latest", "Instances": "*" },
        { "Metric": "datastore.datastoreReadBytes.latest", "Instances": "*" },
        { "Metric": "datastore.datastoreReadIops.latest", "Instances": "*" },
        { "Metric": "datastore.datastoreWriteBytes.latest", "Instances": "*" },
        { "Metric": "datastore.datastoreWriteIops.latest", "Instances": "*" },
        { "Metric": "datastore.numberReadAveraged.average", "Instances": "*" },
        { "Metric": "datastore.numberWriteAveraged.average", "Instances": "*" },
        { "Metric": "datastore.read.average", "Instances": "*" },
        { "Metric": "datastore.totalReadLatency.average", "Instances": "*" },
        { "Metric": "datastore.totalWriteLatency.average", "Instances": "*" },
        { "Metric": "datastore.write.average", "Instances": "*" },
        { "Metric": "mem.state.latest", "Instances": "" }
      ]
    }
  ]
}
EOF

#Create the datasource.yml file for the dashboards
cat << EOF > /opt/github/sfcollector/grafana/provisioning/datasources/datasource.yml
apiVersion: 1
datasources:
  - name: $GRAPHITEVOL
    type: graphite
    access: proxy
    orgId: 1
    url: http://$DOCKERIP:8080
    isDefault: true
    version: 1
    editable: true
    basicAuth: false
EOF

#Change the "datasource": instance in the provisoined dashboards to match the GRAPHITEVOL
echo "Modifying the default 'datasource' values in the pre-packeged dashboards"
DASHBOARDS=$(ls grafana/dashboards/*.json)
sed -i '/-- Grafana --/b; s/\("datasource": "\).*\(".*$\)/\1'$GRAPHITEVOL'\2/g' $DASHBOARDS
