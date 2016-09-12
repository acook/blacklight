#!/usr/bin/env bash

thisscript="ci"
source "$(dirname ${BASH_SOURCE[0]})/_shared.bash"

if [[ ! -x $(which go) ]]; then
  warn " -- can't find Go binary in PATH!"
  warn " -- PATH: $PATH"
fi

blacklight="$(./scripts/build_test.bash)"
./scripts/run_all_examples.bash $blacklight
