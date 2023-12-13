#!/bin/bash

PREV_ROOT = `docker info -f '{{ .DockerRootDir }}'`
systemctl stop docker
systemctl stop docker.socket
systemctl stop containerd

cat '{
  "data-root": "$1"
}' >>  /etc/docker/daemon.json

systemctl start docker
