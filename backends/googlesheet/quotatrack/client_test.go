package quotatrack

import (
	"testing"
	"time"
)

func TestQuota(t *testing.T) {
	testingBackcheckTime := time.Second
	quota := New(testingBackcheckTime)

	// Consume 5
	quota.Consume(3)
	quota.Consume(2)

	// Check to make sure usage is 5
	if quota.Usage() != 5 {
		t.Errorf("Did not get correct usage amount")
	}

	// Check to make sure the time needed until next quota is around time.Second
	if timeUntil := quota.TimeUntilQuotaAvailable(5, 1); timeUntil < testingBackcheckTime/2 || timeUntil < 0 || timeUntil > testingBackcheckTime {
		t.Errorf("Incorrect next available time")
	}

	// Wait for this quota to expire
	time.Sleep(testingBackcheckTime)

	// Usage should be 0
	if quota.Usage() != 0 {
		t.Errorf("Should be 0 quota usage")
	}

	// Check to make sure we are keeping the entries in the backend
	if quota.history.Len() > 0 {
		t.Errorf("Keeping data stored when it shouldn't")
	}
}
