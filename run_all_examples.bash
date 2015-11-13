#!/usr/bin/env bash

source gg 2> /dev/null
files=$(find ./examples -maxdepth 1 -regextype posix-extended -regex '\./.*/[a-z].*.bl' -type f -printf '%p\n')

failures=0

for file in $files; do
  echo
  echo " -- running: $file"
  go run src/github.com/acook/blacklight/*.go $file

  if [ $? -ne 0 ]; then
    failures=$(($failures + 1))
    failed="$failed$(basename $file) "
  fi
done

echo
echo " -- failures: $failures"
echo " -- failed: $failed"

exit $failures
