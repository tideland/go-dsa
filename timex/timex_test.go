// Tideland Go Data Structures and Algorithms - Time Extensions - Unit Tests
//
// Copyright (C) 2009-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package timex_test

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/dsa/timex"
)

//--------------------
// TESTS
//--------------------

// Test time containments.
func TestTimeContainments(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Create some test data.
	ts := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	years := []int{2008, 2009, 2010}
	months := []time.Month{10, 11, 12}
	days := []int{10, 11, 12, 13, 14}
	hours := []int{20, 21, 22, 23}
	minutes := []int{0, 5, 10, 15, 20, 25}
	seconds := []int{0, 15, 30, 45}
	weekdays := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday}

	assert.True(timex.YearInList(ts, years), "Go time in year list.")
	assert.True(timex.YearInRange(ts, 2005, 2015), "Go time in year range.")
	assert.True(timex.MonthInList(ts, months), "Go time in month list.")
	assert.True(timex.MonthInRange(ts, 7, 12), "Go time in month range.")
	assert.True(timex.DayInList(ts, days), "Go time in day list.")
	assert.True(timex.DayInRange(ts, 5, 15), "Go time in day range .")
	assert.True(timex.HourInList(ts, hours), "Go time in hour list.")
	assert.True(timex.HourInRange(ts, 20, 31), "Go time in hour range .")
	assert.True(timex.MinuteInList(ts, minutes), "Go time in minute list.")
	assert.True(timex.MinuteInRange(ts, 0, 5), "Go time in minute range .")
	assert.True(timex.SecondInList(ts, seconds), "Go time in second list.")
	assert.True(timex.SecondInRange(ts, 0, 5), "Go time in second range .")
	assert.True(timex.WeekdayInList(ts, weekdays), "Go time in weekday list.")
	assert.True(timex.WeekdayInRange(ts, time.Monday, time.Friday), "Go time in weekday range .")
}

// TestBeginOf tests the calculation of a beginning of a unit of time.
func TestBeginOf(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	ts := time.Date(2015, time.August, 2, 15, 10, 45, 12345, time.UTC)

	assert.Equal(timex.BeginOf(ts, timex.Second), time.Date(2015, time.August, 2, 15, 10, 45, 0, time.UTC))
	assert.Equal(timex.BeginOf(ts, timex.Minute), time.Date(2015, time.August, 2, 15, 10, 0, 0, time.UTC))
	assert.Equal(timex.BeginOf(ts, timex.Hour), time.Date(2015, time.August, 2, 15, 0, 0, 0, time.UTC))
	assert.Equal(timex.BeginOf(ts, timex.Day), time.Date(2015, time.August, 2, 0, 0, 0, 0, time.UTC))
	assert.Equal(timex.BeginOf(ts, timex.Month), time.Date(2015, time.August, 1, 0, 0, 0, 0, time.UTC))
	assert.Equal(timex.BeginOf(ts, timex.Year), time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC))
}

// TestEndOf tests the calculation of a ending of a unit of time.
func TestEndOf(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	ts := time.Date(2012, time.February, 2, 15, 10, 45, 12345, time.UTC)

	assert.Equal(timex.EndOf(ts, timex.Second), time.Date(2012, time.February, 2, 15, 10, 45, 999999999, time.UTC))
	assert.Equal(timex.EndOf(ts, timex.Minute), time.Date(2012, time.February, 2, 15, 10, 59, 999999999, time.UTC))
	assert.Equal(timex.EndOf(ts, timex.Hour), time.Date(2012, time.February, 2, 15, 59, 59, 999999999, time.UTC))
	assert.Equal(timex.EndOf(ts, timex.Day), time.Date(2012, time.February, 2, 23, 59, 59, 999999999, time.UTC))
	assert.Equal(timex.EndOf(ts, timex.Month), time.Date(2012, time.February, 29, 23, 59, 59, 999999999, time.UTC))
	assert.Equal(timex.EndOf(ts, timex.Year), time.Date(2012, time.December, 31, 23, 59, 59, 999999999, time.UTC))
}

// TestRetrySuccess tests a successful retry.
func TestRetrySuccess(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	count := 0
	err := timex.Retry(func() (bool, error) {
		count++
		return count == 5, nil
	}, timex.ShortAttempt())
	assert.Nil(err)
	assert.Equal(count, 5)
}

// TestRetryFuncError tests an error inside the retried func.
func TestRetryFuncError(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	err := timex.Retry(func() (bool, error) {
		return false, errors.New("ouch")
	}, timex.ShortAttempt())
	assert.ErrorMatch(err, "ouch")
}

// TestRetryTooLong tests a retry timout.
func TestRetryTooLong(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	rs := timex.RetryStrategy{
		Count:          10,
		Break:          5 * time.Millisecond,
		BreakIncrement: 5 * time.Millisecond,
		Timeout:        50 * time.Millisecond,
	}
	err := timex.Retry(func() (bool, error) {
		return false, nil
	}, rs)
	assert.ErrorMatch(err, ".* retried longer than .*")
}

// TestRetryTooOften tests a retry count error.
func TestRetryTooOften(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	rs := timex.RetryStrategy{
		Count:          5,
		Break:          5 * time.Millisecond,
		BreakIncrement: 5 * time.Millisecond,
		Timeout:        time.Second,
	}
	err := timex.Retry(func() (bool, error) {
		return false, nil
	}, rs)
	assert.ErrorMatch(err, ".* retried more than .* times")
}

// EOF
