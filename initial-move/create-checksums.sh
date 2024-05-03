#!/usr/bin/env bash
# Copyright 2024 The Mothership Authors
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail
shopt -s globstar

# Generate checksums for all files in the current directory and subdirectories.
# The output is a file called checksums.txt.
# The format is:
#   <filename> <checksum>
FILENAME="/tmp/checksums.txt"
rm -f "${FILENAME}"
touch "${FILENAME}"

file_gen() {
  local filename="${1}"
  local checksum="$(sha256sum "${filename}" | awk '{print $1}')"
  echo "${filename} ${checksum}" >> "${FILENAME}"
}

for x in **/*.rpm; do
  # Generate checksums in parallel, then output to a file.
  file_gen "${x}" &
done

wait
