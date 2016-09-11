#!/usr/bin/env bash

source "$(dirname $0)/_shared.bash"

timestamp="$(date --utc "+%Y.%m.%d")"
shortsha="$(git rev-parse --short HEAD)"
blacklight="$BLROOT/bin/blacklight_$shortsha-$timestamp"

warn " -- building blacklight binary..."
go build -o "$blacklight" "$BLROOT/src/"*.go

if [[ -x $blacklight ]]; then
  warn " -- binary built: $(basename $blacklight)"
  echo "$blacklight"
else
  warn " -- something went wrong!"
  warn " -- binary not found at: \"$blacklight\""
  exit -1
fi
