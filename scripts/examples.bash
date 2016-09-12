#!/usr/bin/env bash

thisscript="examples"
blacklight="$1"
function usage() { warn "usage: $(scriptname) ./path/to/blacklight example1.bl [example2.bl ...]"; exit -1; }

source "$(dirname ${BASH_SOURCE[0]})/_shared.bash"

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
