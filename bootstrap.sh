#!/bin/bash

# Set up basic_auth
file=basic_auth
if [ ! -f $file ]
then
    echo "Enter a password for the graphite user"
    read passwd
    basic_auth=$(openssl passwd -crypt $passwd)
    echo "graphite:$basic_auth" > ./graphiteconfig/basic_auth
else
    echo "$file file already exists.  Did you already run this?"
fi
