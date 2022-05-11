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

### TCP

**TCP** stands for **Transfer Control Protocol**, and it is the standard that
make possible sending large amounts of data over the internet. It can be
implemented by any programmer, and it is the basis of sending data over the
network.

This protocol works with the IP protocol to transport data over the network.

Data are separated into fragments called **packets**, these packets are
transmitted over the network and glued together in the receiver to create the
original information.

Packets are sent via different mediums, some are faster or shorter than others,
and can be traced with a technique called **packet tracer**. Packet tracing is
something done in networking courses, for example, employing the Cisco Packet
Tracer software.

Be careful as TCP works with data streams, so one packet sent does not mean one
packet received.

#### Handshaking

**Handshaking** is a signal that is sent from one computer to another to
establish a contract. So for example, we can start the communication with one
handshake, and end it when we finish the data transaction with another
handshake.

Handshakes are commonly used in web sockets connections. They are important here
to avoid the overhead of sending HTTP requests for each message. The overhead is
only in certain parts (beginning, end) of the communication instead.

### IP

**IP** stands for **Internet Protocol**, and it is the standard to send those
packages to the correct destination address. The currently used version of this
protocol is IPv4, but IPv6 is the future as it allows practically and infinite
amount of addresses.

If we use IPv6, we won't have to give private addresses to local devices with
DHCP, they can have their own IP address instead. IPv4 only supports an octet or
one byte in the following structure **xxxx.xxxx.xxxx.xxxx** so that is the
address of the receiver that the IP protocol will send that data forward.

## TCP vs UDP

The other established protocol **UDP (User Datagram Protocol)** is used to 
send datagrams that allow faster real time communications. They can drop 
information to be able to send data faster, so common applications are video 
meetings or live-streaming when we don't care about some pixels lost.  

TCP has to assemble the packages sent over the network and perform error
correction. Another protocol for error correction is the **SCTP (Stream Control
Transmission Protocol)**

## FTP

**FTP (File Transfer Protocol)** is used to send files easily via TCP. They 
have users with password login, but it is not secure. The **SFTP (SSH File 
Transfer Protocol)** is to be used for secure applications instead.

## Web Socket

Web Socket is the technology used for client-server real-time communication. It
performs handshakes and enables a port to open to establish the communication.  

Web Sockets applications can vary. They're commonly used for chats, social 
media, collaborative tools, multiplayer games, etc.

## Bibliography

- [Designing and Implementing an Application Layer Network Protocol Using UNIX Sockets and TCP \| Ryan Lattrel](https://www.egr.msu.edu/classes/ece480/capstone/fall12/group02/documents/Ryan-Lattrel_App-Note.pdf)
- [TCP/IP \| Definition \& Facts \| Britannica](https://www.britannica.com/technology/TCP-IP)
- [What is TCP/IP? \| Cloudflare](https://www.cloudflare.com/learning/ddos/glossary/tcp-ip)
- [TCP/IP Protocols \| IBM Documentation](https://www.ibm.com/docs/en/aix/7.2?topic=protocol-tcpip-protocols)
- [Handshaking \| Wikipedia](https://en.wikipedia.org/wiki/Handshaking)
- [User Datagram Protocol \| Wikipedia](https://en.wikipedia.org/wiki/User_Datagram_Protocol)
- [Difference Between TCP and FTP \| Web Hosting Prime](https://webhostingprime.com/difference-between-tcp-and-ftp)
- [Reading Files \| Go by Example](https://gobyexample.com/reading-files)
