#!/usr/bin/env bash

set -e

GEN_PATHS="**/*.go **/*.yaml"
TMP_DIR=$(mktemp -d)

cleanup_git() {
  git -C "$TMP_DIR" clean -df
  git -C "$TMP_DIR" checkout -- .
}

main() {
  cp -r . "$TMP_DIR"
  cleanup_git

  make -C "$TMP_DIR" generate

  CHANGED=$(git -C "$TMP_DIR" status --porcelain ${GEN_PATHS})
  if [ -n "${CHANGED}" ]; then
    printf >&2 "There are generated changes that are not committed:\n%s\n" "$CHANGED"
    exit 1
  else
    echo "Looks good!"
  fi
}

main "$@"
