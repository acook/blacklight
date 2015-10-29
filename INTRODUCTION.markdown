INTRODUCTION TO blacklight
==========================



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

