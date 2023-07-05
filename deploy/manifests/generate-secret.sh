#!/bin/bash

kubectl create secret generic go-webhook --from-env-file=.env --dry-run=client -oyaml > deploy/manifests/go-webhook/secrets.yaml
