# SFCollector

The SolidFire collector is a container based metrics collection and graphing solution for NetApp HCI and SolidFire systems running Element OS 9+

# Current Release
v .6 (beta)

## Updates in .6
* Retooled for Grafana 5.0.0
* Dashboards and datasources are now automatically added through the new provisioning functionality in Grafana 5
* Removed the external volume for the Grafana container, only Graphite uses an (optional) external iSCSI volume for persistent data
* Added the ability to poll for active alerts in the "SolidFire Cluster" dashboard.
* Added support for email alerting based on SolidFire events. Note: alerting queries do not support templating variables so if you have multiple clusters you will need to use "*" for the cluster instance instead of the "$Cluster" variable. The net effect of this is that the alert pane will show alerts from ALL clusters instead of an individually selected cluster.
* New detailed install document
* Added a very basic installation script

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
*Clone the https://github.com/jedimt/sfcollector Github repo 
*Execute the install_hcicollector.sh script and provide the requested input
*Start up the containers (`docker-compose up`)
**Or in detached mode (`docker-compose up -d`)
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
