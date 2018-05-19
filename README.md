# Telnet Chat

[![CircleCI](https://circleci.com/gh/voidfiles/telnet_chat.svg?style=svg)](https://circleci.com/gh/voidfiles/telnet_chat)

A Simple chat server, a toy.

## Implmentation Details

To start I want to make sure I have all the basics covered.

* Config parsing
* Logging
*
Next I want to seperate out the buisness logic of running
at chat server from the connectivity. I can imagine a future choice to use
something other than telnet.

* ChatRoom
  - This should be in charge of the history of messages
  - This should accept and distribute messages to interested parties
* ClientPool
  - This should accept new network Connections
  - It should cleanup broken connections
  - It should listen for inbound messages from Clients and distribute them
  - It should listen for outbound messages and fan them out to Clients
* TelnetServer
  - This should listen for new connections and send them to the ClientPool
* ChatServer
  - This should create an inbound message chan
  - This should create an outbound message chan
  - This should create a net conn chan
  - It should create a ClientPool and wire it up to the inbound,outbound, and new conn chan
  - It should create a ChatRoom and wire it up to the inbound,outbound chans
  - It should create a TelnetServer and wire it up to the new conn chan
  - It should then Run the ChatRoom and wait until something dies

## Getting Started

This will install dependencies. You'll be able to run tests and
build the codebase afterwards.

```bash
make init
```

Run the tests

```bash
make test
```

To build

```bash
make build
```
