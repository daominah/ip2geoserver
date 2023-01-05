#!/usr/bin/env bash

set -e

export DOCKER_IMG_TAG=daominah/ip2geoserver
export DOCKER_CTN_NAME=ip2geoserver

docker build --tag=${DOCKER_IMG_TAG} .
docker rm -f ${DOCKER_CTN_NAME} 2>/dev/null;
docker run -dit --restart always --name ${DOCKER_CTN_NAME} -p 20891:20891 ${DOCKER_IMG_TAG}
