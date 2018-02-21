#!/bin/bash

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$PROJECT_DIR"

UID_GID="$(stat -c "%u:%g" .)"
mkdir bin >/dev/null 2>&1
chown "$UID_GID" bin

docker run --rm -it --user="$UID_GID" -v "$PWD":/go/src/github.com/rootkiwi/screen_share_remote_go \
    rootkiwi/screen_share_remote_go:build go run build.go "$@"

EXPECTED_IMAGE_ID="sha256:d80d1144d9d1d5f379166763ecb72803f1ef10a9761b8bb38696d0542b005b96"
ACTUAL_IMAGE_ID=$(docker inspect --format='{{.Id}}' rootkiwi/screen_share_remote_go:build)

if [ "$EXPECTED_IMAGE_ID" != "$ACTUAL_IMAGE_ID" ]
then
    echo
    echo "docker build image outdated, run following to update:"
    echo "docker pull rootkiwi/screen_share_remote_go:build"
fi
