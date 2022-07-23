package stats

import "github.com/yabosh/go-clock"

//
// Basic rate counter implementation tuned for performance
//
// Counts actions per second for a given number of seconds
//
// This is a basic design based on a very basic ring-buffer
// and is designed for performance. This design has these limitations
//
//   * Single-threaded operation only
//   * Tracks time in seconds
//
// Additionally, this approach provides a basic histogram
// of event samples for 'x' number of timeframes.
//
// Making the assumption that the counter will only be
// accessed by a single thread dramatically simplifies
// the implementation.
//
// The Filter is single-threaded and it will be interacting
// with the rate counter so the rate counter can also be single threaded
//
// Buckets
// -------
// A bucket tracks a number of samples for a given second in time (based on the Unix epoch)
// Each rate counter has a fixed number of buckets each that represent the number
// of events that occurred during a single second of time.  If there are no samples
// for a time interval then there is no bucket for that interval so there may be gaps
// between the timestamps for buckets.

// Counter counts number of samples for a given timestamp
type Counter struct {

	// Unix epoch timestamp that this counter represents
	Timestamp int64 `json:"timestamp"`

	// Number samples for this timestamp
	Count int `json:"samples"`
}

// RateCounter tracks event buckets
type RateCounter struct {

	// Buckets is a simple ring-buffer used to track samples for a given second in time
	Buckets     []*Counter `json:"buckets"`
	Head        int        `json:"head"`
	Tail        int        `json:"tail"`
	bucketCount int
}

// Used to provide time values for RateCounters
// This value can be overridden in unit tests to provide
// deterministic time values
//var rateCounterTime = time.Now

// NewRateCounter returns a new instance of RateCounter initialized for the
// specified number of buckets
func NewRateCounter(maxBuckets int) RateCounter {

	buckets := make([]*Counter, maxBuckets)

	// Initialize all of the buckets
	for i := 0; i < maxBuckets; i++ {
		buckets[i] = &Counter{0, 0}
	}

	return RateCounter{
		Buckets:     buckets,
		Head:        0,
		Tail:        0,
		bucketCount: maxBuckets,
	}

}

// Mark increases the sample count for the current time bucket.
// If the time has changed since the last sample was taken then
// a new bucket is created in the ring buffer.
func (rc *RateCounter) Mark() {

	tick := clock.Now().Unix()

	if rc.Buckets[rc.Head].Timestamp == tick {
		// Current bucket has same timstamp
		rc.Buckets[rc.Head].Count++
	} else {
		// Create a new bucket for this timestamp
		rc.Head = (rc.Head + 1) % rc.bucketCount

		rc.Buckets[rc.Head].Timestamp = tick
		rc.Buckets[rc.Head].Count = 1

		// If the head of the ring buffer reached the tail then move the tail and 'lose' the data at its location
		if rc.Head == rc.Tail {
			rc.Tail = (rc.Tail + 1) % rc.bucketCount
		}
	}

}

// GetHistogram retrieves rates for the last n
// seconds up to the maximum size of the buckets.
//
// The response is an array with exactly 'n' elements
// with samples from now() - n seconds to now() - 1.  If there
// were no samples during a period then a zero is returned
// for that timeslot.  The first element (0) in the array
// represents the time that is time.Now() - n seconds.
//
// This is used to create a histgram of activity within
// the rate counter
//
func (rc *RateCounter) GetHistogram(numSeconds int) []int {

	// Index rate counter samples by timestamp so we
	// can retrieve samples by timestamp
	// Key: timestamp
	// Value: sample count
	index := make(map[int64]int)

	for i := 0; i < rc.bucketCount; i++ {
		index[rc.Buckets[i].Timestamp] = rc.Buckets[i].Count
	}

	// Start with the 'previous' second because the bucket for the
	// current second may only be partially full.
	start := clock.Now().Unix() - int64(numSeconds)
	response := make([]int, numSeconds)

	for i := int64(0); i < int64(numSeconds); i++ {
		response[i] = index[start+i]
	}

	return response
}
