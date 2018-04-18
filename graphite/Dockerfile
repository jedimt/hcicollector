# Graphite stack

# Build from Ubuntu base
FROM ubuntu:14.04.2
#FROM ubuntu:16.04

# This suppresses a bunch of annoying warnings from debconf
ENV DEBIAN_FRONTEND noninteractive

# Install system dependencies
RUN \
  apt-get -qq install -y software-properties-common && \
  add-apt-repository -y ppa:chris-lea/node.js && \
  apt-get -qq update -y && \
  apt-get -qq install -y build-essential curl \
    # Graphite dependencies
    python-dev libcairo2-dev libffi-dev python-pip \
    # Supervisor
    supervisor \
    # nginx + uWSGI
    nginx uwsgi-plugin-python \
    # StatsD
    nodejs

# Install StatsD
RUN \
  mkdir -p /opt && \
  cd /opt && \
  #curl -sLo statsd.tar.gz https://github.com/etsy/statsd/archive/v0.7.2.tar.gz && \
  curl -sLo statsd.tar.gz https://github.com/etsy/statsd/archive/v0.8.0.tar.gz && \
  tar -xzf statsd.tar.gz && \
  #mv statsd-0.7.2 statsd
  mv statsd-0.8.0 statsd

# Install Python packages for Graphite
RUN pip install graphite-api[sentry] whisper carbon

# Optional install graphite-api caching
# http://graphite-api.readthedocs.org/en/latest/installation.html#extra-dependencies
# RUN pip install -y graphite-api[cache]

# Graphite
COPY carbon.conf /opt/graphite/conf/carbon.conf
COPY storage-schemas.conf /opt/graphite/conf/storage-schemas.conf
COPY storage-aggregation.conf /opt/graphite/conf/storage-aggregation.conf
# Supervisord
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
# StatsD
COPY statsd_config.js /etc/statsd/config.js
# Graphite API
COPY graphite-api.yaml /etc/graphite-api.yaml
# uwsgi
COPY uwsgi.conf /etc/uwsgi.conf
# nginx
COPY nginx.conf /etc/nginx/nginx.conf
#COPY basic_auth /etc/nginx/basic_auth

# nginx
# Ports
##Graphite API 8080; Carbon line receiver 2003; Cabon pickle receiver 2004; 
##Carbon cache query 7002; StatsD UDP 8125; StatsD Admin 8126 
EXPOSE 80 \
# graphite-api
8080 \
# Carbon line receiver
2003 \
# Carbon pickle receiver
2004 \
# Carbon cache query
7002 \
# StatsD UDP
8125 \
# StatsD Admin
8126

# Launch stack
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
