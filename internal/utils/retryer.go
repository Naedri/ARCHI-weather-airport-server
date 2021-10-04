package utils

import (
	"time"
)

type RetryEvent struct {
	time    time.Time
	success bool
	fail    error
}

type Retryer struct {
	retryTimes    int
	maxRetryTimes int
	tickInterval  time.Duration
	timeout       time.Duration
	tick          func()
	actionHandler func()
	errorHandler  func()
}

func (r Retryer) something() {
	//r.MaxRetryTimes
}

// will start the
func tick(r *Retryer) {
	//TODO
}

func newRetryer(maxRetryTimes int, tickInterval time.Duration, timeout time.Duration, actionHandler func(), errorHandler func()) (Retryer, error) {
	retryer := new(Retryer)
	retryer.maxRetryTimes = maxRetryTimes
	retryer.tickInterval = tickInterval
	retryer.timeout = timeout
	retryer.actionHandler = actionHandler
	retryer.errorHandler = errorHandler
	retryer.tick = tickRetryer
	retryer.tick()

	return *retryer, nil
}

func newDefaultRetryer(actionHandler func(), errorHandler func()) (Retryer, error) {
	// 12 fois max, 10 seconde d'intervale, timeOut de 8 secondes
	return newRetryer(12, 10*time.Second, 8*time.Second, actionHandler, errorHandler)
}

func tickRetryer() {

}
