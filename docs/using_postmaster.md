# Using Postmaster

This guide explains the basics of using Postmaster to manage events from rubykube components.

Nowadays postmaster support two components, that produce events

* [Barong](https://github.com/rubykube/barong)
* [Peatio](https://github.com/rubykube/peatio)

Read more about peatio Event API [here]().

Read more about barong Event API [here]().

## Concepts

An Event is message produces to message broker with defined structure.

```JSON
{
    "record": {
        "user": {

        }
    }
}
```

## Configuration

## Installation

### Locally

go build ./cmd/postmaster/postmaster.go

./postmaster

```
11:22AM INF waiting for events
11:22AM INF successfully connected to amqp://guest:guest@localhost:5672/
```

### Docker

Use official docker image for postmaster quay.io/openware/postmaster, which is automaticly build on each commit in master branch.

```
docker run quay.io/openware/postmaster --name=postmaster
```