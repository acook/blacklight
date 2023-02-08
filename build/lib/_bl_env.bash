#!/usr/bin/env bash

source "$(dirname "$BASH_SOURCE")/_shared.bash"

export BL_ROOT_PATH
export BL_EXT_PATH
export BL_LOCAL_PATH

BL_ROOT_PATH="$(readlink -e "$SCRIPT_SHARED_DIR/../..")" || die "unable to set env var"
BL_EXT_PATH="$(readlink -e "$BL_ROOT_PATH/ext")" || die "unable to set env var"
BL_LOCAL_PATH="$BL_EXT_PATH/local"

export BL_MAIN_PATH
export BL_OUT_NAME
export BL_BIN_NAME
export BL_BIN_DIR
export BL_BIN_PATH
export BL_TEST_DIR

BL_MAIN_PATH="$BL_ROOT_PATH/src/redlight.c"
BL_OUT_NAME="a.out"
BL_BIN_NAME="redlight"
BL_BIN_DIR="$(readlink -f "$SCRIPT_SHARED_DIR/../out")"
BL_BIN_PATH="$BL_BIN_DIR/$BL_BIN_NAME"
BL_TEST_DIR="$BL_ROOT_PATH/test"

PATH="$PATH:$BL_LOCAL_PATH/bin:/usr/bin:/usr/local/bin:/bin"
PATH="$PATH:$HOME/bin:$HOME/xbin"

# tooling

export BL_ALLOC
export BL_CC
export BL_CCOPTS

export BL_LINKER="ld.lld"
export BL_STRIP="strip"

# detect jemalloc and setup
if command_exists jemalloc-config; then
  BL_ALLOC="-L $(jemalloc-config --libdir) -Wl,-rpath,$(jemalloc-config --libdir) -ljemalloc $(jemalloc-config --libs)"
else
  warn "(nonfatal error) jemalloc-config not found"
  BL_ALLOC=" "
fi
BL_CCOPTS="-std=gnu17 $BL_ALLOC"

if command_exists ecc; then
  BL_CC="ecc"
  BL_CCOPTS="-static $BL_CCOPTS"
  PATH="$PATH:$BL_EXT_PATH/ellcc/bin"
else
  BL_CC="clang"
  #PATH="$PATH"
fi

# vars used by other commands
export CC="$BL_CC"
export CLINKER="$BL_LINKER"
export LD_LIBRARY_PATH="$BL_LOCAL_PATH/lib"

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
