# vSphere to Graphite

Export vSphere stats to graphite.

Written in go as the integration collectd and python plugin posed too much problems (cpu usage and pipe flood).

## Build status

Travis: [![Travis Build Status](https://travis-ci.org/cblomart/vsphere-graphite.svg?branch=master)](https://travis-ci.org/cblomart/vsphere-graphite)

Drone: [![Drone Build Status](https://bot.blomart.net/api/badges/cblomart/vsphere-graphite/status.svg)](https://bot.blomart.net/cblomart/vsphere-graphite)

## Code report

[![Go Report Card](https://goreportcard.com/badge/github.com/cblomart/vsphere-graphite)](https://goreportcard.com/report/github.com/cblomart/vsphere-graphite)

## Example result

Bellow a dashboard realize with grafana.
The backend used in this case is influxdb.

![Example Dashboard](vsphere-graphite-influxdb-grafana-dashboard.png)

## Extenal dependencies

Naturaly heavilly based on [govmomi](https://github.com/vmware/govmomi).

But also on [daemon](github.com/takama/daemon) which provides simple daemon/service integration.

## Configure

You need to know you vcenters, logins and password ;-)

If you set a domain, it will be automaticaly removed from found objects.

Metrics collected are defined by associating ObjectType groups with Metric groups.
They are expressed via the vsphere scheme: *group*.*metric*.*rollup*

ObjectTypes are explained in [this](https://code.vmware.com/web/dp/explorer-apis?id=196) vSphere doc

Performance metrics are explained in [this](https://docs.vmware.com/en/VMware-vSphere/6.5/com.vmware.vsphere.monitoring.doc/GUID-E95BD7F2-72CF-4A1B-93DA-E4ABE20DD1CC.html) vSphere doc

An example of configuration file of contoso.com is [there](./vsphere-graphite-example.json).

You need to place it at /etc/*binaryname*.json (/etc/vsphere-graphite.json per default)

For contoso it would simply be:

  > cp vsphere-graphite-example.json vsphere-graphite.json

Backend paramters can also be set via environement paramterers (see docker)

### Backend parameters

- Type (BACKEND_TYPE): Type of backend to use. Currently "graphite", "influxdb" or "thinfluxdb" (influx client in the project) 

- Hostname (BACKEND_HOSTNAME): hostname were the backend is running (graphite, influxdb, thinfluxdb)

- Port (BACKEND_PORT): port to connect to for the backend (graphite, influxdb, thinfluxdb)

- Username (BACKEND_USERNAME): username to connect to the backend (influxdb and optionally for thinfluxdb)

- Password (BACKEND_PASSWORD): password to connect to the backend (influxdb and optionally for thinfluxdb)

- Database (BACKEND_DATABASE): database to use in the backend (influxdb, thinfluxdb)

- NoArray (BACKEND_NOARRAY): don't use csv 'array' as tags, only the first element is used (influxdb, thinfluxdb)

## Docker

All builds are pushed to docker:

- [cblomart/vsphere-graphite](https://hub.docker.com/r/cblomart/vsphere-graphite/)

- [cblomart/rpi-vsphere-graphite](https://hub.docker.com/r/cblomart/rpi-vsphere-graphite/)

Default tags includes:

- branch (i.e.: master) for latest commit in the branch

- latest for latest release

Configration file can be passed by mounting /etc.

Backend parameters can be set via environment variables to make docker user easier (having graphite or influx as another container).

## Run it

### the container way

Edit the configuration file and set it in the place you like here $(pwd)

  > docker run -t -v $(pwd)/vsphere-graphite.json:/etc/vsphere-graphite.json cblomart/vsphere-graphite:latest

### The old way

#### Deploy

Of course `golang` is needed. Install and set `$GOPATH` such as:
```
mkdir /etc/golang
export GOPATH=/etc/golang
```

Then install with GO:

  > go get github.com/cblomart/vsphere-graphite

The executable should be `$GOPATH/bin/vsphere-graphite` and is now a binary for your architecture/OS

#### Run on Commandline

  > vsphere-graphite

#### Install as a service

  > vsphere-graphite install

#### Run as a service

  > vsphere-graphite start
  >
  > vsphere-graphite status
  >
  > vsphere-graphite stop

#### Remove service

  > vsphere-graphite remove
  
## Contributors

No open source projects would live and thrive without common effort. Here is the section were the ones that help are thanked:

- [sofixa](https://github.com/sofixa)
- [BlueAceTS](https://github.com/BlueAceTS)
- [NoMotion](https://github.com/NoMotion)

Also keep in mind that if you can't contribute code, issues and imporvement request are also a key part of a project evolution!
So don't hesitate and tell us what doens't work or what you miss.

## License

The MIT License (MIT)

Copyright (c) 2016 cblomart

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
