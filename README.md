solidfire-graphite-collector
=============================

Project to collect and store SolidFire cluster metrics in Graphite.  


# solidfire-graphite-collector

Current Release
---------------

Version 1.0.5


Description
-----------

The SolidFire Graphite collector is a simple utility to collect metrics 
from Element OS 8.x and store them in graphite.   

Required Libraries
------------------

| Component                                                        | Version   |
|------------------------------------------------------------------|-----------|
| solidfire-sdk-python <https://github.com/solidfire/solidfire-sdk-python/> | 1.1   |
| Requests <http://docs.python-requests.org/en/master/>            | 2.12.1+   |
| graphyte <https://github.com/Jetsetter/graphyte/>                | 1.1       |
| python-daemon <https://pypi.python.org/pypi/python-daemon/>      | 2.1.2     |
| logging | 0.4.9.6 |


Usage
-----

The script is called from the command line using the parameters shown below.  
You can run multiple instances of the script in parallel if you have more than one 
cluster to monitor.


    usage: python solidfire_graphite_collector.py [-h] [-s SOLIDFIRE] [-u USERNAME]
                 [-p PASSWORD] [-g GRAPHITE] [-t PORT] [-m METRICROOT] [-l LOGFILE]

      -h, --help            show this help message and exit

      -s SOLIDFIRE, --solidfire SOLIDFIRE
                        hostname of SolidFire array from which metrics should
                        be collected

      -u USERNAME, --username USERNAME
                        username for SolidFire array.  default admin

      -p PASSWORD, --password PASSWORD   
                        password for SolidFire array.  default password

      -g GRAPHITE, --graphite GRAPHITE
                        hostname of Graphite server to send to.  default localhost

      -t PORT, --port PORT  port to send message to.  default 2003

      -m METRICROOT, --metricroot METRICROOT  port to send message to.  default netapp.solidfire.cluster

      -l LOGFILE, --logfile LOGFILE  if defined, execution will be logged to this file.



To have it automatically startup on server boot, make use of an rc.d script (or upstart) 
as appropriate for your OS version.   

To stop this script, simply kill the process.  A sample command to do so is shown below:

    ps -ef | grep solidfire_graphite_collector.py | grep -v grep | awk {'print $2'} \
    | xargs kill


Other Scripts
===============

#launcher.py 

Current Release
---------------

Version 1.0.1

Description
-----------

Helper script to use a configuration file to provide arguments for launching multiple 
instances of solidfire_graphite_collector.py at once.

| Component                                                        | Version   |
|------------------------------------------------------------------|-----------|
| python3                                                          | 3.x       |
| configparser | 3.5.0 |

Requires a configuration file (sgc.config) with one or more sections in the form:

        [solidfire_array_hostname]
        username : solidfire_user
        password : solidfire_pass
        graphite : graphite_server_hostname
        port     : graphite_port
        metricroot : netapp.solidfire.cluster

Where graphite and port are optional.

#dashboards/*

Description
-----------
The dashboards directorycontains a set of sample Grafana dashboards 
utilizing the collected metrics.



**License**
-----------

Copyright © 2016 NetApp, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
