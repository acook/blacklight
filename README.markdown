blacklight
==========

> `blacklight` is a general-pupose **multithreading** concatenative stack-based programming language
with first-class **queues** and **objects** with delegation.

blacklight doesn't need "variables",
instead program state is stored in a simple set of stacks.

blacklight doesn't need "lambdas" or function "arguments",
instead you have a vectors of operations on a stack.


Features
--------

blacklight is awesome, here's a few reasons why:

- easy to use builtin parallelism through native multithreading
- data structures are threadsafe by default
- native queue type for inter-thread communication
- stack-oriented programming
- concatenative programming
- objects with delgation *(planned)*
- great for MapReducing, ETLing, and general data processing


Documentation
-------------

- blacklight uses a set of builtin operations for creating and maipulating data.
These operations are documented [in this document on Google Sheets](https://docs.google.com/spreadsheets/d/1Kz5zFMtGjBrdEHrHySFmB5UttQ6lXcKM6C-2iz5VDiM/edit?usp=sharing).

- The [examples directory](https://github.com/acook/blacklight/tree/master/examples) contains several basic scripts to get you started.

- Also check out the blacklight [Wiki](https://github.com/acook/blacklight/wiki/Meet-the-Stacks) to find out more about the underlying concepts.
