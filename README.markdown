![blacklight logo](http://i.imgur.com/N81hd1M.png)

[![Release Tag](https://img.shields.io/github/tag/acook/blacklight.svg?maxAge=2592000)](https://github.com/acook/blacklight/releases)
[![Build Status](https://travis-ci.org/acook/blacklight.svg)](https://travis-ci.org/acook/blacklight)
[![Kanban on Trello](https://img.shields.io/badge/kanban-trello-026AA7.svg)](https://trello.com/b/vygKBL4j)
[![Join the chat at https://gitter.im/acook/blacklight](https://img.shields.io/gitter/room/acook/blacklight.js.svg)](https://gitter.im/acook/blacklight)

> `blacklight` is a *programming language* which is concurrent, stack-based, and concatenative  (BLPL)

> `blacklight` is a **virtual machine** for implementing highly concurrent languages (BLVM)

> `blacklight` is a *data interchange* format for communicating between processes and across networks (BLBC)

Features
--------

blacklight (BLVM) is awesome, here's a few reasons why:

- easy to use builtin parallelism through native concurrency primatives
- threadsafe communication between concurrency units
- rich datatype primitives
- an easy to use homoiconic Forth-like assembly language (BLPL)
- runtime bytecode manipulation and generation
- UTF-8 native datatypes
- multi-architecture and cross-platform (currently: x86_64, ARM, macos, linux, windows)
- (planned) highly optimized vector operations on supported CPUs
- (planned) simple FFI to Rust and C
- (planned) security contexts and permissions

Documentation
-------------

- The [blacklight Wiki](https://github.com/acook/blacklight/wiki) has documentation and links (work in progress).
- The [examples directory](https://github.com/acook/blacklight/tree/master/examples) contains several demonstration scripts to get you started.
