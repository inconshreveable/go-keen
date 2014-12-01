# Keen IO golang client SDK [![godoc reference](http://godoc.org/gopkg.in/inconshreveable/go-keen.v0?status.png)](http://godoc.org/gopkg.in/inconshreveable/go-keen.v0)

## API Stability

The master branch has no API stability guarantees. You can import the latest stable API with:
```go
import "gopkg.in/inconshreveable/go-keen.v0"
````

## Writing Events

This is the very beginnings of a Keen IO client SDK in Go. Currently, only adding events to collections is supported.

The simplest API is to create a client object and then call AddEvent:
```go
package main

import (
        "github.com/inconshreveable/go-keen"
)

type ExampleEvent struct {
        UserId int
        Amount int
        Type string
        Tags []string
}

func main() {
        keenClient := &keen.Client{ ApiKey: "XXX", ProjectToken: "XXX" }
        keenClient.AddEvent("collection_name", &ExampleEvent{
                UserId: 102,
                Amount: 39,
                Type: "ball",
                Tags: []string{ "red", "bouncy" },
        })
}
```

## Batch event reporting

For production use, it makes more sense to add events to an internal buffer which is
flushed to Keen at a regular interval in a single batch upload call. The go-keen library provides
a BatchClient which allows you to do just that while keeping the same, simple API for adding
events. Do note that it does mean that you could lose events if your program exits or crashes before it
flushes the events to Keen.
```go
package main

import (
        "github.com/inconshreveable/go-keen"
        "time"
)

const keenFlushInterval = 10 * time.Second

type ExampleEvent struct {
        UserId int
        Amount int
        Type string
        Tags []string
}

func main() {
        keenClient := &keen.Client{ ApiKey: "XXX", ProjectToken: "XXX" }
        keenBatchClient := keen.NewBatchClient(keenClient, keenFlushInterval)
        keenBatchClient.AddEvent("collection_name", &ExampleEvent{
            UserId: 102,
            Amount: 39,
            Type: "ball",
            Tags: []string{ "red", "bouncy" },
        })
}
```

## TODO
Add support for all other Keen IO API endpoints, espeically querying data.
