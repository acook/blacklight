Redlight
========

`redlight` is a next-gen VM for `blacklight` implemented in C17.

Supported Platforms
-------------------

If you run into any issues compiling `redlight` for the following targets, please file a bug:

- Win32/MinGW x86_64
- Linux x86_64
- Linux AArch64le (ARM 64bit with LittleEndian support)

`redlight` is developed on Solus Linux on x86_64 using Clang.

Building
--------

## Linux
If you download `instant install` and have it in one of the standard binary locations you can run `./build/deps` to download and compile all prerequisites automatically.

Then run the `./build/all` script.

## Windows
The `build/*` scripts have not been tested on Windows.

I have only done this in the past using an old build of ECC, which was unfortunately abandoned in 2017.

To build, make sure you have `ecc` in your system PATH and run:
`ecc -static src/redlight.c`

## Other
The `./build/all` script should work on any system with `bash`, `clang`, and the relevant dependencies in `./ext`.

In order to allow for isolated builds, the `./build/all` script automatically rebuilds its own PATH environment which looks in `./ect/local` for its dependencies.
If you are building your own dependencies set their install prefix to `./ext/local` or add symlinks as needed.

## GCC
`redlight` has also been tested with `gcc` and LLVM's `clang` with `glibc`.
If you have compilation errors ensure you are compiling in gnu17/C17 mode.

Static Analysis
---------------

There is a `build/check` script which runs `cppcheck`, `flawfinder`, and `splint`. (splint support paused)

`infer` support is planned but requires a from-source build with OCAML on Debian stable. Pull requests welcome.

