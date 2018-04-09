#Install HCI Collector####
#Make configuration files

#Dccker compose configuration
cat << EOF > /opt/github/sfcollector/docker-compose.yml
version: "2"
services:
  graphite:
    build: ./graphiteconfig
    restart: always
    ports:
        - "8080:80"
        - "8125:8125/udp"
        - "8126:8126"
        - "2003:2003"
        - "2004:2004"
    volumes: #Trident or local volumes for persistent storage
        - graphite-db:/opt/graphite/storage/whisper
    networks:
        - net_sfcollector

  grafana:
    image: grafana/grafana
    restart: always
    ports:
        - "80:3000"
    networks:
        - net_sfcollector
    environment:
        #Set password for Grafana web interface
        - GF_SECURITY_ADMIN_PASSWORD=<your password>
        #Optional SMTP configuration for alert queries
        #- GF_SMTP_ENABLED=true
        #- GF_SMTP_HOST=smtp.gmail.com:465
        #- GF_SMTP_USER=<email address>
        #- GF_SMTP_PASSWORD=<email password>
        #- GF_SMTP_SKIP_VERIFY=true

  sfcollector-alpine:
    build: ./collector-alpine
    restart: always
    networks:
        - net_sfcollector

  vsphere-collector:
    build: ./vsphere-graphite
    restart: always
    networks:
        - net_sfcollector
    depends_on:
        - graphite

networks:
  net_sfcollector:
    driver: bridge

volumes:
  graphite-db:
    external: true
EOF

#Wrapper script for the SolidFire collector
cat << EOF > /opt/github/sfcollector/collector-alpine/wrapper.sh
#!/usr/bin/env bash
while true
do
/usr/bin/python /solidfire_graphite_collector.py -s 10.193.136.240 -u admin -p solidfire -g graphite &
sleep 60
done
EOF

#Create the storage-schemas.conf file for Graphite
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
cat << EOF > /opt/github/sfcollector/vsphere-graphite/vsphere-graphite.json
{
  "Domain": ".rtp.openenglab.netapp.com",
  "Interval": 60,
  "FlushSize": 100,
  "VCenters": [
    { "Username": "administrator@sflab.local", "Password": "<password>", "Hostname": "sfps-vcsa.rtp.openenglab.netapp.com" }
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
        { "Metric": "virtualDisk.totalWriteLatency.maximum", "Instances": "*" },
        { "Metric": "virtualDisk.totalReadLatency.average", "Instances": "*" },
        { "Metric": "virtualDisk.totalReadLatency.maximum", "Instances": "*" },
        { "Metric": "virtualDisk.numberReadAveraged.average", "Instances": "*" },
        { "Metric": "virtualDisk.numberWriteAveraged.average", "Instances": "*" },
        { "Metric": "cpu.read.summation", "Instance": ""}
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
