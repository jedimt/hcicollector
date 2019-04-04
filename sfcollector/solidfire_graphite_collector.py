#!/usr/bin/python
# solidfire_graphite_collector_v3.py
#
# Version 1.0.3
# Author: Aaron Patten
# Original author: Colin Bieberstein
# Original contributors: Pablo Luis Zorzoli, Davide Obbi
#
# Copyright  2018 NetApp, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
import argparse
import time
import graphyte
from solidfire.factory import ElementFactory
import solidfire.common
import logging


def send_cluster_faults(sf_element_factory, prefix):
    """
    send active cluster fault counts by: warning, error, critical
    exclude best practices, and only include current faults
    """
    fault_list = sf_element_factory.list_cluster_faults(False,"current").to_json()['faults']
    group = {'critical':0, 'error':0, 'warning':0 }
    for d in fault_list:
        if d['severity'] not in group:
            group[ d['severity'] ] = 1
        else:
            group[ d['severity'] ] += 1

    if to_graphite:
        for key in group:
            graphyte.send(prefix + '.fault.' + key, to_num(group[key]))
    else:
        for key in group:
            LOG.warning('fault.' + key, str(to_num(group[key])))


def send_cluster_stats(sf_element_factory, prefix):
    """
    send a subset of GetClusterStats API call results to graphite.
    """
    metrics = ['clientQueueDepth', 'clusterUtilization', 'readOpsLastSample',
               'readBytesLastSample', 'writeOpsLastSample', 'writeBytesLastSample',
               'actualIOPS', 'latencyUSec', 'normalizedIOPS', 'readBytes',
               'readLatencyUSec', 'readOps', 'unalignedReads', 'unalignedWrites',
               'writeLatencyUSec', 'writeOps', 'writeBytes']

    cluster_stats_dict = sf_element_factory.get_cluster_stats().to_json()['clusterStats']

    clusterUtilizationDec = float(cluster_stats_dict['clusterUtilization'])
    clusterUtilizationScaled = clusterUtilizationDec

    if to_graphite:
        graphyte.send(prefix + '.clusterUtilizationScaled', clusterUtilizationScaled)

    for key in metrics:
        if to_graphite:
            graphyte.send(prefix + '.' + key, to_num(cluster_stats_dict[key]))
        else:
            LOG.warning(key + ' ' + str(cluster_stats_dict[key]))

def send_cluster_capacity(sf_element_factory, prefix):
    """
    send a subset of GetClusterCapacity API call results and derived metrics to graphite.
    """
    metrics = ['activeBlockSpace', 'activeSessions', 'averageIOPS',
               'clusterRecentIOSize', 'currentIOPS', 'maxIOPS',
               'maxOverProvisionableSpace', 'maxProvisionedSpace',
               'maxUsedMetadataSpace', 'maxUsedSpace', 'nonZeroBlocks',
               'peakActiveSessions', 'peakIOPS', 'provisionedSpace',
               'snapshotNonZeroBlocks', 'timestamp', 'totalOps',
               'uniqueBlocks', 'uniqueBlocksUsedSpace', 'usedMetadataSpace',
               'usedMetadataSpaceInSnapshots', 'usedSpace', 'zeroBlocks']

    result = sf_element_factory.get_cluster_capacity().to_json()['clusterCapacity']
    for key in metrics:
        if to_graphite:
            graphyte.send(prefix + '.' + key, to_num(result[key]))
        else:
            LOG.warning(key + ' ' + str(result[key]))

    # Calculate & send derived metrics
    non_zero_blocks = to_num(result['nonZeroBlocks'])
    zero_blocks = to_num(result['zeroBlocks'])
    unique_blocks = to_num(result['uniqueBlocks'])
    unique_blocks_used_space = to_num(result['uniqueBlocksUsedSpace'])
    if non_zero_blocks != 0:
        thin_factor = float((non_zero_blocks + zero_blocks)) / float(non_zero_blocks)
    else:
        thin_factor = 1
    if to_graphite:
        graphyte.send(prefix + '.thin_factor', thin_factor)
    else:
        LOG.warning(key + ' ' + str(result[key]))

    if unique_blocks != 0:
        dedupe_factor = float(non_zero_blocks) / float(unique_blocks)
    else:
        dedupe_factor = 1
    if to_graphite:
        graphyte.send(prefix + '.dedupe_factor', dedupe_factor)
    else:
        LOG.warning('dedupe_factor ' + str(dedupe_factor))
    if unique_blocks_used_space != 0:
        #compression_factor = (unique_blocks * 4096.0) / unique_blocks_used_space
        compression_factor = (unique_blocks * 4096.0) / (unique_blocks_used_space * .93)
    else:
        compression_factor = 1
    if to_graphite:
        graphyte.send(prefix + '.compression_factor', compression_factor)
    else:
        LOG.warning('compression_factor ' + str(compression_factor))

    efficiency_factor = thin_factor * dedupe_factor * compression_factor
    if to_graphite:
        graphyte.send(prefix + '.efficiency_factor', efficiency_factor)
    else:
        LOG.warning('efficiency_factor ' + str(efficiency_factor))


