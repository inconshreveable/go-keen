package keen

import (
	"fmt"
	"time"
)

const (
	addEventTimeout = time.Second
)

type eventTuple struct {
	collection string
	payload    interface{}
}

type BatchClient struct {
	*Client
	flushInterval time.Duration
	flush         chan int
	events        chan *eventTuple
}

func NewBatchClient(client *Client, flushInterval time.Duration) *BatchClient {
	c := &BatchClient{
		Client:        client,
		flush:         make(chan int),
		flushInterval: flushInterval,
		events:        make(chan *eventTuple),
	}
	go c.loop()
	return c
}

func (c *BatchClient) AddEvent(collection string, event interface{}) error {
	select {
	case c.events <- &eventTuple{collection: collection, payload: event}:
	case <-time.After(addEventTimeout):
		return fmt.Errorf("Timeout while trying to add event for batch processing")
	}
	return nil
}

func (c *BatchClient) Flush() {
	c.flush <- 1
}

func (c *BatchClient) loop() {
	go func() {
		for _ = range time.Tick(c.flushInterval) {
			c.flush <- 1
		}
	}()

	batch := make(map[string][]interface{})
	for {
		select {

		// add events to the batch
		case e := <-c.events:
			list, ok := batch[e.collection]
			if !ok {
				list = make([]interface{}, 0)
			}
			batch[e.collection] = append(list, e.payload)

		// flush
		case <-c.flush:
			// no metrics to report
			if len(batch) == 0 {
				continue
			}

			// make the request
			c.AddEvents(batch)

			// batch is empty now
			batch = make(map[string][]interface{})
		}
	}
}
