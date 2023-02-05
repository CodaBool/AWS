#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

path_to_dockerfile="$1"

file_hashes="$(
  cd $path_to_dockerfile \
  && find . -type f -not -path './.**' \
    -not -path '*/venv/*' \
    -not -path '*/node_modules/*' \
  | sort \
  | xargs md5sum
)"

hash="$(echo "$file_hashes" | md5sum | cut -d' ' -f1)"
echo '{ "hash": "'"$hash"'" }'