# Pattern for Universal Protocols (PUP)

# Abstract

This document describes a decentralized computing system. The system
combines simple concepts to support decentralized consensus,
coordination across multiple timelines, and a cosmos-sized address
space. The UI is a decentralized versioning filesystem that merges
standard storage capabilities with special files that execute
Turing-complete computations when accessed. Underlying this is a
pub/sub messaging system with a decentralized cache storing
hash-chained messages indefinitely. The combination of these features
provides a large, decentralized multi-user machine. The system's
decentralized design inherently makes it resilient to the failure of
any particular legal entity or hardware component.

# Overview

PUP, puplang, pupd, and pupgrid are components of a decentralized
computing system. The system's goal is to support a better
collaborative understanding of our universe, ranging from everyday
computations like accounting or document editing to the complexities
of organizational or community governance and the uncertainties
inherent in modern spacetime physics. The system builds on simple
concepts to support a cosmos-sized address space, facilitate
decentralized consensus, and coordinate actions among multiple parties
and timelines.

The system is inherently resilient to failure of any single legal
entity or hardware component.

The system's user interface is implemented as a decentralized
versioning filesystem. This filesystem combines conventional storage
capabilities with special files or directories (inodes) capable of
Turing-complete computation when accessed. Inode functions can be
programmed in any language and transparently executed on remote
hardware.

The filesystem's underlying layer is a pub/sub messaging system;
message topics match filesystem directory paths. The messages
published in each topic are stored indefinitely in cache space across
the system.  Storage is in the form of hash-chained messages that are
in concept similar to log-structured filesystem records.

Combined, these features provide a computing facility that is similar to a
very large decentralized multi-user machine.

# PUP

PUP (Pattern for Universal Protocols) is a blueprint for composing a
decentralized computing system using different systems, languages, and
machine architectures. 

PUP concepts include:

- address: hash-based reference to a function or data item
- syscode: address that references a physical or virtual machine 
- statecode: address that references a machine state
- opcode: address that references a function
- message: basic unit of communication and storage
- puplang: message syntax and grammar
- pub/sub fabric:  message distribution system
- cache: decentralized store of prior messages
- test cases: method of detecting consensus
- filesystem: built on pub/sub topics
- special files: execute functions when accessed

## Address

PUP addresses are content-addressable unique identifiers that
reference functions and data.  An address consists of an
integer that references the hash algorithm, followed by the hash
itself.  If rendered in human-readable text, there is nromally a colon
between the integer and the hash, with no spaces.  

As of this writing, the hash algorithm assigned numbers are:

| i | algorithm common name |
|---|-----------------------|
| 0 | md5                   |
| 1 | sha1                  |
| 2 | sha256                |
| 3 | sha512                |

For example, if an address uses the sha256 hash algorithm, then the
human-readable rendering of the address might look like this:

```
2:a5a5318e7a548a20e755e9ccd29a3eaa1de3b590659fb59c760d693aa4fe0bb1
```

## Syscode

XXX


PUP (Pattern for Universal Protocols) system. These addresses are
based on a cryptographic hash of the precise description of the
functions. This mechanism of addressing ensures that functions remain
unambiguous, even if they get executed in different languages or
systems. It is an important part of PUP's decentralized and efficient
communication framework.

## PUP messages

A PUP message is the fundamental unit of communication. The first line
of a PUP message is the syscode.  A syscode
always starts with a hash that references the function to be used to
parse the remainder of the message.  The hash is composed of an
integer that references the hash algorithm, followed by a colo the hash
itself 
header format for a PUP message:


show me the header format for an MQTT message

The MQTT message utilizes a structured format that includes:

1. A fixed header: It is the least amount of information necessary for all messages. The first byte of this header is created by a "message type" coupled with some flags concerning the message. The subsequent bytes render a "remaining length", which details the number of bytes related to this message.

2. A variable header: This section holds optional protocol metadata which varies based on the message type. For instance, for a "CONNECT" message, the variable header would include the protocol name ("MQTT"), the protocol level, connect flags, and a keep alive timer.

3. The payload: It contains the actual data. For a "CONNECT" message it would include payload data transferred from Client to Server as part of the "CONNECT" message (like client identifier, will topic, username and password etc.), and vice versa.

Please be aware that an MQTT message does not always include these three parts – some messages might lack the "variable header" and the "payload", consisting solely of the "fixed header".


imagine the header format for a PUP message.  it will be similar to
the header format for an MQTT or NATS message

