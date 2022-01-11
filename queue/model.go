package queue

import (
	"sync"
	"time"
)

type (
	Metadata struct {
		Timestamp time.Time
		ID        string
	}

	Message struct {
		Payload  interface{}
		MetaData Metadata
	}

	gobbit struct {
		subs  map[string]chan *Message
		l     sync.Mutex
		order uint64
	}
)
