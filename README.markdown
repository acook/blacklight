Redlight
========

`redlight` is a next-gen VM for `blacklight` implemented in C11.

Supported Platforms
-------------------

If you run into any issues compiling `redlight` for the following targets, please file a bug:

- Win32/MinGW/Musl x86_64
- Linux/Musl x86_64
- Linux/Musl AArch64le (ARM 64bit with LittleEndian support)

Support for the following platforms is desired but untested at this time:

- Linux/Musl ppc64le (POWER8+)
- AIX/Musl ppc64le (POWER8+)

`redlight` is developed on BunsenLabs (a Debian-based distro) and Windows 10 on x86_64 using ELLCC and Musl.

Building
--------

The `build/all` script expects all dependencies (including `ecc`) to be located under `./ext/`.

## Debian / Ubuntu
If you are on a distro which uses `apt` and Debian-like package names you can run `build/deps` to download and compile all prerequisites.

Then run the `build/all` script.

## Windows
You can download WinGW ELLCC from its older versions subdirectory.
The `build/all` script has not been tested on Windows.

To build, make sure you have `ecc` in your system PATH and run:
`ecc -static src/redlight.c`

## Other
The `build/all` script should work on any system with `bash`, `ecc`, `jemalloc-config`, and the relevant dependencies in `./ext`.

In addition, it automatically adds `./ext/local/bin` to the PATH so `make install` or `ln -s` tools there.

## GCC / Clang
`redlight` has been tested with `gcc` and LLVM's `clang`.
If you have compilation errors ensure you are compiling in gnu11/C11 mode.

Static Analysis
---------------

There is a `build/check` script which runs `cppcheck`, `flawfinder`, and `splint`.

`infer` support is planned but requires a from-source build with OCAML on Debian stable. Pull requests welcome.

