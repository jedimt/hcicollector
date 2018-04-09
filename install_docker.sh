#Install DockerCE for Ubuntu 16.04.3

# Get Docker installed on the host
apt update

apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"

apt update

apt -y install docker-ce=17.12.1~ce-0~ubuntu

#Install docker-compose
apt -y install docker-compose
