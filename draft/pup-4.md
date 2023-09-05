# Pattern for Universal Protocols (PUP)
### draft/pup-4.md

# Discussion

XXX describe origin and timeline of concepts -- non-profit work,
Challenger, IBM lab shots, trading floors, isconf, bootstrapping and
turing papers, governance in the 2000s, SARS, climate, etc.  mine
ghent talk for ideas.

XXX forward references before the Concepts heading, backward
references after 

XXX move PUP to the end of the document, after the concepts section,
into and Implementation section.  Implementation section should
be a procedure for implementing a grid.

XXX need a generic name like "grid" for the system, that isn't yet in
common use in computing or organizational theory.  retitle this doc
and repo accordingly.  

XXX this doc is likely a book, not a paper.



PUP is a notional reference implementation of a more general concept.
This book is a guide for building a UIT-based system, and PUP is a
specific implementation of that system.

The more general concept is a distributed, decentralized, peer-to-peer
system for handling information and interaction.  The system is
designed to be scalable, fault-tolerant, and self-healing.  It is
designed to be a general-purpose computing platform, able to decide
anything that is decideable by machine (Turing-equivalent), but also
able to decide things that are not decideable by machine
(a superset of Turing-complete, a hypercomputer).  



# Concepts

## Mapping

Mapping is a concept from set theory.  A mapping is a relationship
between two sets.  

A mapping can be one-to-one, where each element of the first set is
mapped to exactly one element of the second set.  For example, a set
of numbers can be mapped to a set of other numbers in a one-to-one
relationship -- "1" is mapped to "2", "2" is mapped to "3", and so on.  
XXX picture

A mapping can be one-to-many, where each element of the first set is
mapped to multiple elements in the second set.  For example, a set of
numbers can be mapped to a set of other numbers in a one-to-many
relationship -- "1" is mapped to "2", "3", and "4", "2" is mapped to
"3", "4", and "5", and so on. XXX picture

A mapping can be many-to-one, where multiple elements of the first set
are mapped to a single element in the second set.  For example, a set
of numbers can be mapped to a set of other numbers in a many-to-one
relationship -- "1" and "2" are mapped to "3", "2" and "3" are mapped
to "4", and so on. XXX picture

A mapping can be many-to-many, where multiple elements of the first
set are mapped to multiple elements in the second set.  For example, a
set of numbers can be mapped to a set of other numbers in a
many-to-many relationship -- "1" and "2" are mapped to "3", "4", and
"5", "2" and "3" are mapped to "4", "5", and "6", and so on. XXX
picture

## Function

