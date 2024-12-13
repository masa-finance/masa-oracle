#!/bin/bash

pushd tests/e2e/work/
    docker compose up -d
popd

E2E=true go test -v ./tests/e2e