#!/usr/bin/env bash

set -eux

DATE=$(date +"%Y.%m.%d")

PLATFORM=$(uname -o)

if [[ $PLATFORM == "Darwin" ]]; then
    SINCE=$(date -v-1d +'%Y-%m-%d')
else
    SINCE=$(date -d '-1 day' +'%Y-%m-%d')
fi

COMMITS=$(git rev-list --count HEAD --since $SINCE)

echo $DATE.$COMMITS
