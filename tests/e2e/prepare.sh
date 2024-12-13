#!/bin/bash
set -ex

TEE_IMAGE="${TEE_IMAGE:-masaengineering/tee-worker:main}"

if [ ! -d tests/e2e/work ]; then
  mkdir -p tests/e2e/work
else
  rm -rf tests/e2e/work/*
fi

mkdir -p tests/e2e/work

cp -rfv tests/e2e/fixtures/docker-compose.template.yaml tests/e2e/work/docker-compose.yaml
cp -rfv tests/e2e/fixtures/env tests/e2e/work/.env

docker build -t masa-node -f Dockerfile .

docker volume create work_masa
docker volume create work_masa2

pubkey1=$(docker run --rm -v work_masa:/home/masa/ -e PRINT_PUBKEY=true masa-node)
pubkey2=$(docker run --rm -v work_masa2:/home/masa/ -e PRINT_PUBKEY=true masa-node)

sed -i 's/%%IMAGE%%/masa-node/g' tests/e2e/work/docker-compose.yaml
sed -i 's/%%TEEIMAGE%%/'$TEE_IMAGE'/g' tests/e2e/work/docker-compose.yaml
sed -i 's/%%NODE1PUB%%/'$pubkey2'/g' tests/e2e/work/docker-compose.yaml
sed -i 's/%%NODE2PUB%%/'$pubkey1'/g' tests/e2e/work/docker-compose.yaml