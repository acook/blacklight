#pragma once

#include "./bl_types.h"

typedef byte opcode;
static const opcode Ref = (1 << 7);  // xor with reference types
static const opcode Invalid = 0;
static const opcode Text = 34;
static const opcode Number = 35;
// static const opcode Double = 36;
// static const opcode Decimal = ??;
// static const opcode Atom = 39;
static const opcode Extended = 255;
// static const opcode Rune = ??;
