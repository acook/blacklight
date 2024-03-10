#!/usr/bin/env bash

source "$(dirname "${BASH_SOURCE[0]}")/_shared.bash"

timestamp="$(date --utc "+%Y.%m.%d")"
shortsha="$(git rev-parse --short HEAD)"
blacklight="$BLROOT/bin/blacklight_$timestamp-$shortsha"

warn " -- building blacklight binary..."
go build -o "$blacklight" "$BLROOT/src/"*.go
exitstatus=$?

if [[ $exitstatus == 0 && -x $blacklight ]]; then
  warn " -- binary built: $(basename "$blacklight")"
  
  file "$blacklight" >&2
  ls -sh "$blacklight" >&2
  
  # this filename is captured by ci.bash
  echo "$blacklight"
else
  warn " -- something went wrong!"
  warn " -- go exit status: $exitstatus"
  warn " -- binary not found at: \"$blacklight\""
  exit 255
fi

