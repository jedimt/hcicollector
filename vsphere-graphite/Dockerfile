FROM ubuntu

RUN apt-get update && apt-get install -y curl

RUN curl -L https://github.com/cblomart/vsphere-graphite/releases/download/v0.4/vsphere-graphite_linux_amd64_v0.4.tgz --output /vsphere-graphite_linux_amd64_v0.4.tgz && \
    tar -zxvf /vsphere-graphite_linux_amd64_v0.4.tgz    && \
    mv /vsphere-graphite /usr/local/bin/vsphere-graphite && \
    mv /vsphere-graphite.json /etc/vsphere-graphite.json && \
    rm /vsphere-graphite_linux_amd64_v0.4.tgz

ADD vsphere-graphite.json /etc/vsphere-graphite.json
ADD start-vsphere-graphite.sh /start-vsphere-graphite.sh 

#ENTRYPOINT [ "vsphere-graphite" ]
ENTRYPOINT [ "/start-vsphere-graphite.sh" ]

