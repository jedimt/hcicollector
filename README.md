sfcollector
=============================

sfcollector is a containerized collector for SolidFire clusters and is based off the following projects
-> cbiebers/solidfire-graphite-collector (original collector logic)
-> jmreicha/graphite-docker (graphite + grafana)

Current Release
---------------

Version .1 (beta)

Description
-----------

The SolidFire collector is a fully packaged metrics collection and graphing solution 
for Element OS 8+ based on three container. 
	1. SFCollector-> runs a python script to scrape results from SolidFire clusters 
	2. Graphite database -> keeps all time series data from the SFCollector
	3. Grafana -> Graphing engine

The collector stores metrics in graphite and presents those metrics 
through a set of pre-configured Grafana dashboards.  Optionally, the Netapp Docker Volume
Plugin (NDVP) can be used for persistent storage of metrics on a NetApp system.

#Quick and Dirty Install and Configuration
 *(Optional) Install the NetApp Docker Volume Plugin (NDVP)
 *Download and install docker-compose ('sudo pip install -U docker-compose)
 *Clone this repo ('git clone https://github.com/jedimt/sfcollector')
 *Modify bootstrap.sh script (`cd sfcollector && chmod +x bootstrap.sh`)
 *Run the bootstrap.sh script (`./bootstrap.sh`)
 *Modify the collector/wrapper.sh script supplying the SolidFire MVIP address,
  and a user name and password
 *Modify docker-compose.yml to point at persistent storage volumes  
 *Start up the containers (`docker-compose up`)
 **Or in detached mode (`docker-compose up -d`)
 *Add the graphite data source to Grafana
 *Add the preconfigured Grafana dashboards from the 'dashboards' directory
 
A detailed install guide can be found at 
https://docs.google.com/document/d/1ZWiBs0_pYRTywlzlV0eV_Qnb9wwxH7u__dEUNQOCEFw/edit

**Details for the docker-compose.yml file**

In the `docker-compose.yml` file we are setting the default grafana password to
`P@ssw0rd`.  You can either modify the compose file update the account password
in the GUI after logging in the first time.

This set up uses basic_auth to secure graphite, you can view more info here -
http://nginx.org/en/docs/http/ngx_http_auth_basic_module.html.  To turn basic
auth off you can modify the Dockerfile and nginx config to remove references
the basic_auth settings.

The basic_auth generation relies on openssl for creating the user and
hash for the password so one of these tools must be installed for the
basic_auth component to work.

