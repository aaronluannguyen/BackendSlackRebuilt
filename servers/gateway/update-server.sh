#!/usr/bin/env bash
docker rm -f summaryserver
docker pull aaronluannguyen/summary-server
docker run -d --name summaryserver \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=/etc/letsencrypt/live/api.aaronnluannguyen.me/fullchain.pem \
-e TLSKEY=/etc/letsencrypt/live/api.aaronnluannguyen.me/privkey.pem \
aaronluannguyen/summary-server