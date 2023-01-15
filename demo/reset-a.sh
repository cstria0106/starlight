#!/bin/bash
## This script reset the experiment environment
## Please cd into the project's root directory and run this script using this command:
##  sudo ./demo/reset.sh

STARLIGHT_SNAPSHOTTER_ROOT=/var/lib/starlight-grpc/

# Stop Starlight and containerd
systemctl stop containerd
systemctl stop starlight
pkill -9 'containerd' | true
pkill -9 'starlight-grpc' | true


# Clear containerd folder
rm -rf /var/lib/containerd

# Clear starlight folder
if [ -d "${STARLIGHT_SNAPSHOTTER_ROOT}sfs/" ] ; then
    find "${STARLIGHT_SNAPSHOTTER_ROOT}sfs/" \
         -maxdepth 1 -mindepth 1 -type d -exec sudo umount -f "{}/m" \;
fi
rm -rf "${STARLIGHT_SNAPSHOTTER_ROOT}"*

# Remove Redis data folder
rm -rf /tmp/test-redis-data
rm -rf /tmp/test-pg-data


# Restart the service
./out/starlight-grpc  run --server=10.0.2.2:8090 --insecure --log-level=debug &
containerd &
