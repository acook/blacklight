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

# SCRIPT FUNCTIONS

# usage: displayname <path>
# takes a path and generates "last_folder/filename" string
displayname() {
  basename -z $(dirname $(readlink -m "$1"))
  echo -ne "/"
  basename "$1"
}
# formats the current script name
scriptname() { displayname "$SCRIPT_CURRENT_PATH"; }
# determine the currently executing script via caller
# can't be nested in other functions
scriptcaller() { readlink -e "$(caller | cut -d " " -f2-)"; }
# for conditionals, determines if caller is the same as the main parent script
scriptsame() { [[ $SCRIPT_MAIN_PATH == "$(readlink -e $(caller | cut -d " " -f2-))" ]]; }
# used internally to set the current script global
_set_scriptcurrent() {
  local fallback=${BASH_SOURCE[2]}
  local script=${1:-$fallback}

  SCRIPT_CURRENT_PATH=$(readlink -m "$script");
}
# source a script only once
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
# source a script once or more
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

# DISPLAY FUNCTIONS
say()  { echo -ne " -- ($(scriptname) @ $(ts)) : $*\n"; }
warn() { say "$*" >&2; }
# usage: sayenv <VARNAME>
sayenv() { say "$1=$(eval "echo -ne \$$1")"; }

# EXIT STATUS FUNCTIONS
ok()   { say "\e[32m(ok) $*\e[0m"; exit 0; }
die()  { warn "\e[31m(die) $*\e[0m"; exit 1; }
# usage: die_status <status> [message]
die_status() { warn "\e[31m(died with status code $1) ${*:2}\e[0m"; exit $1; }

# WRAPPER FUNCTIONS

# if cd fails then we should exit
safe_cd() { cd $1 || die "couldn't cd! $1"; }
# used for conditionals to determine presence of a command or executable
command_exists() { command -v $1 > /dev/null 2>&1; }
# usage: run "title" <command> [args]
# display command to run, confirm it exists, run it, output a warning on failure
run() {
  say "running $1 command: \`${@:2}\`"
  if command_exists $2; then
    "${@:2}"
    ret=$?
    [[ $ret ]] || warn "$1 command exited with status code $?"
    return $ret
  else
    warn "command \`$2\` not found"
    return 255
  fi
}
# usage: run_or_die "title" <command> [args]
# as run, but die if command missing or exits with an error
run_or_die() {
  say "running $1 command: \`${@:2}\`"
  command_exists $2 || die "command \`$2\` not found"
  $2 "${@:3}" || die_status $? "$2 command"
}

# UTILITY FUNCTIONS

# usage: realpath <path>
# attempts to resolve all symlinks until the origin path is discovered
# circular symlinks will keep it looping forever
realpath() {
  p="$1"
  # loop until the file is no longer a symlink (or doesn't exist)
  while [[ -h $p ]]; do
    d="$( cd -P "$( dirname "$p" )" && pwd )"
    p="$(readlink -e "$p")"
    # if $p was a relative symlink
    # we need to resolve it relative to the path where the symlink file was located
    [[ $p != /* ]] && p="$d/$p"
  done
  echo "$( cd -P "$(dirname "$p")" && pwd )"
}

# usage: elapsed <start_time> <end_time>
# takes two unix timestamps with nanoseconds
# returns the difference in human-readable format
elapsed() {
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
