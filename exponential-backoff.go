package expb

import (
	"fmt"
	"math/rand"
	"time"
)

type StatusChecker func(q interface{}) (ok, retryable bool)

type ExponentialBacker struct {
	Debug       bool
	Do          Producer
	RetryCount  uint32
	StatusCheck StatusChecker
}

type Callback func(result interface{}, err error)
type Producer func() (result interface{}, err error)

func ExponentialBackOff(bk *ExponentialBacker, cb Callback) {
	retries := uint32(0)

	for {
		res, err := bk.Do()
		ok, retryable := bk.StatusCheck(res)

		if ok || !retryable || retries >= bk.RetryCount {
			cb(res, err)
			return
		}

		ms := time.Duration(1e9 * rand.Float64()) + ((1 << retries) * time.Second)
		if bk.Debug {
			fmt.Printf("trying again in %v\n", ms)
		}

		duration := time.Duration(ms)
		time.Sleep(duration)

		retries += 1
	}
}
