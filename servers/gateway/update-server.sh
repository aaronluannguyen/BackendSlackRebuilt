#!/usr/bin/env bash
docker rm -f gatewayinfo344
docker pull aaronluannguyen/gatewayinfo344
docker run -d --name gatewayinfo344 \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=/etc/letsencrypt/live/api.aaronnluannguyen.me/fullchain.pem \
-e TLSKEY=/etc/letsencrypt/live/api.aaronnluannguyen.me/privkey.pem \
aaronluannguyen/gatewayinfo344