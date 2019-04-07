#!/bin/bash

if ! go build aurora/main.go; then
	echo "failed to build"
	exit -1
fi
if ! sudo setcap 'cap_net_bind_service=+ep' main; then
	echo "failed to give main permissions"
	exit -1
fi
