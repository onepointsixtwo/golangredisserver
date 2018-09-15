package ttltimer

import (
	"time"
)

type TTLTimer struct {
	timer *time.Timer
	end   time.Time
}

func New(seconds int) *TTLTimer {
	duration := time.Duration(seconds) * time.Second
	return &TTLTimer{time.NewTimer(duration), time.Now().Add(duration)}
}

func (timer *TTLTimer) Reset(seconds int) {
	duration := time.Duration(seconds) * time.Second
	timer.timer.Reset(duration)
	timer.end = time.Now().Add(duration)
}

func (timer *TTLTimer) Stop() {
	timer.timer.Stop()
}

func (timer *TTLTimer) GetTimerChannel() <-chan time.Time {
	return timer.timer.C
}

func (timer *TTLTimer) RemainingTTL() int {
	duration := timer.end.Sub(time.Now())
	return int(duration.Seconds())
}
