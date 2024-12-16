#!/bin/bash

pushd tests/e2e/work/
    docker compose down
popd

rm -rf tests/e2e/work/*

docker volume rm work_masa || true
docker volume rm work_masa2 || true
docker rmi masa-node || true