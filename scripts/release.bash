#!/usr/bin/env bash

source "$(dirname ${BASH_SOURCE[0]})/_shared.bash"

mkdir release

timestamp="$(date --utc "+%Y.%m.%d")"
shortsha="$(git rev-parse --short HEAD)"
blacklight="$BLROOT/release/blacklight_$timestamp-$shortsha"


warn " -- building blacklight binary for macOS..."
bl_macos="$blacklight-macos"
env GOOS=darwin GOARCH=amd64 go build -v -a -o "$bl_macos" "$BLROOT/src/"*.go

warn " -- building blacklight binary for Windows..."
bl_win="$blacklight-windows.exe"
env GOOS=windows GOARCH=amd64 go build -v -a -o "$bl_win" "$BLROOT/src/"*.go

warn " -- building blacklight binary for Linux AMD64..."
bl_linux="$blacklight-linux"
env GOOS=linux GOARCH=amd64 go build -v -a -o "$bl_linux" "$BLROOT/src/"*.go

# upgrade this to arm64 when I get a pi3
warn " -- building blacklight binary for Linux ARM7x32..."
bl_arm="$blacklight-linuxARM"
env GOOS=linux GOARCH=arm GOARM=7 go build -v -a -o "$bl_arm" "$BLROOT/src/"*.go

warn " -- binary prefix: $blacklight"

warn " -- macOS binary: $(basename $bl_macos)"
warn " -- Windows binary: $(basename $bl_win)"
warn " -- Linux binary: $(basename $bl_linux)"
warn " -- LinuxARM binary: $(basename $bl_arm)"

warn " -- done!"
