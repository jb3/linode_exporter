#!/usr/bin/env sh

set -eux

DATE=$(date +"%Y.%m.%d")
COMMITS=$(git rev-list --count HEAD --since $(date -v-1d +'%Y-%m-%d'))

echo $DATE.$COMMITS
