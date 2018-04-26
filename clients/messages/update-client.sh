#!/usr/bin/env bash
docker rm -f client344
docker pull aaronluannguyen/client344
docker run -d --name client344 -p 80:80 -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro aaronluannguyen/client344