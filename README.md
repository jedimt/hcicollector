# SFCollector

SFcollector is a containerized collector for SolidFire clusters and is based off the following projects
* [solidfire-graphite-collector](https://github.com/cbiebers/solidfire-graphite-collector) Original Python collector script 
* [graphite-docker](https://github.com/jmreicha/graphite-docker) Graphite and Grafana containers
* [vsphere-graphite] (https://github.com/cblomart/vsphere-graphite) vSphere collector for Graphite

# Current Release
v .4 (beta)

## Updates in .4
* Added a vSphere collectored based heavily on the work of cblomart's vsphere-graphite collector
* Dashboard updates
* New dashboards for vSphere components 

# Description
The SFCollector is a fully packaged metrics collection and graphing solution for Element OS 8+ based on three container. 
* SFCollector -> runs a python script to scrape results from SolidFire clusters 
* Graphite database -> keeps all time series data from the SFCollector
* vsphere-graphite -> Pulls in vCenter metrics
* Grafana -> Graphing engine

The collector stores metrics in graphite and presents those metrics through a set of pre-configured Grafana dashboards.  Optionally, the Netapp Docker Volume Plugin (NDVP) can be used for persistent storage of metrics on a NetApp system.

![SFCollector architecture overview](http://www.jedimt.com/wp-content/uploads/2017/09/sfcollector-overview.jpeg)

### Prerequisites

* A Linux Docker host with docker-compose installed
* Optionally, the NetApp Docker Volume Plugin (NDVP)

### Installing

## Quick and Dirty Installation and Configuration

```
*(Optional) Install the NetApp Docker Volume Plugin (NDVP)
*Download and install docker-compose ('sudo pip install -U docker-compose)
*Clone this repo ('git clone https://github.com/jedimt/sfcollector')
*Modify bootstrap.sh script (`cd sfcollector && chmod +x bootstrap.sh`)
*Run the bootstrap.sh script (`./bootstrap.sh`)
*Modify the collector/wrapper.sh script supplying the SolidFire MVIP address,
and a user name and password
*Modify vsphere-graphite.json with your vCenter credentials and IP address 
*Modify docker-compose.yml to point at persistent storage volumes  
*Start up the containers (`docker-compose up`)
**Or in detached mode (`docker-compose up -d`)
*Add the graphite data source to Grafana
*Add the preconfigured Grafana dashboards from the 'dashboards' directory
```

A more complete installation and configuration guide "SolidFireStatsCollectionwithGraphiteandGrafana.pdf" is included in the repository.

## Author

**Aaron Patten**

*GitHub* - [Jedimt](https://github.com/jedimt)

*Blog* - [Jedimt.com](http://jedimt.com)

*Twitter* - [@jedimt](https://twitter.com/jedimt)

## Acknowledgments

* This would not have been possible if not for the prior work of cblocmart, cbiebers and jmreicha

