#!/usr/bin/env bash

thisscript="examples"
blacklight="$1"
shift
examples="$*"
function usage() { warn "usage: $(scriptname) ./path/to/blacklight example1.bl [example2.bl ...]"; exit -1; }

source "$(dirname ${BASH_SOURCE[0]})/_shared.bash"

if [[ ! -x $blacklight ]]; then
  warn " -- binary not found at: \"$blacklight\""
  usage
fi

failures=0
failed=""
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

ended_at=$(timer)

echo
echo " -- failures: $failures"
echo " -- failed: ${failed-<all passed>}"

elapsed $started_at $ended_at

exit $failures
