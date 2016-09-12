#!/usr/bin/env bash

# Treat unset variables as an error
set -o nounset

# setup environment (if available)
source gg 2> /dev/null

# utility functions

function scriptname() {
  basename ${BASH_SOURCE[1]-$thisscript}
}

function realscriptpath() {
  # usage: realscriptpath [relative_path]
  # relative path is optional and expected to be reliable
  SOURCE="${BASH_SOURCE[0]}"
  while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
    DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
  done
  echo "$( cd -P "$(dirname "$SOURCE")/$1" && pwd )"
}

function warn() { echo "$@" 1>&2; }

function timer() { date "+%s.%N"; }

function elapsed() {
  started_at=$1
  ended_at=$2

  if [[ -x $(which bc) ]]; then
    dt=$(echo "$ended_at - $started_at" | bc)
    dd=$(echo "$dt/86400" | bc)
    dt2=$(echo "$dt-86400*$dd" | bc)
    dh=$(echo "$dt2/3600" | bc)
    dt3=$(echo "$dt2-3600*$dh" | bc)
    dm=$(echo "$dt3/60" | bc)
    ds=$(echo "$dt3-60*$dm" | bc)

    printf " -- time elapsed: %d:%02d:%02d:%02.4f\n" $dd $dh $dm $ds
  else
    warn "START: $started_at"
    warn "END: $ended_at"
  fi
}


# Make sure we're in the right directory
export OLDDIR="$(pwd)"
export BLROOT="$(realscriptpath ..)"
cd $BLROOT
