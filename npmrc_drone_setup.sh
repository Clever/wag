#!/usr/bin/env bash
set -eu

sed -i.bak s/\${npm_token}/$npm_token/ .npmrc_drone
mv .npmrc_drone test/js/.npmrc
