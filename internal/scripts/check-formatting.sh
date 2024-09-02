#!/usr/bin/env bash

set -e

TMP_DIR=$(mktemp -d)

cleanup_git() {
  git -C "$TMP_DIR" clean -df
  git -C "$TMP_DIR" checkout -- .
}

main() {
  cp -r . "$TMP_DIR"
  cleanup_git

  make -C "$TMP_DIR" format

  CHANGED=$(git -C "$TMP_DIR" diff --name-only)
  if [ -n "${CHANGED}" ]; then
    printf >&2 "The following file(s) are not formatted:\n%s\n" "$CHANGED"
    exit 1
  else
    echo "Looks good!"
  fi
}

main "$@"
