# TCP/IP Protocol

## Introduction

Computer applications need to send data to other computers in order to transform
the information into something useful in the other end of the communication. For
this to be possible, it's key to understand the standards that govern the
internet, a.k.a. the biggest network, or the network of networks.

The concerning standard this time is TCP/IP that allows to send data streams
from one computer or device to another via the internet.

Many considerations must be taken to undertake these designs in real life, like
the 7-layer network architecture [^1], and the error detection and
correction [^2].

We should also understand how data is packet, encoded, transmitted. So for 
example, we can design a protocol that adds a line break `\n` to tell the 
receiver to split the data into tokens delimited by the line feed character.

[^1]: The 7-layer architecture is more academic than real, a more pragmatic approach is taken in real implementations

[^2]: Read my course project [Reed-Muller Codes](https://dev.mathsoftware.engineer/cp-unah-mm544-reed-muller-codes) 

## Protocol Definition

The TCP/IP protocols are fundamental for transferring data over the internet. 
These are detailed below. 