def send_node_stats(sf_element_factory, prefix):
    """
    send a subset of ListNodeStats API call results to graphite.
    Note:   Calls ListAllNodes to get node name to use in metric path.
    """
    metrics_list = ['cpu', 'usedMemory', 'networkUtilizationStorage',
                    'networkUtilizationCluster', 'cBytesOut', 'cBytesIn', 'sBytesOut',
                    'sBytesIn', 'mBytesOut', 'mBytesIn', 'readOps', 'writeOps']

    node_list = sf_element_factory.list_all_nodes().to_json()['nodes']
    nodeinfo_by_id = list_to_dict(node_list, key="nodeID")

    nodestats = sf_element_factory.list_node_stats().to_json()['nodeStats']['nodes']
    for ns_dict in nodestats:
        node_name = nodeinfo_by_id[ns_dict['nodeID']]['name']
        for key in metrics_list:
            if to_graphite:
                graphyte.send(prefix + '.' + node_name + '.' + key, to_num(ns_dict[key]))
            else:
                LOG.warning(node_name + ' ' + key + ' ' + str(ns_dict[key]))


def send_volume_stats(sf_element_factory, prefix):
    """
    send a subset of ListVolumeStatsByVolume results to graphite.
    Note: Calls ListVolumes to get volume names for use in metric path.
    """
    metrics_list = ['volumeSize', 'zeroBlocks', 'nonZeroBlocks', 'volumeUtilization',
                    'actualIOPS', 'averageIOPSize', 'throttle', 'burstIOPSCredit',
                    'clientQueueDepth', 'latencyUSec',
                    'writeBytes', 'writeOps', 'writeLatencyUSec', 'unalignedWrites',
                    'readBytes', 'readOps', 'readLatencyUSec', 'unalignedReads',
                    'readBytesLastSample', 'readOpsLastSample', 'writeBytesLastSample',
                    'writeOpsLastSample']

    volume_list = sf_element_factory.list_volumes(include_virtual_volumes=False).to_json()['volumes']
    volinfo_by_id = list_to_dict(volume_list, key="volumeID")

    volstats = sf_element_factory.list_volume_stats_by_volume(include_virtual_volumes=False).to_json()['volumeStats']
    for vs_dict in volstats:
        vol_name = volinfo_by_id[vs_dict['volumeID']]['name']
        vol_accountID = volinfo_by_id[vs_dict['volumeID']]['accountID']
        vol_accountName = sf_element_factory.get_account_by_id(vol_accountID).to_json()['account']['username']
        for key in metrics_list:
            if to_graphite:
                graphyte.send(prefix + '.accountID.' + str(vol_accountName) + \
                              '.volume.' + vol_name + '.' + key, to_num(vs_dict[key]))
            else:
                LOG.warning('accountID ' + str(vol_accountName) + \
                            ' volume ' + vol_name + ' ' + key + ' ' + str(vs_dict[key]))


def send_drive_stats(sf_element_factory, prefix):
    """
    calculates summary statistics about drives by status and type at both cluster
    and node levels and submits them to graphite.
    Calls ListDrives and ListAllNodes
    """
    # Cluster level stats
    drive_list = sf_element_factory.list_drives().to_json()['drives']
    for status in ['active', 'available', 'erasing', 'failed', 'removing']:
        value = count_if(drive_list, 'status', status)
        if to_graphite:
            graphyte.send(prefix + '.drives.status.' + status, value)
        else:
            LOG.warning('drives.status ' + status + ' ' + str(value))
    for type in ['volume', 'block', 'unknown']:
        value = count_if(drive_list, 'type', type)
        if to_graphite:
            graphyte.send(prefix + '.drives.type.' + type, value)
        else:
            LOG.warning('drives.type ' + type + ' ' + str(value))
    # Node level stats
    node_list = sf_element_factory.list_all_nodes().to_json()['nodes']
    nodeinfo_by_id = list_to_dict(node_list, key="nodeID")
    for node in nodeinfo_by_id:
        node_name = nodeinfo_by_id[node]['name']
        for status in ['active', 'available', 'erasing', 'failed', 'removing']:
            value = count_ifs(drive_list, 'status', status, 'nodeID', node)
            if to_graphite:
                graphyte.send(prefix + '.node.' + node_name + '.drives.status.' + status, value)
            else:
                LOG.warning('node ' + node_name + ' drives.status ' + status + ' ' + str(value))
        for drive_type in ['volume', 'block', 'unknown']:
            value = count_ifs(drive_list, 'type', drive_type, 'nodeID', node)
            if to_graphite:
                graphyte.send(prefix + '.node.' + node_name + '.drives.type.' + drive_type, value)
            else:
                LOG.warning('node ' + node_name + ' drives.type ' + drive_type + ' ' + str(value))


