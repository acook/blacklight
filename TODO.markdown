blacklight TODO
===============

> Man, I really got my work cut out for me. -- me

I imagine 3 major phases of blacklight:

- **bianca**  - the feature-complete and working version of blacklight according to my original notes
- **azzurra** - a major overhaul which fixes design flaws I discover during **bianca**
- **nera**    - bytecode interpreter (or LLVM) with a focus on speed and portability

Currently we're in the midst of a headlong rush to **bianca** to prove out the base concepts.

Commandline
-----------

- --version : displays version and build date
- --help    : displays commandline usage and url
- --check   : parses and lexes, but doesn't eval


Vectors
-------

- eval op for WVs


Numbers
-------

- floats
- cell-like operations on Ns
- byte operations on Ns
- conversions between N and C


Logic & Loops
-------------

- is, datatype literals
- while


IO & System
-----------

- file descriptors: especially for stdin, stdout, and stderr
- files
- sockets
- system commands
- test Windows/Linux support


Text
----

- chars should work with vector's `app` op
- backslash chars
- double-backslash chars for control characters and unicode
- t-to-v to display Tag information
- datatypes.Print() vs datatypes.String()
- warn


Objects
-------

- object datatype
- set/get slots by name
- parent slot (with `child` op) will check parents for a given slot name
- object-stack


codename: bluelight
-------------------

A wishlist of features for a major overhaul of blacklight.

- **rework syntax**
- - no stack, queue, or object literal syntax (just need to update spec)
- - blocks are delimited with paired square brackets
- - call `new` on datatype literals instead of dedicated operations
- **cleaner stack implementation**
- - no meta-stack
- - one system-stack per thread/user/instance
- - can make abitrary stacks the @stack "cd mystack"
- **tighter object integration**
- - one system-object per thread/user/instance
- - setwords and getwords apply directly to the current object
- - can make arbitrary objects the @object "cd myobject"
- - delegate unknown slot names to a list of other objects
- - delegate a particular slot name to a particular object
- **improved queues**
- - infinite queues
- - separate blocking-channel datatype
- **text**
- - double-quoted strings
- - internal string escapes
- **logic and loops**
- - repeat n times
- **fibers and events**
- - a thread dedicated to a reactor loop which schedules fibers and checks for events
- **metaprogramming**
- - WVs are called blocks, may be a different underlying type
- - Lisp-like macros using blocks
- - Lisp-like list manipulations of blocks
- **better IO**
- - trunc, readwrite, cursors
- **generic sequences**
- - group together sequnce-types
- **object persistence mechanism**
- - group datatype
- - groups contain auto-incremented slots that store objects
- - store/retr objects by id
- - search stored objects lazily

