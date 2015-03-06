package main

import (
	"fmt"
	"net/http"

	expb "github.com/odeke-em/exponential-backoff"
)

func tryGet(uri string) expb.Producer {
	return func() (interface{}, error) {
		return http.Get(uri)
	}
}

func consume(result interface{}, err error) {
	fmt.Printf("result: %v err: %v\n", result, err)
}

func newBackoff(url string, retryCount uint32) *expb.ExponentialBacker {
	req := tryGet(url)
	return &expb.ExponentialBacker{
		Do:          req,
		StatusCheck: httpStatus,
		RetryCount:  retryCount,
	}
}

func httpStatus(v interface{}) (ok, retryable bool) {
	res := v.(*http.Response)
	statusCode := res.StatusCode
	fmt.Println("statuscode", statusCode)
	if statusCode >= 200 && statusCode <= 299 {
		ok = true
		return
	}
	if statusCode == http.StatusForbidden {
		retryable = true
		return
	}
	return
}

func main() {
	backer := newBackoff("https://golang.org/pkg/net/httpx", 5)
	expb.ExponentialBackOff(backer, consume)
	fmt.Println("expb")
}
