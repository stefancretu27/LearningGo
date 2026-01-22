package go_didactical_apps

/**
* 1. A programming paradigm is a fundamental style or approach to programming. It defines how you structure and organize your code, and how you think about solving problems.
* Some common paradigms include:
*	- Imperative: Code describes how to do things step by step (e.g., C, Go).
*	- Declarative: Code describes what to do, not how (e.g., SQL, HTML).
*	- Object-Oriented: Organizes code around objects and classes (e.g., Java, C++).
*	- Functional: Emphasizes pure functions and immutability (e.g., Haskell, parts of JavaScript).
*	- Procedural: A subtype of imperative, focused on procedures/functions (e.g., C, Pascal).
*	- Concurrent/Parallel: Focused on executing multiple tasks simultaneously (e.g., Go, Erlang).
*
* Go's Paradigms
*	- Imperative
*	- Procedural
*	- Concurrent (via goroutines and channels)
*	- It’s not object-oriented in the traditional sense (no classes or inheritance), but it supports composition and interfaces, which enable similar patterns.
*
* 2. Idioms are like micro-patterns that are tailored to the language's syntax, semantics, and philosophy.
*
* 3. About Go
*	- Go is designed to have a simpler syntax and to compile faster than C/C++, with focus on concurrency
*	- Compilation: likewise C, it is fully compiled. The compilation output is a single static binary, with all dependecies linked within
*	- OOP: not supported
*	- Overloading: not supported
*	- Inheritance: not supported
*	- Classes: only structs, which can nest other structs => composition/aggregation
*	- Data access specifier: unlike c++, which has public/private/protected, Go uses starting capital letter for publicly accessible data, 
*	and small starting letter for private, both at struct level, and at package level
*	- Manual memory management: no, as it uses a garbage collector, included in the runtime library, that is compiled together with the executable
*	- Templates: uses generics, a similar concept
*/