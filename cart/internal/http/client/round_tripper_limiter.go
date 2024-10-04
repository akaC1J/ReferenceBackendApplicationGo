package client

import (
	"net/http"
	"time"
)

type LimiterRoundTripper struct {
	transport    http.RoundTripper
	reqPerSecond uint
	sem          chan struct{}
}

func NewLimiterRoundTripper(roundTripper http.RoundTripper, reqPerSecond uint) http.RoundTripper {
	rt := &LimiterRoundTripper{
		transport:    roundTripper,
		reqPerSecond: reqPerSecond,
		sem:          make(chan struct{}, reqPerSecond),
	}

	for i := uint(0); i < reqPerSecond; i++ {
		rt.sem <- struct{}{}
	}
	go rt.refilling()
	return rt
}

func (rt *LimiterRoundTripper) refilling() {
	//нам нужен тикер во время жизни всего приложения, реализуем без стопа
	// иначе неочевидно что должно происходить при остановке
	tick := time.Tick(time.Second / time.Duration(rt.reqPerSecond))
	for range tick {
		select {
		case rt.sem <- struct{}{}:
		default:
		}
	}

}
func (rt *LimiterRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	<-rt.sem
	resp, err := rt.transport.RoundTrip(req)
	return resp, err
}
