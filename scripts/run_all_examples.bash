#!/usr/bin/env bash

set -o nounset                              # Treat unset variables as an error
cd "$(dirname $0)/.."
export dir="$(pwd)"

source gg 2> /dev/null

function timer() {
  date +%s.%N
}

function elapsed() {
  started_at=$1
  ended_at=$2
  dt=$(echo "$ended_at - $started_at" | bc)
  dd=$(echo "$dt/86400" | bc)
  dt2=$(echo "$dt-86400*$dd" | bc)
  dh=$(echo "$dt2/3600" | bc)
  dt3=$(echo "$dt2-3600*$dh" | bc)
  dm=$(echo "$dt3/60" | bc)
  ds=$(echo "$dt3-60*$dm" | bc)

  printf " -- time elapsed: %d:%02d:%02d:%02.4f\n" $dd $dh $dm $ds
}

timestamp="$(date --utc "+%Y.%m.%d")"
shortsha="$(git rev-parse --short HEAD)"
blacklight="./bin/blacklight_$shortsha-$timestamp"

echo " -- building blacklight binary..."
go build -o "$blacklight" src/*.go

if [[ -x $blacklight ]]; then
  echo " -- binary built at: \"$blacklight\""
else
  echo " -- something went wrong!"
  echo " -- binary not found at: \"$blacklight\""
  exit -1
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
