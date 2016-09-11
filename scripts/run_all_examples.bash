#!/usr/bin/env bash

blacklight="$1"
source "$(dirname $0)/_shared.bash"

this="run_all_examples"
function usage() { warn "usage: $this ./path/to/blacklight"; exit -1; }

if [[ ! -x $blacklight ]]; then
  warn " -- binary not found at: \"$blacklight\""
  usage
fi

failures=0
examples=$(find ./examples -maxdepth 1 -regextype posix-extended -regex '\./.*/[a-z].*.bl' -type f -printf '%p\n')
started_at=$(timer)

for file in $examples; do
  echo
  echo " -- running: $file"
  echo test | $blacklight $file

  if [ $? -ne 0 ]; then
    failures=$(($failures + 1))
    failed="$failed$(basename $file) "
  fi
done

echo
echo " -- failures: $failures"
echo " -- failed: ${failed-<all passed>}"

ended_at=$(timer)

elapsed $started_at $ended_at

exit $failures
