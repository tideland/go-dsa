// Tideland Go Data Structures and Algorithms - Map/Reduce
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mapreduce // import "tideland.dev/go/dsa/mapreduce"

//--------------------
// IMPORTS
//--------------------

import (
	"hash/adler32"
	"runtime"
)

//--------------------
// IDENTIFIABLE
//--------------------

// Identifiable has to be implemented by the data handled
// by map/reduce.
type Identifiable interface {
	// ID returns the identifier for the mapping.
	ID() string
}

// IdentifiableChan is a channel for the transfer of identifiable data.
type IdentifiableChan chan Identifiable

// Close the channel for identifiable data.
func (c IdentifiableChan) Close() {
	close(c)
}

//--------------------
// MAP/REDUCE
//--------------------

// MapReducer has to be implemented to control the map/reducing.
type MapReducer interface {
	// Input has to return the input channel for the
	// date to process.
	Input() IdentifiableChan

	// Map maps a key/value pair to another one and emits it.
	Map(in Identifiable, emit IdentifiableChan)

	// Reduce reduces the values delivered via the input
	// channel to the emit channel.
	Reduce(in, emit IdentifiableChan)

	// Consume allows the MapReducer to consume the
	// processed data.
	Consume(in IdentifiableChan) error
}

// MapReduce applies a map and a reduce function to keys and values in parallel.
func MapReduce(mr MapReducer) error {
	mapEmitChan := make(IdentifiableChan)
	reduceEmitChan := make(IdentifiableChan)

	go performReducing(mr, mapEmitChan, reduceEmitChan)
	go performMapping(mr, mapEmitChan)

	return mr.Consume(reduceEmitChan)
}

//--------------------
// PRIVATE
//--------------------

// closerChan signals the closing of channels.
type closerChan chan struct{}

// closerChan closes given channel after a number of signals.
func newCloserChan(kvc IdentifiableChan, size int) closerChan {
	signals := make(closerChan)
	go func() {
		ctr := 0
		for {
			<-signals
			ctr++
			if ctr == size {
				kvc.Close()
				close(signals)
				return
			}
		}
	}()
	return signals
}

// performReducing runs the reducing goroutines.
func performReducing(mr MapReducer, mapEmitChan, reduceEmitChan IdentifiableChan) {
	// Start a closer for the reduce emit chan.
	size := runtime.NumCPU()
	signals := newCloserChan(reduceEmitChan, size)

	// Start reduce goroutines.
	reduceChans := make([]IdentifiableChan, size)
	for i := 0; i < size; i++ {
		reduceChans[i] = make(IdentifiableChan)
		go func(in IdentifiableChan) {
			mr.Reduce(in, reduceEmitChan)
			signals <- struct{}{}
		}(reduceChans[i])
	}

	// Read map emitted data.
	for kv := range mapEmitChan {
		hash := adler32.Checksum([]byte(kv.ID()))
		idx := hash % uint32(size)
		reduceChans[idx] <- kv
	}

	// Close reduce channels.
	for _, reduceChan := range reduceChans {
		reduceChan.Close()
	}
}

// Perform the mapping.
func performMapping(mr MapReducer, mapEmitChan IdentifiableChan) {
	// Start a closer for the map emit chan.
	size := runtime.NumCPU() * 4
	signals := newCloserChan(mapEmitChan, size)

	// Start map goroutines.
	mapChans := make([]IdentifiableChan, size)
	for i := 0; i < size; i++ {
		mapChans[i] = make(IdentifiableChan)
		go func(in IdentifiableChan) {
			for kv := range in {
				mr.Map(kv, mapEmitChan)
			}
			signals <- struct{}{}
		}(mapChans[i])
	}

	// Dispatch input data to map channels.
	idx := 0
	for kv := range mr.Input() {
		mapChans[idx%size] <- kv
		idx++
	}

	// Close map channels.
	for i := 0; i < size; i++ {
		mapChans[i].Close()
	}
}

// EOF
