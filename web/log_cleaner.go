package web

import (
	"cronsun/db/entries"
	"time"

	"cronsun/log"
)

func RunLogCleaner(cleanPeriod, expiration time.Duration) (close chan struct{}) {
	t := time.NewTicker(cleanPeriod)
	close = make(chan struct{})
	go func() {
		for {
			select {
			case <-t.C:
				cleanupLogs(expiration)
			case <-close:
				return
			}
		}
	}()

	return
}

func cleanupLogs(expiration time.Duration) {
	err := entries.ClearJobLogs(expiration)
	if err != nil {
		log.Errorf("[Cleaner] Failed to remove expired logs: %s", err.Error())
		return
	}

}
