package ttltimer

import (
	"testing"
)

func TestTTLTimer(t *testing.T) {
	timer := New(10)
	remaining := timer.RemainingTTL()

	// Bit of a clunky test, but since the remaining TTL is seconds and the time immediately starts
	// ticking down, we would expect it to be equal to 9 when we immedately check the TTL because it'll round down
	if remaining == 9 {
		t.Log("Remaining time was value of 9 seconds as expected")
	} else {
		t.Errorf("Expected time remaining to be greater than 9 but was %v", remaining)
	}
}
