#!/usr/bin/env bash

source "$(dirname "$BASH_SOURCE")/_shared.bash"

export BL_ROOT_PATH="$(readlink -e "$SCRIPT_SHARED_DIR/../..")" || die "unable to set env var"
export BL_EXT_PATH="$(readlink -e "$BL_ROOT_PATH/ext")" || die "unable to set env var"
export BL_LOCAL_PATH="$BL_EXT_PATH/local"

export BL_MAIN_PATH="$BL_ROOT_PATH/src/redlight.c"
export BL_OUT_NAME="a.out"
export BL_BIN_NAME="redlight"
export BL_BIN_DIR="$(readlink -f "$SCRIPT_SHARED_DIR/../out")"
export BL_BIN_PATH="$BL_BIN_DIR/$BL_BIN_NAME"
export BL_TEST_DIR="$BL_ROOT_PATH/test"

# tooling
if command_exists ecc; then
  export BL_CC="ecc"
  export BL_CCOPTS="-static $BL_ALLOC"
  export PATH="$BL_EXT_PATH/ellcc/bin"
else
  export BL_CC="clang"
  export BL_CCOPTS="$BL_ALLOC"
  export PATH=""
fi

export PATH="$PATH:$BL_LOCAL_PATH/bin:/usr/bin:/usr/local/bin:/bin"
export BL_STRIP="strip"

# detect jemalloc and setup
if command_exists jemalloc-config; then
  export BL_ALLOC="-L $(jemalloc-config --libdir) -Wl,-rpath,$(jemalloc-config --libdir) -ljemalloc $(jemalloc-config --libs)"
else
  warn "(nonfatal error) jemalloc-config not found"
  export BL_ALLOC=""
fi

# vars used by other commands
export CC="$BL_CC"
export CLINKER="$BL_CC"

# random functions
bigsay() {
  if command_exists figlet; then
    if [[ -f $BL_EXT_PATH/figlet/chunky.flf ]]; then
      font="-f chunky -d $BL_EXT_PATH/figlet"
    else
      font=""
    fi
    figlet -t $font "$*"
  else
    echo "$*"
  fi
}
banner() {
  if command_exists figlet; then
    if [[ -f $BL_EXT_PATH/figlet/chunky.flf ]]; then
      font="-f chunky -d $BL_EXT_PATH/figlet"
    else
      font=""
    fi
    figlet -t -c $font "$*"
  else
    echo "$*"
  fi
}
bl_banner() {
  colorreset ; colorbg black ; colorfg violet
  ansieol ; echo ; ansieol ; echo ; ansieol ; echo ; ansieol ; echo ; ansieol ; echo
  ansiup 5
  if command_exists figlet; then
    banner "blacklight"
  else
    echo "
 __     __              __     __ __         __     __
|  |--.|  |.---.-.----.|  |--.|  |__|.-----.|  |--.|  |_
|  _  ||  ||  _  |  __||    < |  |  ||  _  ||     ||   _|
|_____||__||___._|____||__|__||__|__||___  ||__|__||____|
                                     |_____|"
  fi
  colorreset ; echo
}

bl_banner

# Make sure we're in the right directory
safe_cd "$BL_ROOT_PATH" > /dev/null

if scriptsame; then
  printenv | grep "^BL_*"
fi
