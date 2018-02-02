# Changelog

# Current Release
.v5 (beta)

## Changes for Current Release
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
