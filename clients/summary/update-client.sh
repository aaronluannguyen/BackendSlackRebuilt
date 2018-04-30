#!/usr/bin/env bash
docker rm -f summaryclient
docker pull aaronluannguyen/summary-client
docker run -d --name summaryclient -p 80:80 -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro aaronluannguyen/summary-client