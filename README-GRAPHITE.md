###Graphite Docker

This will set up all of the components that Graphite needs to run.  It also
runs Grafana and connects to the graphite container.

A few assumptions for configurations are made, if you have any specific needs
you can modify the default configs and rebuild the containers with the updated
configs.

**Instructions for running**

 * Download and install docker-compose (`sudo pip install -U docker-compose`)
 * Clone this repo (`git clone https://github.com/jmreicha/graphite-docker`)
 * Modify install script (`cd graphite-docker && chmod +x bootstrap.sh`)
 * Run script (`./bootstrap.sh`)
 * Run the graphite stack (`docker-compose up`)
 * Or in detached mode (`docker-compose up -d`)
 * Open up a browser and navigate to address where these containers are running

**Details**

In the `docker-compose.yml` file we are setting the default grafana password to
`password`.  You can either modify the compose file update the account password
in the GUI after logging in the first time.

This set up uses basic_auth to secure graphite, you can view more info here -
http://nginx.org/en/docs/http/ngx_http_auth_basic_module.html.  To turn basic
auth off you can modify the Dockerfile and nginx config to remove references
the basic_auth settings.

The basic_auth generation relies on openssl for creating the user and
hash for the password so one of these tools must be installed for the
basic_auth component to work.

**Notes**

In a production environment it would be a good idea to mount in a separate
volume for data in to the /data as a mount point so that if the OS runs out of
space the volume can be attached somewhere else.

You may need to run the script with sudo if the /data volume has restricted
permissions.
