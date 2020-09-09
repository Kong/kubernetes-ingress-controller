#!/bin/bash

# Example script to generate certs

openssl req -new -x509 -nodes -newkey ec:<(openssl ecparam -name secp384r1) \
  -keyout cluster.key -out cluster.crt \
  -days 1095 -subj "/CN=kong_clustering"