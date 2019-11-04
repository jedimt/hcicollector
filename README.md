# HCICollector

The HCI Collector is a container-based metrics collection and graphing solution for NetApp HCI and SolidFire systems running Element OS 9+

# Current Release
v .7 (beta)

## Updates in .7
* TODO

See the changelog for updates in previous versions

# Description
The SFCollector is a fully packaged metrics collection and graphing solution for Element OS 9+ based on these containers: 
* sfcollector -> runs a python script to scrape results from SolidFire clusters 
* vmwcollector -> vSphere stats collector, written in Go
* graphite -> keeps all time series data from the HCICollector
* grafana -> graphing engine

The collector stores metrics in graphite and presents those metrics through a set of pre-configured Grafana dashboards.  Optionally, the Netapp [Trident](https://netapp.io/2018/01/26/one-container-integration/) project can be used for persistent storage of metrics on a NetApp system.

![HCICollector architecture overview](https://github.com/jedimt/hcicollector/blob/master/hcicollector_architecture_overview.jpg)

## Prerequisites
* Docker host running 17.03+ 
* Account information for vCenter (optional) and SolidFire components to collect against 
 
## Quick and Dirty Installation and Configuration

```
# install Docker CE, docker-compose and iSCSI client packages
# enable and start docker, dnsmasq and open-iscsi service
# clone the hcicollector repository
git clone https://github.com/jedimt/hcicollector
# execute the install_hcicollector.sh script and provide the requested input 
cd hcicollector; sudo ./install_hcicollector.sh
# w/o Trident storage, create local volume using the volume name chosen in install wizard
#   docker volume create --name=chosen-name
# start up the containers
sudo docker-compose up
# to run in detached mode: sudo docker-compose up -d
```

For more information please consult the following material included in this repository:
* Installation and configuration guide "TR-4694-0618-Visualizing-NetApp-HCI-Performance.pdf" 
* Overview PowerPoint deck "HCICollector Overview.pptx"
* Installation demo and dashboard walkthrough "hcicollector-v.6-Install_and_Overview.mp4"

### Stand-Alone Use of Collector Script

Should you want to use collect Element cluster metrics in your own project or existing Graphite
environment, you may use solidfire_graphite_collector.py: `python3 solidfire_graphite_collector.py -h`

## Author

**Aaron Patten**

*GitHub* - [Jedimt](https://github.com/jedimt)

*Blog* - [Jedimt.com](http://jedimt.com)

*Twitter* - [@jedimt](https://twitter.com/jedimt)

## Acknowledgments

This would not have been possible if not for the prior work of cblomart, cbiebers and jmreicha
* [solidfire-graphite-collector](https://github.com/cbiebers/solidfire-graphite-collector) Original Python collector script 
* [docker-graphite-statsd](https://github.com/graphite-project/docker-graphite-statsd) Graphite and Statsd container
* [vsphere-graphite](https://github.com/cblomart/vsphere-graphite) vSphere collector for Graphite
