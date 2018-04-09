# Changelog

# Current Release
.v6 (beta)

## Changes for Current Release
* Retooled for Grafana 5.0.0
* Dashboards and datasources are now automatically added through the new provisioning functionality in Grafana 5
* Removed the external volume for the Grafana container, only Graphite uses an (optional) external iSCSI volume for persistent data
* Added the ability to poll for active alerts in the "SolidFire Cluster" dashboard. 
* Added support for email alerting based on SolidFire events. Note: alerting queries do not support templating variables so if you have multiple clusters you will need to use "*" for the cluster instance instead of the "$Cluster" variable. The net effect of this is that the alert pane will show alerts from ALL clusters instead of an individually selected cluster. 
* New detailed install document
* Added a very basic installation script

## Updates in .5
* Extensive dashboard updates. Dashboards now available on [grafana.com](https://grafana.com/dashboards?search=HCI)
* Added additional metrics to collection
* Updated to Trident from NDVP for persistent storage 

## Updates in .4
* Added a vSphere collectored based heavily on the work of cblomart's vsphere-graphite collector
* Dashboard updates
* New dashboards for vSphere components 

## Updates in .3
* Changed the collector container to Alpine which dramatically cut down container size and build time.
* Other minor changes

### Changes for .v2
* Added "&" in wrapper.sh script to make the collector calls async. Previously the script was waiting for the collector script to finish before continuing the loop. This caused the time between collections to stack which caused holes in the dataset. Now stats should be returned every minute.
* Changed graphs to use the summerize function for better accuracy.

### Changes for .v1
* Initial release
