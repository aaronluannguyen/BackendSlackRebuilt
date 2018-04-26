#!/usr/bin/env bash
docker rm -f clientinfo344
docker pull aaronluannguyen/clientinfo344
docker run -d --name clientinfo344 -p 80:80 -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro aaronluannguyen/clientinfo344