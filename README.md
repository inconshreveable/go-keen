# Keen IO golang client SDK [![godoc reference](http://godoc.org/gopkg.in/inconshreveable/go-keen.v0?status.png)](http://godoc.org/gopkg.in/inconshreveable/go-keen.v0)

### Community-Supported SDK
This is an _unofficial_ community supported SDK. If you find any issues or have a request please post an [issue](https://github.com/inconshreveable/go-keen/issues).

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
        keenClient := &keen.Client{ WriteKey: "XXX", ProjectId: "XXX" }
        keenBatchClient := keen.NewBatchClient(keenClient, keenFlushInterval)
        keenBatchClient.AddEvent("collection_name", &ExampleEvent{
            UserId: 102,
            Amount: 39,
            Type: "ball",
            Tags: []string{ "red", "bouncy" },
        })
}
```

## Simple queries

For single queries, use the ```keen.Query()``` interface. For example: 

```go
package main

import (
    "fmt"
    keen "github.com/inconshreveable/go-keen"
)

type Query struct {
    EventCollection string  `json:"event_collection,omitempty"`
    TargetProperty  string  `json:"target_property,omitempty"`
    Timeframe       string  `json:"timeframe,omitempty"`
    GroupBy         string  `json:"group_by,omitempty"`
    Interval        string  `json:"interval,omitempty"`
    Percentile      float64 `json:"percentile,omitempty"`
}

func main() {
    var client = &keen.Client{ReadKey: "XXX", ProjectId: "YYY", WriteKey: "ZZZ"}

    resp, err := client.Query("count", Query{EventCollection: "<event_collection>", Timeframe: "this_14_days", Interval: "daily"})
    if err != nil {
        panic(err)
    }

    fmt.Printf("count: %v\n", resp)
}
```

## TODO
Add support for all other Keen IO API endpoints, especially querying data.


## LICENSE
MIT
