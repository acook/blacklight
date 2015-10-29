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


Meet the Stacks
---------------

Stacks are kinda, like, a thing.

### The Three Stacks

When you start blacklight it automatically initializes a few things for you, the most apparent of which will be your home-stack. This is the very first system-stack that your operations will be dealing with but it’s hardly the last. You can easily create more of these system-stacks, but you usually only work with one at a time, the currently active stack is called the `@stack`.

When you create a new system-stack the old one doesn’t go away forever, but it’s put onto a special stack called the `meta-stack` or `$stack`. It is, like the name implies, a stack of stacks. The old `system-stack` is now called the prev-stack or the `^stack`. Sometimes it’s handy to reference one of the three stacks directly and you can use use their prefix sigil to do so: `$` dollar-sign for the meta-stack, `@` at-sign for the current-stack, and `^` caret for the next-stack.

These three stacks are the ones you’ll use the most and are the only three stacks with special syntax. The `@stack` and `^stack` are actually relative stacks, reference, and any time you create a new system-stack, those references change.

If you have a system-stack on the `@stack` you can use the same operations on them as a normal user-stack!

#### system-stack

These are your bread and butter, they store all your data, and are the core of blacklight. They are stacks with a little something extra and there are a ton of built-in operations that manipulate these stacks.

These are what your `@stack` and `^stack` are made of.

Make a new one with the `$new` operation and then use the `swap`, `drop`, `rot`, `over` and other operations to manipulate it.

*Any of blacklight’s datatypes can be stored in a system stack.*

#### user-stack

The user-stack is a normal datatype in blacklight for doing any sort of stack-based algorithms and storage. They’re not as fancy as system-stacks, but they’re thread-safe and easy to use.

Make a new one with the `s-new` operation, then you can `push` items onto them, or `pop` items off of them.

*Any of blacklight’s datatypes can be stored on a user-stack.*

#### meta-stack

There is only one meta-stack (per instance, per user, per thread) which governs all of your system-stacks. There’s a handful of very simple but very powerful operations that you can use on the meta-stack, but be careful, with great power comes completely screwing up your blacklight!

*Only system-stacks can be stored on the meta-stack.*


Documentation
-------------

blacklight uses a set of builtin operations for creating and maipulating data.

These operations are documented [in this document on Google Sheets](https://docs.google.com/spreadsheets/d/1Kz5zFMtGjBrdEHrHySFmB5UttQ6lXcKM6C-2iz5VDiM/edit?usp=sharing).


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

### Other concatenative Languages

- [Forth](https://en.wikipedia.org/wiki/Forth_(programming_language)) - this is the big one, the granddaddy of them all
- [JS-FORTH](https://repl.it/languages/forth) - fun to play around with in your browser
- [Factor](https://en.wikipedia.org/wiki/Factor_(programming_language)) - a cool modern concatenative language
- [Postscript](https://en.wikipedia.org/wiki/PostScript) - most printers use Postscript to do their jobs


blacklight is influenced by Forth, REBOL, Factor, Self, and Redis.
