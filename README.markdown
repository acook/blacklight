[![blacklight logo](http://i.imgur.com/N81hd1M.png)](https://github.com/acook/blacklight#readme)

[![Release Tag](https://img.shields.io/github/tag/acook/blacklight.svg?maxAge=2592000)](https://github.com/acook/blacklight/releases)
[![Build Status](https://acook.semaphoreci.com/badges/blacklight/branches/bianca.go.svg?style=shields&key=ea65e38e-7450-4e89-ae6d-53d6eb4bb433)](https://acook.semaphoreci.com/projects/blacklight)
[![Build status](https://ci.appveyor.com/api/projects/status/7h1e1sly5024l6im/branch/master?svg=true)](https://ci.appveyor.com/project/acook/blacklight/branch/master)
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
- (in progress) highly optimized vector operations on supported CPUs
- (in progress) simple FFI to Rust and C
- (planned) security contexts and permissions

Documentation
-------------

- The [blacklight Wiki](https://github.com/acook/blacklight/wiki) has documentation and links (work in progress).
- The [examples directory](https://github.com/acook/blacklight/tree/master/examples) contains several demonstration scripts to get you started.

BLPOC
-----

The current implementation of `blacklight` is a proof-of-concept. It's functional but intended primarily for proving out features, strategies, and specifications. Once The ABI is stable it will be reimplemented with optimization and compatibility in mind against a full test suite. As is, there is very little about `blacklight` that isn't subject to change to better reflect the results of research and experimentation. 
