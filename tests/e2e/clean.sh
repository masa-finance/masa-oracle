#!/bin/bash
set -ex

pushd tests/e2e/work/
    docker compose down
popd

rm -rf tests/e2e/work/*

docker volume rm work_masa
docker volume rm work_masa2
docker rmi masa-node