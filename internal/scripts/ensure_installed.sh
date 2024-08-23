#!/usr/bin/env bash

set -e

LOCAL_BIN_DIR="${LOCAL_BIN_DIR:-./bin}"

_binary() {
  if [ ! -f "${LOCAL_BIN_DIR}/${1}" ]; then
    echo "$1 was not found in $LOCAL_BIN_DIR" >&2
    make "install/$1"
  fi
}

# It's cheaper to run yarn install then do any other checks.
_yarn() {
  make "install/yarn"
}

case "$1" in
  binary) _binary "$2" ;;
  yarn) _yarn "$2" ;;
  *) echo "invalid source provided: $1" >&2 && exit 1 ;;
esac
