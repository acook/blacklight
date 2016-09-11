#!/usr/bin/env bash

# Treat unset variables as an error
set -o nounset

# setup environment (if available)
source gg 2> /dev/null

# utility functionso

scriptpath() {
  # usage: scriptpath [relative_path]
  # relative path is optional and expected to be reliable
  SOURCE="${BASH_SOURCE[0]}"
  while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
    DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
  done
  echo "$( cd -P "$(dirname "$SOURCE")/$1" && pwd )"
}

warn() { echo "$@" 1>&2; }

# Make sure we're in the right directory
export OLDDIR="$(pwd)"
export BLROOT="$(scriptpath ..)"
cd $BLROOT