The header format for a PUP (Pattern for Universal Protocols) message could likely include a pre-defined sequence of fields, similar to MQTT or NATS, which could possibly include the following:

1. **Message ID**: A unique identifier for each message - this could potentially be a hash, to ensure all messages are unique.

2. **Timestamp**: The time at which the message was created.

3. **Source and Destination Addresses**: Where the message is coming from and where it's intending to go.

4. **Payload Length**: This indicates the size of the subsequent payload to enable the receiver to correctly parse it.

5. **Hash/Index**: As indicated in the context, a PUP message starts with a hash that includes essential information such as the machine, format, and other details.

Please note, the actual implementation may vary or have additional components based on PUP's design and communication requirements.

A PUP
message starts with a hash that indicates machine, format,
among other details, before the messages content. This helps avoid
compatibility issues as it holds critical information about the format
and function, enabling developers to encode and decode messages
according to target system's needs. 

Although described in simple terms, precise details about the format and content may vary depending on the implementation, usage, and specifics of the PUP-based system.


PUP messages are routed using a pub/sub system across different
systems and applications. Different systems can publish or subscribe
to the topics they are interested in, this subscription model
facilitates inter-component communication. When a message is published
to a specific topic, it's stored in the 'pupd' cache. Nodes that
subscribe to the topic can then retrieve this message from the cache,
enabling asynchronous communication. The pub/sub model is
decentralized and scalable, where each node in the grid could act as a
publisher, subscriber, or both. Messages representing different types
of operations are published to different 'topics', creating a dynamic
and flexible infrastructure enabling direct communication between
nodes without needing to know specifics about the routing paths of the
messages. The pupd cache serves as a temporary storage component that
stores published messages and transmits them to subscribers when
interest in a particular topic is expressed. This all contributes to
the efficiency and scalability of the PUP system allowing for a
seamless data flow.

PUP messages are routed

As a messaging protocol, PUP manages the routing of messages between
different applications and systems.

PUP utilizes a local cache (timeline) at each node, which contains previous function calls and their results. It essentially treats systems like distributed databases, each storing a portion of the global state.

The pub/sub and filesystem aspect is embodied in its structure and interaction model. PUP treats topics akin to filesystem directory paths, with each containing sequential, hash-chained messages, similar to writing changes in a log-structured filesystem. The pub/sub model allows different systems to publish or subscribe to the events they are interested in, which can facilitate communication between different components.

As for inode-based computation, PUP's unique function addresses based on cryptographic hashes ensure efficient data routing and processing in a manner that resonates with the inode concept seen in filesystems. Each function, either local or remote, in the PUP system has a unique address akin to inodes, serving as a fundamental unit of data routing.
  
Also, in terms of decentralization, PUP uses a consensus mechanism to manage conflicts and merge timelines in different nodes to ensure data integrity across the system.

In summation, PUP as a pattern encapsulates these aspects, forming a comprehensive and efficient messaging and communication framework.

Pup is the underlying architecture that drives PupGrid, adhering to the philosophy of 'everything's a file.' It offers a consistent interaction model with different types of data, making it developer-friendly. PupLang, akin to Job Control Language (JCL), is the language responsible for linking various software and hardware functions together. PupD is a daemon process that manages the distributed cache, ensuring data consistency and availability in the network. PupGrid is a distributed, decentralized machine built using Pup, aiming for efficient, flexible, and decentralized communication.

PUP (Pattern for Universal Protocols) is the core component or underlying architecture designed to facilitate seamless interaction between diverse systems, programming languages, and machine architectures. It is the driving force in promoting efficient, flexible, and decentralized communication and managing a universe-sized address space unrivalled by traditional methods.

# puplang

pupLang is akin to the Job Control Language (JCL) and plays a pivotal role in linking together different software and hardware functions. It represents the task-level programming language designed explicitly for the PUP system, linking various aspects of the system into consistent processes.

# pupd

pupd is essentially a daemon process that assists in conducting network functions. It manages the distributed cache, ensuring data availability and consistency across the vast network, crucial for smooth and reliable operations.

# pupgrid

PupGrid is a distributed, decentralized machine constructed using PUP. It aims to create a unified and efficient method of communication across varied systems while supporting decentralized consensus processes. Apart from facilitating regular computations, it's versatile enough to deal with more complex problems encompassing governance and spacetime physics. Overall, PUPGrid offers scalability, resilience, and manageability, making it an appropriate choice for diverse and evolving environments.

# requirements

# axioms

# assumptions

# implementation
