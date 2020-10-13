#!/usr/bin/env bash
set -euxo pipefail
echo "building docker image and pushing to project $1"
docker build -t calculator/calculator .
docker tag calculator/calculator gcr.io/"$1"/calculator/calculator:latest
docker push gcr.io/"$1"/calculator/calculator:latest
echo "done"
