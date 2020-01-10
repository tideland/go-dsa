// Tideland Go Data Structures and Algorithms - Map/Reduce - Unit Tests
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mapreduce_test

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/dsa/identifier"
	"tideland.dev/go/dsa/mapreduce"
)

//--------------------
// TESTS
//--------------------

// TestMapReduce runs a scenario, where orders are analyzed
// and a list of the analyzed articles will be returned
func TestMapReduce(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	// Create MapReducer and let the show begin.
	mr := &OrderMapReducer{200000, make(map[int][]*OrderItem), make(map[string]*OrderItemAnalysis), assert}
	err := mapreduce.MapReduce(mr)
	// Asserts.
	assert.Nil(err)
	assert.Equal(len(mr.items), len(mr.analyses), "all items are analyzed")
	for _, analysis := range mr.analyses {
		quantity := 0
		amount := 0.0
		discount := 0.0
		items := mr.items[analysis.ArticleNo]
		for _, item := range items {
			unitDiscount := (item.UnitPrice / 100.0) * item.DiscountPerc
			totalDiscount := unitDiscount * float64(item.Count)
			totalAmount := (item.UnitPrice - unitDiscount) * float64(item.Count)

			quantity += item.Count
			amount += totalAmount
			discount += totalDiscount
		}
		assert.Equal(quantity, analysis.Quantity, "quantity per article")
		assert.About(amount, analysis.Amount, 0.01, "amount per article")
		assert.About(discount, analysis.Discount, 0.01, "discount per article")
	}
}

// BenchmarkMapReduce tests the performance of map/reduce.
func BenchmarkMapReduce(b *testing.B) {
	assert := asserts.NewPanic()
	mr := &OrderMapReducer{b.N, make(map[int][]*OrderItem), make(map[string]*OrderItemAnalysis), assert}
	mapreduce.MapReduce(mr)
}

//--------------------
// HELPERS
//--------------------

type OrderMapReducer struct {
	count    int
	items    map[int][]*OrderItem
	analyses map[string]*OrderItemAnalysis
	assert   *asserts.Asserts
}

// Input has to return the input channel for the
// date to process.
func (o *OrderMapReducer) Input() mapreduce.IdentifiableChan {
	input := make(mapreduce.IdentifiableChan)
	// Generate test orders and push them into a the
	// input channel.
	articleMaxNo := 9999
	unitPrices := make([]float64, articleMaxNo+1)
	for i := 0; i < articleMaxNo+1; i++ {
		unitPrices[i] = rand.Float64() * 100.0
	}
	go func() {
		defer close(input)
		for i := 0; i < o.count; i++ {
			order := &Order{
				OrderNo:    identifier.NewUUID(),
				CustomerNo: rand.Intn(999) + 1,
				Items:      make([]*OrderItem, rand.Intn(9)+1),
			}
			for j := 0; j < len(order.Items); j++ {
				articleNo := rand.Intn(articleMaxNo)
				order.Items[j] = &OrderItem{
					ArticleNo:    articleNo,
					Count:        rand.Intn(9) + 1,
					UnitPrice:    unitPrices[articleNo],
					DiscountPerc: rand.Float64() * 4.0,
				}
				o.items[articleNo] = append(o.items[articleNo], order.Items[j])
			}
			input <- order
		}
	}()
	return input
}

// Map maps a key/value pair to another one and emits it.
func (o *OrderMapReducer) Map(in mapreduce.Identifiable, emit mapreduce.IdentifiableChan) {
	order := in.(*Order)
	// Analyse items and emit results.
	for _, item := range order.Items {
		unitDiscount := (item.UnitPrice / 100.0) * item.DiscountPerc
		totalDiscount := unitDiscount * float64(item.Count)
		totalAmount := (item.UnitPrice - unitDiscount) * float64(item.Count)
		analysis := &OrderItemAnalysis{item.ArticleNo, item.Count, totalAmount, totalDiscount}
		emit <- analysis
	}
}

// Reduce reduces the values delivered via the input
// channel to the emit channel.
func (o *OrderMapReducer) Reduce(in, emit mapreduce.IdentifiableChan) {
	memory := make(map[string]*OrderItemAnalysis)
	// Collect emitted analysis data.
	for i := range in {
		analysis := i.(*OrderItemAnalysis)
		if existing, ok := memory[i.ID()]; ok {
			existing.Quantity += analysis.Quantity
			existing.Amount += analysis.Amount
			existing.Discount += analysis.Discount
		} else {
			memory[i.ID()] = analysis
		}
	}
	// Emit the result to the consumer.
	for _, analysis := range memory {
		emit <- analysis
	}
}

// Consume allows the MapReducer to consume the
// processed data.
func (o *OrderMapReducer) Consume(in mapreduce.IdentifiableChan) error {
	for i := range in {
		o.assert.Nil(o.analyses[i.ID()], "each analysis only once")
		o.analyses[i.ID()] = i.(*OrderItemAnalysis)
	}
	return nil
}

// Order item type.
type OrderItem struct {
	ArticleNo    int
	Count        int
	UnitPrice    float64
	DiscountPerc float64
}

// Order type.
type Order struct {
	OrderNo    identifier.UUID
	CustomerNo int
	Items      []*OrderItem
}

// ID returns the order number as identifier.
func (o *Order) ID() string {
	return o.OrderNo.String()
}

func (o *Order) String() string {
	msg := "ON: %v / CN: %4v / I: %v"
	return fmt.Sprintf(msg, o.OrderNo, o.CustomerNo, len(o.Items))
}

// Order item analysis type.
type OrderItemAnalysis struct {
	ArticleNo int
	Quantity  int
	Amount    float64
	Discount  float64
}

// IDy returns the article number as identifier.
func (a *OrderItemAnalysis) ID() string {
	return strconv.Itoa(a.ArticleNo)
}

func (oia *OrderItemAnalysis) String() string {
	msg := "AN: %5v / Q: %4v / A: %10.2f € / D: %10.2f €"
	return fmt.Sprintf(msg, oia.ArticleNo, oia.Quantity, oia.Amount, oia.Discount)
}

// EOF
