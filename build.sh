#!/bin/bash

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$PROJECT_DIR"

UID_GID="$(stat -c "%u:%g" .)"
mkdir bin >/dev/null 2>&1
chown "$UID_GID" bin

docker run --rm -it --user="$UID_GID" -v "$PWD":/go/src/github.com/rootkiwi/screen_share_remote_go \
    rootkiwi/screen_share_remote_go:build go run build.go "$@"

EXPECTED_IMAGE_ID="sha256:565113a9c8ed59e099805e7ac0dc7299e9343e819d934c833435f301e0d268a7"
ACTUAL_IMAGE_ID=$(docker inspect --format='{{.Id}}' rootkiwi/screen_share_remote_go:build)

if [ "$EXPECTED_IMAGE_ID" != "$ACTUAL_IMAGE_ID" ]
then
    echo
    echo "docker build image outdated, run following to update:"
    echo "docker pull rootkiwi/screen_share_remote_go:build"
fi
