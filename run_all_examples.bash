#!/usr/bin/env bash

source gg 2> /dev/null

timestamp="$(date --utc "+%Y.%m.%d")"
shortsha="$(git rev-parse --short HEAD)"
blacklight="./bin/blacklight_$shortsha-$timestamp"

echo " -- building blacklight binary..."
go build -o "$blacklight" src/github.com/acook/blacklight/*.go

if [[ -x $blacklight ]]; then
  echo " -- binary built at: \"$blacklight\""
else
  echo " -- something went wrong!"
  echo " -- binary not found at: \"$blacklight\""
  exit -1
fi

failures=0
examples=$(find ./examples -maxdepth 1 -regextype posix-extended -regex '\./.*/[a-z].*.bl' -type f -printf '%p\n')

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

exit $failures
