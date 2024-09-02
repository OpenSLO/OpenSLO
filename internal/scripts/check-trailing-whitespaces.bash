#!/usr/bin/env bash

fileExtensionsToIgnore=('.ico' '.png' '.desc')

found=0
for file in $(git ls-tree -r HEAD --name-only); do
  for ext in "${fileExtensionsToIgnore[@]}"; do
    if [[ "$file" == *"$ext" ]]; then
      continue 2
    fi
  done
  if [ ! -f "$file" ]; then
    continue
  fi
  if grep -qE ' +$' "$file"; then
    if ((found != 1)); then
      echo "Trailing whitespaces found!" >&2
    fi
    grep -nE ' +$' "$file" | awk -F: '{print $1}' | while IFS= read -r line; do
      echo "$file:$line" >&2
    done
    found=1
  fi
done
if ((found == 1)); then
  exit 1
fi
