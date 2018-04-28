#!/usr/bin/env bash

export THIS=()

type $THIS

this_push() { THIS+=("$1"); }