A function is a one-to-one or many-to-one [mapping](#mapping) between
two sets.  A function maps each element of the domain to exactly one
element of the range.  

In this book, we'll use the terms "input" and "output" to refer to the
domain and range of a function.  

For example, the function f(x, y) = x + y maps the input domain of
integers x and y to the output range of integers that are the sum of x
and y. In this example, the answer to the question "What is f(1, 2)?"
has the single answer "3".  XXX picture

An example of a non-function is the one-to-many [mapping](#mapping) of
"1" to "2", "3", and "4".  This is not a function because each element
of the input domain is mapped to multiple elements of the output
range.  The answer to the question "What is f(1)?" is "2", "3", or
"4" -- there is no single answer. XXX need a picture of this.

## Bidirectional function

A bidirectional [function](#function) is a one-to-one
[mapping](#mapping) where each element of the domain is mapped to
exactly one element of the range and vice versa.  In other words, each
input produces exactly one output, and each output is produced by
exactly one input -- it's possible to easily "work backwards" from the
output to the input.

For example, the function f(x) = x + 1 is a bidirectional function,
because each range integer is the exactly one more than each domain
integer, and each domain integer is exactly one less than each range
integer.  The answer to the question "If x is 3, what is f(x)?" is
"4", and the answer to the question "If f(x) is 4, what is x?" is
"3". XXX picture

## One-way function

A one-way [function](#function) is a many-to-one [mapping](#mapping)
where each element of the range is mapped to exactly one element of
the domain, but there many by many elements of the domain that map to
each element of the range.  In other words, multiple inputs might
produce the same output, so it's not possible to easily "work
backwards" from an output to discover the input.

For example, the function f(x, y) = x + y is a one-way function,
because each output range integer is the sum of an infinite number of
pairs of input domain integers -- if the sum is 4, then the pairs (1,
3), (2, 2), (3, 1), (0, 4), (-1, 5), and so on, all map to 4.  There
is no way to know which pair of input integers was used to generate
the output of 4.  XXX picture

One-way functions are called "trapdoor functions" in cryptography:
They make it easy to encrypt input data to get output ciphertext, but
difficult to decrypt the ciphertext without access to one or more of
the inputs.  For example, if f(x, key) = y is a trapdoor function,
then it is easy to encrypt data by providing the input data and the
key, but it is difficult to decrypt y without the key.  
XXX picture

## Decision problem

A decision problem is a problem that can be answered with a true/false
answer.  For example, "2 + 2 = 4" is a decision problem that has the
answer "true".  "2 + 2 = 5" is a decision problem that has the answer
"false".  "What is 2 + 2?" is not a decision problem, because it
cannot be answered with a true/false answer.
XXX picture

A decision problem can be thought of as a [one-way
function](#one-way-function) that maps a domain of inputs to a range
of outputs that has exactly two elements -- true and false.
XXX picture

## Function problem

A function problem is a problem that can be answered with a
[function](#function). For example, "What is 2 + 2?" is a function
problem.  Function problems are a superset of [decision
problems](#decision-problem).
XXX picture

## Algorithm

An algorithm is a finite sequence of steps that can be followed to
solve a problem.  For example, the algorithm "add x + y, then send
the output to Alice" has two steps; addition and sending.  
XXX picture

An algorithm can be represented as a [function](#function) that maps
the input domain of the problem to the output range of the solution.
XXX picture

An algorithm might be a [bidirectional
function](#bidirectional-function), where it's possible to easily
"work backwards" from the output to the input.  For example, the
algorithm "add x + y, then send 'add', x, and the output to Alice" is
a bidirectional function, because Alice can see that it was an
addition operation, so she can subtract x from the output to get y. It
is generally the case that another algorithm is needed to find the
input that produces a particular output.  In this case, Alice uses the
reverse algorithm "if the first word is 'add', then subtract the
second word from the third word to get the output".  XXX picture

An algorithm might be a [one-way function](#one-way-function), where
it's not possible to easily "work backwards" from the output to the
input.  For example, the algorithm "add x + y, then send x and the
output to Alice" is a one-way function, because Alice doesn't know
that there was an addition, so she can't subtract x from the output to
get y.  There is no algorithm that can be followed to find the input
that produces a particular output.
XXX picture

## Halting problem

> *There is no algorithm that can be followed to detect whether an arbitrary algorithm will halt or run forever.*
>
> *There is no procedure that can be followed to detect whether an arbitrary procedure will produce desired results.*

The halting problem was a key insight in Alan Turing's 1936 paper that
defined the theoretical basis for all modern computing systems
[turing][turing].  In the paper, he addressed a [decision
problem](#decision-problem) that asks, "given a [machine](#machine)
and an input to that machine, will the machine halt (stop executing)
or will it run forever?"  

Turing proved that there is no machine that can solve the halting
problem for all machines and all inputs.  A simple proof of this is
that if such an algorithm existed, it could be used to create this
paradox: "If I am a program that will halt when given my own code as
input, then loop forever, otherwise halt."

In modern computing terms, there is no program that can be written to
detect whether an arbitrary program will function correctly: "If I am a
program that will function correctly when given my own code as input,
then crash, otherwise function correctly."

This has implications for software testing: it is impossible to write
a program that can detect whether another program will function
correctly.  This is why software testing is so difficult: it is
impossible to know whether a program will function correctly in all
cases.  The best that can be done is to test a program with a large
number of inputs, and hope that it will function correctly for those
inputs not tested.  

The halting problem can be translated into other walks of life by
generalizing the paradox:  "If I am a procedure that will produce a
desired output when given my own procedure as input, then produce an
undesired output, otherwise produce a desired output."  


## Predicting 

The [halting problem](#halting-problem) describes one of the
factors that contribute to the difficulty of predicting the future
state of a complex system.  

Another contributing factor is the [uncertainty [uncertainty
principle][uncertainty] in physics, which states that it is impossible
to know both the position and momentum of a particle at the same time:
it is impossible to predict the future state of a particle with
absolute certainty, therefore it is impossible to completely predict
the future state of any larger system.

A third contributing factor is the [butterfly effect][butterfly], which
states that a small change in the initial state of a system can lead
to large changes in the future state of the system.  


## Computable function

A computable function is a [function](#function) that can be computed
by an algorithm. All algorithms are functions, but not all functions
are computable.

For example, the function "find the square root of any number" is not
a computable function, because there is no algorithm that can be
followed to find the exact square root of any arbitrary number.  This
is due to the infinite number of digits in the square root of an
irrational number -- the algorithm would never terminate.

XXX venn diagram
A computable function can be represented as a one-to-one or many-to-one
[mapping](#mapping) that maps the input domain of the problem to the
output range of all possible solutions.  If the function is
not computable, then the mapping is not possible.  XXX picture

## Machine

A machine is a real or virtual automaton that is capable of executing
code.  A machine is capable of deciding some subset of the
[computable functions](#computable-function). 

Because code is data and data is code, a machine might be any piece of
data that can be interpreted as a machine.  For example, a virtual
machine's disk image is a piece of data that can be interpreted as a
machine by a VM runtime.  

Any program written in a high-level language is a piece of data that
can be interpreted as a virtual machine by the language's runtime. For
example, the a python program is a piece of data that can be
interpreted as a virtual machine by the python runtime.  

A machine-language program compiled for execution on a UNIX or Linux
system is a piece of data that can be interpreted as a virtual machine
by the kernel of that operating system.  

A machine-language program compiled for bare-metal execution is a
piece of data that can be interpreted as a virtual machine by the
physical CPU.  

Likewise, a disk image of a physical machine is a piece of data that
can be interpreted as a virtual machine by the physical CPU.

In all cases, a machine is an automaton that is capable of executing
code in a language that it recognizes -- see Chomsky's hierarchy for
theoretical examples.  In practice, a python or C application program
is a machine that recognizes the inputs to the program, an operating
system is a machine that recognizes programs compiled for that OS, and
a CPU is a machine that recognizes machine-language programs compiled
for that hardware.

Taking this concept to its logical conclusion, a machine can be
anything from a subatomic particle to an entire universe.  In all
cases, the laws of physics are the language that any machine
recognizes.  See Mark Burgess' work on Promise Theory and
SmartSpacetime for a framework for reasoning about machines at all
levels of abstraction.  [promise][promise] [smartspacetime][smartspacetime]

## Scalable

A scalable system is one that has the ability to adapt its capacity to
meet growth in demand.  

## Fault-tolerant

A fault-tolerant system is one that can continue operation, possibly
at a reduced level, rather than completely failing when some part of
the system fails.

## Self-healing

A self-healing system is one that can perceive that it is not
operating correctly and, without external intervention, make the
necessary adjustments to restore itself to normal operation.

## Decidable

A decidable problem is a (decision problem)[#decision-problem] that
can be solved algorithmically by a (machine)[#machine].  

## Undecidable

An undecidable problem is a (decision problem)[#decision-problem] that
cannot be solved algorithmically by a [machine](#machine).  The
halting problem is an example of an undecidable problem.


computable function versus decidable problem?

In computer science and computability theory, a computable function is one where there exists an algorithm that can output the value of the function for any possible input. The algorithm doesn't have to be efficient but it has to always produce the correct result eventually for any legal input.

A decidable problem, on the other hand, is a decision problem (a problem with a yes or no answer) for which an algorithm exists that can provide the correct answer for any input. In other words, it's a problem that can be solved algorithmically. 

So while computable functions are about producing a specific output for a given input, decidable problems are about determining whether a statement is true or false. If there's an algorithm that can solve the problem in finite time, then the problem is decidable, or the function is computable. 

An additional note is that all problems solvable by a Turing machine (which defines the concept of computability) are decidable, but there are problems that are undecidable, meaning no algorithm can solve them, such as the halting problem you mentioned in the context above.

## Turing machine

A Turing machine is a hypothetical machine that can solve any
decidable problem.  It is the theoretical basis for all modern
computing systems. [turing][turing]


## Turing-complete

A system is said to be "Turing-complete" if it has the ability to
simulate a Turing machine. This means that, given enough time and
resources, it has the capability to perform any computation that a
Turing machine can.

## Turing-equivalent

A system is said to be "Turing-equivalent" if it can not only simulate
a Turing machine but can also be simulated by a Turing machine. This
implies that a Turing-equivalent system and a Turing machine are
essentially equal in terms of their computational powers and
abilities.

While all Turing-equivalent systems are
[Turing-complete](#turing-complete), not all Turing-complete systems
are necessarily Turing-equivalent because they may not be able to be
simulated by a Turing machine.


## Hypercomputer

A hypercomputer is a hypothetical machine that can solve supertask
problems, which are problems that ordinary Turing machines cannot
solve.

## Decentralized Computing

Decentralized computing is a computing paradigm where the computing
resources are distributed across multiple [machines](#machine). The
machines communicate with each other to coordinate their work.  The
machines are not centrally controlled or coordinated.  The machines
need not belong to a single organization or be owned by a single
person or other legal entity.  Examples of decentralized computing
include the Internet, the World Wide Web, and the Bitcoin network.

Distributed computing is a subset of decentralized computing.  In
distributed computing, the machines might be owned by multiple legal
entities, but are centrally controlled or coordinated by one.
Examples of distributed computing include the Google search engine,
the Amazon Web Services cloud computing platform, and the SETI@home
and folding@home distributed computing projects.

## Decentralized Storage

Decentralized storage is a necessary component of [decentralized
computing](#decentralized-computing).  Decentralized storage is a
storage paradigm where the data is distributed across multiple
[machines](#machine). The machines communicate with each other to
ensure that the data is replicated, and that loss or disruption of any
single legal entity or physical location does not result in loss of
data.  

## Participant

A participant is any entity that can interact with a
[machine](#machine). A participant can be a human, another
[machine](#machine), or any other creature, system, or thing, real or
virtual, that is capable of influencing or being influenced by the
behavior or state of another participant.

## Decentralized Community

A decentralized community is a community that is composed of multiple
[participants](#participant) who are not centrally controlled or
coordinated.  The participants need not be members of a single
organization or legal entity.  Examples of decentralized communities
include certain social and cultural movements, the global community of
scientists, and the global community of software developers.

How well a decentralized community can coordinate its efforts depends
on the [decentralized systems](#decentralized-computing) that the
participants rely on.  For example, the Internet and the World Wide
Web have enabled the global community of software developers to
coordinate their efforts to create the Linux operating system, the
Apache web server, the Python programming language, and many other
software projects, including much of the software and hardware that
makes up the Internet itself.  The Internet in turn enabled the global
community of bioinformatics researchers to coordinate their efforts to
sequence the human genome, and more recently to sequence the genomes
of the SARS-CoV-2 coronavirus and other organisms.

## WebAssembly

WebAssembly (WASM) is a binary instruction format for a stack-based
virtual machine.  WASM is designed as a portable target for
compilation of high-level languages like C/C++/Rust, enabling
deployment on the web for both new and legacy client and server
applications.  WASM is a key enabling technology for [decentralized
computing](#decentralized-computing) and [decentralized
communities](#decentralized-community). 

The WASM VM is now included in all major web browsers, including
Chrome, Firefox, Safari, and Edge.  WASM is roughly a replacement for
JavaScript, with the advantage of much better performance.  User
experience is similar to that of JavaScript: the user clicks on a link
to a web page, and the browser executes any code on the page, whether
it is JavaScript or WASM.  Users may not even be aware that they are
using a WASM-based application.  Due to its inclusion in web browsers,
WASM is already the most widely deployed virtual machine in the world.

It's fair to say that WASM delivers on the original promise of Java:
"write once, run anywhere".  Java failed to deliver on that promise
due to both legal and technical issues, including the need for users
to install a Java VM on their systems.  

JavasScript succeeded where Java failed, in large part because it was
included in all web browsers starting with Netscape Navigator in 1995.
But JavaScript has its own issues, including poor performance and a
need for developers to write or rewrite their applications in
JavaScript.

WASM is the next step in the evolution of not just web applications,
but software in general.  Prior to WASM, getting better performance
than JavaScript usually meant writing an application in a compiled
language like C, C++, Go, or Rust, and then compiling it to binary
code for every target platform -- Windows, Linux, MacOS, Android, iOS,
etc. Users would then have to find, download, and install the
application on their systems.  With WASM, developers can write
applications in any language, compile them once, and users can then
run those applications in a web browser by clicking on a link.

This portability has led to WASM's adoption as a general-purpose
portable virtual machine.  WASM is now used to run applications in
many non-web environments; standalone runtimes are available in most
major programming languages.  

See [wasm][wasm] for more information.


## Community Systems Working Group 

The Community Systems Working Group (CSWG) is a group of
[participants](#participant) who are working together to create
[decentralized systems](#decentralized-computing) that will enable
[decentralized communities](#decentralized-community) to coordinate
their efforts to solve the world's most pressing problems.  [cswg][cswg]

CSWG began as a spin-out from the Nation of Makers [nom][nom]
organization in 2022.  The Nation of Makers was in a project of the
Obama Administration's Office of Science and Technology Policy (OSTP)
in 2016.  

It is the intent of the CSWG to become self-hosting; that is, to use
the decentralized systems that it creates to coordinate its own
efforts.  The CSWG is currently bootstrapping this process by using
existing centralized systems such as GitHub, Slack, and Google Docs.

## Host

A host is a single computer.  In this document, we use the term "host" 
instead of "machine" to avoid confusion with the more general concept
of a [machine](#machine).

## Runtime

The XXX runtime is a log-based virtual machine.  In this machine
type, the state of the machine is represented as a node in a directed
acyclic graph (DAG).  Each state change is represented as a transition
from one node to another.  Each node contains the changes caused by
the preceeding transition.  

## Transition function

A transition function:
- takes a state node as input and returns a new state node as output
- is deterministic; that is, given the same input, it will always
  return the same output  
- is pure; that is, it has no side effects that are not represented in
  the output state node
- is referentially transparent; that is, the transition function and
  its input can be replaced by its output without changing the
  content of future state nodes




This machine type can be thought of as a variant of:
- XXX pick or combine
- a Belt machine in which the belt length is unbounded
- a Turing machine in which the entire known universe is the tape
- a Turing machine universe at a given point in time is on one cell of
  the tape.  
- a lambda calculus machine in which the entire known universe is the
  lambda expression



The runtime has access to both the state graph and the chunk cache.
The runtime can facilitate I/O with external actors by providing a
standard interface for external actors to interact with the state
graph and the chunk cache.

## Sandbox

The sandbox is the puplang runtime on a single pup host.  

## Address

An address is derived from the hash of content.  The content is stored in the state
graph.  The content is immutable.

## State Graph

The state graph is a set of [mappings](#mapping).  Each mapping
describes a one-way [transition function](#transition-function) from
one [state](#state) to another state.  The state graph is immutable.

The state graph is a checkpoint of log messages.  We only store
[addresses](#address) in the state graph.  We do not store content in
the state graph.

## Side Effects

A side effect is any change to external state.  External state is
anything outside of the state graph and the chunk cache.  

We make the assumption that exchanging [addresses](#address) with
external actors does not cause side effects, while exchanging content
does.

## Content

Content is stored in the [chunk cache](#chunk-cache).  Content is
immutable.

## Log Messages

Log messages are immutable.  Each log message includes the [address](#address) of
one or more previous log messages.  

## Primitives

Primitives are the lowest level of abstraction in PUP.  PUP's
primitives include [log messages](#log-messages), [addresses](#address),
and [content](#content).

# PUP


PUP (Pattern for Universal Protocols) is the working name for the set
of protocols and standards that the
[CSWG](#community-systems-working-group) is developing.  

## PUPgrid

PUPgrid is the name of the [decentralized
system](#decentralized-computing), based on the [PUP](#pup) protocols
and standards, that the [CSWG](#community-systems-working-group) is
developing.

PUPgrid is designed to be a "universal substrate" for other
decentralized systems, supporting applications ranging from
business systems to social networks and personal productivity tools.
Existing applications can be ported to PUPgrid, often by recompiling
them into [WASM](#webassembly) modules.

Migrating existing applications to PUPgrid will enable them to take
advantage of the resilience and scalability of decentralized systems,
and to more easily interoperate with other PUPgrid applications.

## Puplang

Puplang is the programming language of [PUP](#pup).  Puplang is
explicitly designed to execute on [decentralized
systems](#decentralized-computing) in support of [decentralized
communities](#decentralized-community).  

Puplang is designed to be used as a "glue language" to tie together
other languages and systems, as well as a file and data format for
storing and transmitting data.

As part of the design goal of explicitly supporting communities,
puplang is designed to recognize and interact with individual human
[participants](#participant) in its syntax and semantics.  This is in
contrast to other programming languages which leave the job of
representing human participants to the application programmer.

Puplang is is also unique in that it is designed to serve as an
"assembly language for the universe" -- that is, it uses an address
space that is large enough to address every atom in the universe
several times over.  This 

support the creation of [decentralized systems](#decentralized-computing)
that can scale to any arbitrary size.
As part of this design goal, puplang is designed to be used as a








# References



[mapping]: <https://en.wikipedia.org/wiki/Map_(mathematics)> 'Wikipedia. "Map (mathematics)."' 


[promise]: <http://dx.doi.org/10.1007/11568285_9> 'An Approach to Understanding Policy Based on Autonomy and Voluntary Cooperation. Mark Burgess. DOI: 10.1007/11568285_9'




[turing]: <https://londmathsoc.onlinelibrary.wiley.com/doi/abs/10.1112/plms/s2-42.1.230> 'Turing, A.M. (1937). "On Computable Numbers, with an Application to the Entscheidungsproblem". Proceedings of the London Mathematical Society. 2. 42 (1): 230–65. doi:10.1112/plms/s2-42.1.230.'

XXX the state graph contains only addresses of:
- old state
- new state

XXX the state graph contains only addresses of:
- old state
- transition function
- operand set
- new state

XXX the state graph contains only:
- old state address
- address of transition chunk

XXX messaging protocol is not pub/sub.  It is sub/ad/pull.  
XXX the only side effect is a pull

XXX messaging protocol is not pub/sub.  It is sub/ad/order/hereis.
XXX the only side effect is a hereis

XXX messaging protocol is not pub/sub.  It is sendme(state, func, args)/hereis(state, func, args, result).
- but we need to account for side effect. So it is sendme(state, func, args)/hereis(state, func, args, result, sideeffect).  XXX nope -- need to handle address-only and chunk messages separately