def list_to_dict(list_of_dicts, key):
    """
    pivots a list of dicts into a dict of dicts, using key.
    """
    x = dict((child[key], dict(child, index=index)) for (index, child) in \
             enumerate(list_of_dicts))
    return x


def count_if(my_list, key, value):
    """
    return number of records in my_list where key==value pair matches
    """
    counter = (1 for item in my_list if item.get(key) == value)
    return sum(counter)


def count_ifs(my_list, key, value, key2, value2):
    """
    return number of records in my_list where both key==value pairs matches
    ToDo:   convert to grab any number of key=value pairs
    """
    counter = (1 for item in my_list if ((item.get(key) == value) and \
                                         (item.get(key2) == value2)))
    return sum(counter)


def to_num(metric):
    """
    convert string to number (int or float)
    """
    x = 0
    try:
        x = float(metric)
    except ValueError:
        try:
            x = float(metric)
        except ValueError:
            x = float('NaN')
    finally:
        return x


# Parse commandline arguments
parser = argparse.ArgumentParser()
parser.add_argument('-s', '--solidfire',
                    help='hostname of SolidFire array from which metrics should be collected')
parser.add_argument('-u', '--username', default='admin',
                    help='username for SolidFire array. default admin')
parser.add_argument('-p', '--password', default='password',
                    help='password for SolidFire array. default password')
parser.add_argument('-o', '--timeout', default=15,
                    help='Timeout for SolidFire API calls to complete.')
parser.add_argument('-g', '--graphite', default='localhost',
                    help='hostname of Graphite server to send to. default localhost. "debug" sends metrics to logfile')
parser.add_argument('-t', '--port', type=int, default=2003,
                    help='port to send message to. default 2003. if the --graphite is set to debug can be omitted')
parser.add_argument('-m', '--metricroot', default='netapp.solidfire.cluster',
                    help='graphite metric root. default netapp.solidfire.cluster')
parser.add_argument('-l', '--logfile', 
                    help='logfile.')
args = parser.parse_args()

to_graphite = True
# Logger module configuration
LOG = logging.getLogger('solidfire_graphite_collector.py')
if args.logfile:
    logging.basicConfig(filename=args.logfile, level=logging.DEBUG, format='%(asctime)s %(message)s')
    LOG.warning("Starting Collector script as a daemon.  No console output possible.")
else:
    logging.basicConfig(level=logging.DEBUG, format='%(asctime)s %(message)s')

# Initialize graphyte sender
if args.graphite == "debug":
    LOG.warning("Starting collector in debug mode. All the metrics will be shipped to logfile")
    to_graphite = False
else:
    graphyte.init(args.graphite, port=args.port, prefix=args.metricroot)

LOG.info("Metrics Collection for array: {0}".format(args.solidfire))
try:
    sfe = ElementFactory.create(args.solidfire, args.username, args.password)
    sfe.timeout(args.timeout)
except solidfire.common.ApiServerError as e:
    LOG.warning("ApiServerError: {0}".format(str(e)))
    sfe = None
except Exception as e:
    LOG.warning("General Exception: {0}".format(str(e)))
    sfe = None

while sfe != None:
    try:
        cluster_name = sfe.get_cluster_info().to_json()['clusterInfo']['name']
        send_cluster_stats(sfe, cluster_name)
        send_cluster_faults(sfe, cluster_name)
        send_cluster_capacity(sfe, cluster_name)
        send_volume_stats(sfe, cluster_name)
        send_drive_stats(sfe, cluster_name)
        send_node_stats(sfe, cluster_name + '.node')
    except solidfire.common.ApiServerError as e:
        LOG.warning("ApiServerError: {0}".format(str(e)))
    except Exception as e:
        LOG.warning("General Exception: {0}".format(str(e)))
