#!/bin/bash

# create docker volumes
if ! docker volume create dockerPostgresVolume; then
    echo "Warning - failed to create docker volume: dockerPostgresVolume"
fi
if ! docker volume create dockerPGAdminVolume; then
    echo "Warning - failed to create docker volume: dockerPGAdminVolume"
fi

# create certs
if ! mkdir certs; then
    echo "Warning - failed to create directory for certificates, already exist?"
fi
if ! openssl req -new -newkey rsa:4096 -x509 -sha256 -days 365 -nodes -out certs/server.rsa.crt -keyout certs/server.rsa.key -subj "/C=US/ST=Texas/O=Aurora/"; then
    echo "ERROR - failed to create certificates, exiting"
    exit -1
fi

# build initial binary
if ! go build aurora/main.go; then
    echo "ERROR - failed to create go binary, exiting"
    exit -1
fi
if ! sudo setcap 'cap_net_bind_service=+ep' main; then
    echo "ERROR - failed to give binary bind access, exiting"
    exit -1
fi

# ready
echo "Setup Complete"
echo "To start server, make sure docker container inside of database is running."
echo "Then, run ./main"
exit 0