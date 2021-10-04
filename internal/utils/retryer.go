package utils

import "time"

type Retryer struct {
	maxRetryTimes int
	tickInterval  time.Duration
	timeout       time.Duration
	tick          func()
	retryAction   func() (bool, error)
	errorHandler  func()
}

func newRetryer(maxRetryTimes int, tickInterval time.Duration, timeout time.Duration, retryAction func() (bool, error), errorHandler func()) *Retryer {
	retryer := new(Retryer)
	retryer.maxRetryTimes = maxRetryTimes
	retryer.tickInterval = tickInterval
	retryer.timeout = timeout
	retryer.retryAction = retryAction
	retryer.errorHandler = errorHandler
	retryer.tick = func() {
		for i := 0; i < maxRetryTimes; i++ {
			time.Sleep(tickInterval)
			state, err := retryAction()
			if err != nil {
				errorHandler()
			}
			if state {
				return
			}
		}
	}
	return retryer
}

func newDefaultRetryer(actionHandler func() (bool, error), errorHandler func()) *Retryer {
	// 12 fois max, 10 secondes d'intervale, timeOut de 8 secondes
	return newRetryer(12, 10*time.Second, 8*time.Second, actionHandler, errorHandler)
}
