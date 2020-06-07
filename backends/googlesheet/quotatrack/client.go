// Package quotatrack makes it easy to track quota usage over the last x seconds
package quotatrack

import (
	"container/list"
	"time"
)

// Quota Tracks the count of something over a period of the last x time
type Quota struct {
	// Every histroy entry represents a quota hit at that time
	history *list.List
	// The time back in history to look for past entries
	backcheckTime time.Duration
}

// New Creates a new quota counting the uses in the last backcheckTime
func New(backcheckTime time.Duration) *Quota {
	return &Quota{
		history:       list.New(),
		backcheckTime: backcheckTime,
	}
}

// Usage returns the number of consumptions over the last backcheckTime
func (q *Quota) Usage() uint64 {
	usage := uint64(0)
	h := q.history.Front()
	for h != nil {
		if h.Value.(time.Time).After(time.Now().Add(-q.backcheckTime)) {
			usage++
			h = h.Next()
		} else {
			// Remove this entry because we don't need to track it
			toRemove := h
			h = h.Next()
			q.history.Remove(toRemove)
		}
	}

	return usage
}

// Consume consumes "count" entries from the quota at this point in time
func (q *Quota) Consume(count uint64) {
	for i := uint64(0); i < count; i++ {
		q.history.PushBack(time.Now())
	}
}

// TimeUntilQuotaAvailable returns time until "count" quota is available given a limit of "limit" in our backcheckTime
func (q *Quota) TimeUntilQuotaAvailable(limit, count uint64) time.Duration {
	// If the limit is less than the count, this will never be available, return -1
	if limit < count {
		return -1
	}

	// Keep going through oldest entries until we get enough count
	h := q.history.Front()

	// First remove any old entries
	q.Usage()

	timeNeeded := time.Second * 0
	counted := uint64(0)
	for h != nil {
		counted++
		// Add time for this to expire
		timeNeeded += h.Value.(time.Time).Sub(time.Now().Add(-q.backcheckTime))

		if counted >= count {
			break
		}
	}

	return timeNeeded
}
