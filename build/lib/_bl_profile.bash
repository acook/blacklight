#!/usr/bin/env bash

source "$(dirname "$BASH_SOURCE")/_shared.bash" 2> /dev/null
include "bl_env"

if [ -f /etc/bash_completion ] && ! shopt -oq posix; then
  . /etc/bash_completion 2> /dev/null
fi

export PAGER=less
export CLICOLOR=1
export HISTCONTROL=erasedups
export HISTSIZE=10000
shopt -s histappend

# load up aliases and functions
source ~/.bash_aliases

__bl_prompt() {
    local EXITSTATUS="$?"
    PS1=""

    local lfcr='\n\[\r\]'
    local fg='\['"$(colorfg 244)"'\]'
    local bright='\['"$(colorfg 231)"'\]'
    local warn='\['"$(colorfg red3)"'\]'
    local okay='\['"$(colorfg green3)"'\]'
    local blfg='\['"$(colorfg violet)"'\]'

    if [ $EXITSTATUS == 0 ]; then
        PS1+="${okay}"
    else
        PS1+="${warn}"
    fi
    PS1+="[ $(ts) ] "
    PS1+="[ proccess exited with status code $EXITSTATUS ]"
    PS1+="${lfcr}${lfcr}"

    PS1+="${blfg}blacklight build console "
    PS1+="${fg}\w"

    local branch="$(git branch 2> /dev/null | \
      awk -v on_color="$bright" -v branch_color="$fg"\
        '$1 =="*" { printf on_color " @ " branch_color $2 " " $3 }')"
    PS1+="$branch"

    local dirty="$(git status 2> /dev/null |
      awk 'index($0,"Untracked files:") { unknown = 1 }
           index($0,"modified:") { changed = 1 }
           index($0,"new file:") { new = 1 }
           END {
             if (unknown) printf "?"
             else if (changed) printf "!"
             else if (new) printf "."
           }')"

    PS1+="${bright}${dirty}${lfcr}"

    PS1+="${bright}\$"
    PS1+="\[$(colorreset)\] "
}
export PROMPT_COMMAND=__bl_prompt
