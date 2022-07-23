package stats

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yabosh/go-clock"
)

//
// RateCounter tests
//

func Test_rc_get_histogram_no_gaps(t *testing.T) {
	defer clock.ResetClock()

	// given a rate counter with samples over multiple seconds
	clock.SetFixedClock()
	clock.AdvanceSeconds(1)
	rc := NewRateCounter(5)
	rc.Mark()
	rc.Mark()
	rc.Mark()
	clock.AdvanceSeconds(1)
	rc.Mark()
	rc.Mark()
	clock.AdvanceSeconds(1)
	rc.Mark()
	clock.AdvanceSeconds(1)

	// When a histogram is retrieved
	hist := rc.GetHistogram(5)

	// then it should be empty
	assert.Equal(t, []int{0, 0, 3, 2, 1}, hist, "Response expected: [0,0,3,2,1]")
}

func Test_rc_get_histogram_gaps(t *testing.T) {
	defer clock.ResetClock()

	// given a rate counter with samples over multiple seconds
	clock.SetFixedClock()
	clock.AdvanceSeconds(1)
	rc := NewRateCounter(5)
	rc.Mark()
	rc.Mark()
	rc.Mark()
	clock.AdvanceSeconds(3)
	rc.Mark()
	rc.Mark()
	clock.AdvanceSeconds(1)

	// When a histogram is retrieved
	hist := rc.GetHistogram(5)

	// then it should be empty
	assert.Equal(t, []int{0, 3, 0, 0, 2}, hist, "Response expected: [0,0,3,2,1]")
}

func Test_rc_get_histogram_wrap_around(t *testing.T) {
	defer clock.ResetClock()

	// given a rate counter with samples over multiple seconds
	clock.SetFixedClock()
	clock.AdvanceSeconds(1)
	rc := NewRateCounter(5)
	rc.Mark()
	rc.Mark()
	rc.Mark()
	clock.AdvanceSeconds(7)
	rc.Mark()
	rc.Mark()
	clock.AdvanceSeconds(1)
	rc.Mark()
	clock.AdvanceSeconds(1)

	// When a histogram is retrieved
	hist := rc.GetHistogram(5)

	// then it should be empty
	assert.Equal(t, []int{0, 0, 0, 2, 1}, hist, "Response expected: [0,0,3,2,1]")
}

func Test_rc_get_histogram_larger_than_sample_size(t *testing.T) {
	defer clock.ResetClock()

	// given a rate counter with samples over multiple seconds
	clock.SetFixedClock()
	clock.AdvanceSeconds(1)
	fmt.Println(clock.Now().Format(time.RFC3339))
	rc := NewRateCounter(5)
	rc.Mark()
	rc.Mark()
	clock.AdvanceSeconds(1)
	rc.Mark()
	clock.AdvanceSeconds(1)

	// When a histogram is retrieved
	hist := rc.GetHistogram(10)

	// then it should be empty
	assert.Equal(t, []int{0, 0, 0, 0, 0, 0, 0, 0, 2, 1}, hist, "Response expected: [0,0,3,2,1]")
}
