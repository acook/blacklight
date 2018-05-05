#!/usr/bin/env bash

set -o nounset

if [[ -z ${_BASH_SHARED_LIB+unset} ]]; then
  declare -a _BASH_SHARED_LIB
  _BASH_SHARED_LIB=("$(readlink -e "$BASH_SOURCE")")
else
  return 0
fi

echo " -- ($(basename $(dirname $(readlink -m "${BASH_SOURCE[-1]}")))/$(basename "${BASH_SOURCE[-1]}") @ $(date "+%Y-%m-%d %T")) : setting up..." >&2

SCRIPT_SHARED_PATH="$(readlink -e "$BASH_SOURCE")"
SCRIPT_SHARED_NAME="$(basename "$SCRIPT_SHARED_PATH")"
export SCRIPT_SHARED_DIR="$(dirname "$SCRIPT_SHARED_PATH")"
export SCRIPT_ORIG_PWD="$(pwd -P)"

export SCRIPT_MAIN_PATH="$(readlink -e "$0")"
export SCRIPT_MAIN_NAME="$(basename "$SCRIPT_MAIN_PATH")"
export SCRIPT_MAIN_DIR="$(dirname "$SCRIPT_MAIN_PATH")"
export SCRIPT_MAIN_EXE="$(basename "$SCRIPT_MAIN_DIR")/$SCRIPT_MAIN_NAME"

export SCRIPT_CURRENT_PATH=$SCRIPT_SHARED_PATH

# time-related functions
ts()      { date "+%Y-%m-%d %T"; }              # local timestamp for output
ts_file() { date --utc "+%Y-%m-%d-%H-%M-%S"; }  # timestamp for filenames, in UTC for consistency
ts_unix() { date "+%s.%N"; }                    # unix timestamp with nanoseconds for use with elapsed()

# datatype functions
array_contains () {
  local e match="$1"
  shift
  for e; do [[ "$e" == "$match" ]] && return 0; done
  return 1
}

# script-related functions
displayname() {
  basename -z $(dirname $(readlink -m "$1"))
  echo -ne "/"
  basename "$1"
}
scriptname() { displayname "$SCRIPT_CURRENT_PATH"; }
scriptcaller() { readlink -e "$(caller | cut -d " " -f2-)"; } # this can't be nested
scriptsame() { [[ $SCRIPT_MAIN_PATH == "$(readlink -e $(caller | cut -d " " -f2-))" ]]; }
_set_scriptcurrent() { 
  local fallback=${BASH_SOURCE[2]}
  local script=${1:-$fallback}

  SCRIPT_CURRENT_PATH=$(readlink -m "$script"); 
}
include() { 
  local fullpath="$SCRIPT_SHARED_DIR/_$1.bash"
  if [[ ! -f $fullpath ]]; then
    die "unable to include \`$fullpath\`: file not found"
  fi
  if [[ ! " ${_BASH_SHARED_LIB[@]} " =~ " ${1} " ]]; then
    _BASH_SHARED_LIB+=("$1")
    _set_scriptcurrent "$fullpath"
    source "$fullpath" || die "error including $fullpath"
    _set_scriptcurrent
  fi
}
load() { 
  if [[ ! -f $1 ]]; then
    die "unable to load \`$1\`: file not found"
  fi
  _set_scriptcurrent "$1"
  source "$1" || die "error loading \`$1\`"; 
  _set_scriptcurrent
}

trace() { # for debugging bash functions
  local frame=0
  echo -ne "TRACE ($frame): "
  while caller $frame; do
    ((frame++));
    echo -ne "TRACE ($frame): "
  done
  echo BASH
}

# output-related functions
say()  { echo -ne " -- ($(scriptname) @ $(ts)) : $*\n"; }
warn() { say "$*" >&2; }
sayenv() { say "$1=$(eval "echo -ne \$$1")"; }

# exit-related functions
ok()   { say "(ok) $*"; exit 0; }
die()  { warn "(die) $*"; exit 1; }
die_status() { warn "(died with status code $1) ${*:2}"; exit $1; }

# wrapper functions
safe_cd() { cd $1 || die "couldn't cd! $1"; }
command_exists() { command -v $1 > /dev/null 2>&1; }
run() { 
  say "running $1 command: \`${@:2}\`"
  if command_exists $2; then
    eval "${@:2}" || warn "$1 command exited with status code $?"
  else
    warn "command \`$2\` not found"
  fi
}
run_or_die() {
  say "running $1 command: \`${@:2}\`"
  command_exists $2 || die "command \`$2\` not found"
  eval "${@:2}" || die_status $? "$2 command"
}

function realpath() {
  p="$1"
  # loop until the file is no longer a symlink (or doesn't exist)
  # circular symlinks will keep it looping forever
  while [[ -h $p ]]; do 
    d="$( cd -P "$( dirname "$p" )" && pwd )"
    p="$(readlink -e "$p")"
    # if $p was a relative symlink
    # we need to resolve it relative to the path where the symlink file was located
    [[ $p != /* ]] && p="$d/$p" 
  done
  echo "$( cd -P "$(dirname "$p")" && pwd )"
}

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

_set_scriptcurrent
