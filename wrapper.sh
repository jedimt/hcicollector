#!/usr/bin/env bash
while true
do
 /usr/bin/python /solidfire_graphite_collector_v2.py -s 172.27.40.200 -u admin -p solidfire -g 172.17.0.2
 /usr/bin/python /solidfire_graphite_collector_v2.py -s 172.27.40.205 -u admin -p solidfire -g 172.17.0.2 
sleep 60
done

