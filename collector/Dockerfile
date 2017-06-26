# Docker container for solidfire-graphite-collector
# aaron.patten@netapp.com | @jedimt
FROM ubuntu:16.04
MAINTAINER aaron.patten@netapp.com
RUN apt update \
        && apt install -y python-pip git \
        && apt clean
RUN pip install solidfire-sdk-python==1.2.0.108 \
        requests \
        graphyte \
        logging
ADD solidfire_graphite_collector_v2.py /solidfire_graphite_collector_v2.py
ADD wrapper.sh /wrapper.sh
ENTRYPOINT /wrapper.sh

