# Changelog

# Current Release
.v4 (beta)

## Changes for Current Release
* Added a vSphere collector heavily based on cblomarts vsphere-graphite project. 
* Added dashboards specific to vSphere and VM metrics (still a work in progress) 

## Updates in .3
* Changed the collector container to Alpine which dramatically cut down container size and build time.
* Other minor changes

### Changes for .v2
* Added "&" in wrapper.sh script to make the collector calls async. Previously the script was waiting for the collector script to finish before continuing the loop. This caused the time between collections to stack which caused holes in the dataset. Now stats should be returned every minute.
* Changed graphs to use the summerize function for better accuracy.

### Changes for .v1
* Initial release
