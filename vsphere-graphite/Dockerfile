FROM ubuntu

RUN apt-get update && apt-get install -y curl

RUN curl -L https://github.com/cblomart/vsphere-graphite/releases/download/v0.3c/vsphere-graphite_linux_amd64_v0.3c.tgz --output /vsphere-graphite_linux_amd64_v0.3c.tgz && \
    tar -zxvf /vsphere-graphite_linux_amd64_v0.3c.tgz    && \
    mv /vsphere-graphite /usr/local/bin/vsphere-graphite && \
    mv /vsphere-graphite.json /etc/vsphere-graphite.json && \
    rm /vsphere-graphite_linux_amd64_v0.3c.tgz

ADD vsphere-graphite.json /etc/vsphere-graphite.json

ENTRYPOINT [ "vsphere-graphite" ]
