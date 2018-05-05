#!/usr/bin/env bash

source "$(dirname ${BASH_SOURCE[0]})/_shared.bash"

export BL_ROOT_PATH="$(readlink -e "$SCRIPT_SHARED_DIR/../..")" || die "unable to set env var"
export BL_EXT_PATH="$(readlink -e "$BL_ROOT_PATH/ext")" || die "unable to set env var"
export BL_LOCAL_PATH="$BL_EXT_PATH/local"

export BL_MAIN_PATH="$BL_ROOT_PATH/src/redlight.c"
export BL_OUT_NAME="a.out"
export BL_BIN_NAME="redlight"
export BL_BIN_DIR="$(readlink -f "$SCRIPT_SHARED_DIR/../out")"
export BL_BIN_PATH="$BL_BIN_DIR/$BL_BIN_NAME"

# set PATH
export PATH="$BL_EXT_PATH/ellcc/bin:$BL_LOCAL_PATH/bin:/usr/bin:/usr/local/bin:/bin"

# detect jemalloc and setup
if command_exists jemalloc-config; then
  export BL_ALLOC="-L $(jemalloc-config --libdir) -Wl,-rpath,$(jemalloc-config --libdir) -ljemalloc $(jemalloc-config --libs)"
else
  warn "(nonfatal error) jemalloc-config not found"
  export BL_ALLOC=""
fi

# tooling
export BL_CC="ecc" # FIXME: allow other compilers
export BL_CCOPTS="-o $BL_OUT_NAME -static $BL_ALLOC"
export BL_STRIP="strip"

# vars used by other commands
export CC="$BL_CC"
export CLINKER="$BL_CC"


# Make sure we're in the right directory
safe_cd $BL_ROOT_PATH

if scriptsame; then
  printenv | grep "^BL_*"
fi
