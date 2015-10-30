blacklight
==========

> blacklight: the language for people who want to code in it

blacklight is a general-pupose multithreading concatenative stack-based programming language
with first-class queues and prototype-style objects.

blacklight doesn't need "variables",
instead program state is stored in a simple set of stacks.

blacklight doesn't need "lambdas" or function "arguments",
instead you have a vectors of operations on a stack.

Features
--------

blacklight is ideal for MapReducing, ETLing, and processing.

*Below list refers to features planned for the first release.*

- data structures are threadsafe by default
- stack-oriented programming
- concatenative programming
- builtin parallelism through native multithreading (in-progress)
- objects with prototypical delgation (planned)


Documentation
-------------

blacklight uses a set of builtin operations for creating and maipulating data.

These operations are documented [in this document on Google Sheets](https://docs.google.com/spreadsheets/d/1Kz5zFMtGjBrdEHrHySFmB5UttQ6lXcKM6C-2iz5VDiM/edit?usp=sharing).

Also check out the `INTRODUCTION.markdown` to find out more about the underlying concepts.


Futher Reading
---------------

There is a rich history of languages like blacklight and as a result there is a lot of pre-existing information to introduce you to these concepts.

### Resources

- [Why Concatenative Programming Matters (An Introduction)](http://evincarofautumn.blogspot.com/2012/02/why-concatenative-programming-matters.html)
- [The Joy of Concatenative Languages](http://www.codecommit.com/blog/cat/the-joy-of-concatenative-languages-part-1)
- [Concatenative Languages on C2](http://c2.com/cgi/wiki?ConcatenativeLanguage)
- [Concatenative Programming on Wikipedia](https://en.wikipedia.org/wiki/Concatenative_programming_language)
- [Stack-Oriented Programming on Wikipedia](https://en.wikipedia.org/wiki/Stack-oriented_programming_language)
- [Cool Stuff Built in Forth](http://www.forth.org/successes.html) (hint: NASA used it)
- [Concatenative Programming Wiki](http://concatenative.org)

### Other concatenative Languages

- [Forth](https://en.wikipedia.org/wiki/Forth_(programming_language)) - this is the big one, the granddaddy of them all
- [JS-FORTH](https://repl.it/languages/forth) - fun to play around with in your browser
- [Factor](http://factorcode.org/) - a cool modern concatenative language
- [Postscript](https://en.wikipedia.org/wiki/PostScript) - most printers use Postscript to do their jobs


blacklight is influenced by Forth, REBOL, Factor, Self, and Redis.
