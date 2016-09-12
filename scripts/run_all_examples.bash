#!/usr/bin/env bash

thisscript="run_all_examples"
blacklight="$1"

source "$(dirname ${BASH_SOURCE[0]})/_shared.bash"
function usage() { warn "usage: $(scriptname) ./path/to/blacklight"; exit -1; }

if [[ ! -x $blacklight ]]; then
  warn " -- binary not found at: \"$blacklight\""
  usage
fi

examples=$(find ./examples -maxdepth 1 -regextype posix-extended -regex '\./.*/[a-z].*.bl' -type f -printf '%p\n')
time ./scripts/examples.bash "$blacklight" "$examples"
