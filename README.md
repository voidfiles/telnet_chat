# Telnet Chat

[![CircleCI](https://circleci.com/gh/voidfiles/telnet_chat.svg?style=svg)](https://circleci.com/gh/voidfiles/telnet_chat)

A Simple chat server, a toy.

## Run

After building you can run the command like this.

```
_work/telnet_chat_darwin_amd64 -configpath ./examples/config.toml
```

or you can use

```
make run
```

You can then telnet like this

```
telnet localhost 6000
```

You can also use the HTTP server with [httpie](https://httpie.org/).

```
http POST http://localhost:8000/api/messages client_id:=1 text="A message"
http GET http://localhost:8000/api/messages
```

## Config

Here is an example config file.

```
telnetIp = "localhost"
telnetPort = "6000"
httpIp = "localhost"
httpPort = "8000"
logPath = "/tmp/logfile.json"
```

## Design

**ChatRoom**:
- Stores a histroy of messages
- Methods for adding messages and retrieving histroy
- Communicates via channels

**ClientPool**:
- Maintains client connections
- Distributes messages to all clients

**ChatServer**:
- Listens for inbound telnet connections
- Passes them to the ClientPool

**HttpServer**:
- Exposes an http interface to the ChatRoom

## Implmentation Details

To start I want to make sure I have all the basics covered.

* Config parsing
* Logging

Next I want to seperate out the buisness logic of running a chat server from the connectivity. I can imagine a future choice to use something other than telnet.

* ChatRoom
  - This should be in charge of the history of messages
  - This should accept and distribute messages to interested parties
* ClientPool
  - This should accept new network Connections
  - It should cleanup broken connections
  - It should listen for inbound messages from Clients and distribute them
  - It should listen for outbound messages and fan them out to Clients
* ChatServer
  - This should listen for new connections and send them to the ClientPool
* HttpServer
  - Should hook into methods provided by ChatRoom
  - To Expose message history
  - To allow
* Main
  - This should create an inbound message chan
  - This should create an outbound message chan
  - It should create a ClientPool and wire it up to the inbound and outbound chans
  - It should create a ChatRoom and wire it up to the inbound and outbound chans
  - It should create a ChatServer and wire it up to the ClientPool
  - It should then run the ClientPool, ChatRoom and ChatServer

## Considerations

If this were going to be something I would put into production I would consider a few things.

* Persistance
  - All the data is lost the second the server shuts down
  - If it was running on a single instance I might consider somethign like boltdb
* Bounds
  - I haven't fully vetted the project for bounds
  - I think their might be a few unbounded operations
    - Specificaly the client fan out
* E2E tests
  - I think I could fully vet the service e2e
  - I'd want to build a system to do that

## Getting Started

This will install dependencies. You'll be able to run tests and build the codebase afterwards.

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

# Notes

Writing a telnet server is new for me. So, I used [golang-chat](https://github.com/kljensen/golang-chat) as inspiration for how to create a telnet server.
