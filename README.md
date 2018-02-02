# SFCollector

The SolidFire collector is a container based metrics collection and graphing solution for NetApp HCI and SolidFire systems running Element OS 9+

# Current Release
v .5 (beta)

## Updates in .5
* Extensive dashboard updates. Dashboards now available on [grafana.com](https://grafana.com/dashboards?search=HCI)
* Added additional metrics to collection
* Updated to Trident from NDVP for persistent storage 
See the changelog for updates in previous versions

# Description
The SFCollector is a fully packaged metrics collection and graphing solution for Element OS 9+ based on these containers: 
* SFCollector -> runs a python script to scrape results from SolidFire clusters 
* vsphere-graphite -> vSphere stats collector, written in Go
* Graphite database -> keeps all time series data from the SFCollector
* Grafana -> Graphing engine

The collector stores metrics in graphite and presents those metrics through a set of pre-configured Grafana dashboards.  Optionally, the Netapp [Trident](https://netapp.io/2018/01/26/one-container-integration/) project can be used for persistent storage of metrics on a NetApp system.

![SFCollector architecture overview](http://www.jedimt.com/wp-content/uploads/2017/09/sfcollector-overview.jpeg)

## Prerequisites
* Docker host running 17.03+ 
* Account information for vCenter (optional) and SolidFire components to collect against 

## Quick and Dirty Installation and Configuration

```
*(Optional) Install NetApp Trident
*Download and install docker-compose ('sudo pip install -U docker-compose)
*Clone this repo ('git clone https://github.com/jedimt/sfcollector')
*Modify bootstrap.sh script (`cd sfcollector && chmod +x bootstrap.sh`)
*Run the bootstrap.sh script (`./bootstrap.sh`)
*Modify the ./collector-alpine/wrapper.sh script supplying the SolidFire MVIP address,
and a user name and password
*Rename ./vsphere-collector/vsphere-graphite-example.json to vsphere-graphite.json and modify with your vCenter credentials and IP address 
*Modify docker-compose.yml to point at persistent storage volumes (either on docker host or via Trident)  
*Start up the containers (`docker-compose up`)
**Or in detached mode (`docker-compose up -d`)
*Add the graphite data source to Grafana
*Add the preconfigured Grafana dashboards from the 'dashboards' directory or from grafana.com
```

A more complete installation and configuration guide "SFCollector_Install_and_Configure.pdf" is included in the repository.

## Author

**Aaron Patten**

*GitHub* - [Jedimt](https://github.com/jedimt)

*Blog* - [Jedimt.com](http://jedimt.com)

*Twitter* - [@jedimt](https://twitter.com/jedimt)

## Acknowledgments

This would not have been possible if not for the prior work of cblomart, cbiebers and jmreicha
* [solidfire-graphite-collector](https://github.com/cbiebers/solidfire-graphite-collector) Original Python collector script 
* [graphite-docker](https://github.com/jmreicha/graphite-docker) Graphite and Grafana containers
* [vsphere-graphite](https://github.com/cblomart/vsphere-graphite) vSphere collector for Graphite
