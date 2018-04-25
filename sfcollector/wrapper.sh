#!/usr/bin/env bash
while true
do
/usr/bin/python /solidfire_graphite_collector.py -s 10.193.136.240 -u grafana -p netapp123 -g graphite &
sleep 60
done
