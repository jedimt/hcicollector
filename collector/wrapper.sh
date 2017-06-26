#!/usr/bin/env bash
#usage: python solidfire_graphite_collector.py [-h] [-s SOLIDFIRE] [-u USERNAME]
#usage: python solidfire_graphite_collector.py [-h] [-s SOLIDFIRE] [-u USERNAME]
#             [-p PASSWORD] [-g GRAPHITE] [-t PORT] [-m METRICROOT] [-l LOGFILE]
#
#  -h, --help            show this help message and exit
#
#  -s SOLIDFIRE, --solidfire SOLIDFIRE
#                    hostname of SolidFire array from which metrics should
#                    be collected
#
#  -u USERNAME, --username USERNAME
#                    username for SolidFire array.  default admin
#
#  -p PASSWORD, --password PASSWORD   
#                    password for SolidFire array.  default password
#
#  -g GRAPHITE, --graphite GRAPHITE
#                    hostname of Graphite server to send to.  default localhost
#
#  -t PORT, --port PORT  port to send message to.  default 2003
#
#  -m METRICROOT, --metricroot METRICROOT  port to send message to.  default netapp.solidfire.cluster
#
#  -l LOGFILE, --logfile LOGFILE  if defined, execution will be logged to this file.

while true
do
 /usr/bin/python /solidfire_graphite_collector_v2.py -s 172.27.40.200 -u admin -p solidfire -g graphite 
 #To monitor additional clusters add another line
 #/usr/bin/python /solidfire_graphite_collector_v2.py -s 172.27.40.205 -u admin -p solidfire -g graphite 
sleep 60
done

