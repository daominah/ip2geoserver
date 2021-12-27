#!/usr/bin/env bash

set -e

export dockerImgTag=daominah/echo_ip_httpsvr
export dockerCtnName=echo_ip_httpsvr

docker build --tag=${dockerImgTag} .
docker rm -f echo_ip_httpsvr
docker run -dit --restart always --name ${dockerCtnName} -p 20891:20891 ${dockerImgTag}
