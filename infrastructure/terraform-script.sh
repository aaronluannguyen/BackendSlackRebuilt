#!/usr/bin/env bash
export DO_PAT=9c1177bca3ac3f4ab580712c8036a96cf4ce0863a58784c792b0a4d3405c064e
export SSH_FINGERPRINT=58:3f:9d:d0:5d:34:f2:0a:0d:e1:10:bc:f5:f8:e4:24
terraform plan \
-var "do_token=${DO_PAT}" \
-var "pub_key=$HOME/.ssh/id_rsa.pub" \
-var "pvt_key=$HOME/.ssh/id_rsa" \
-var "ssh_fingerprint=$SSH_FINGERPRINT"
terraform apply \
-var "do_token=${DO_PAT}" \
-var "pub_key=$HOME/.ssh/id_rsa.pub" \
-var "pvt_key=$HOME/.ssh/id_rsa" \
-var "ssh_fingerprint=$SSH_FINGERPRINT"