package helpers

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetFreePort(t *testing.T) {
	tests := []struct {
		count int
	}{
		{
			count: 10,
		},
		{
			count: 100,
		},
		{
			count: 500,
		},
	}

	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.count), func(t *testing.T) {
			wg := sync.WaitGroup{}
			wg.Add(tt.count)
			ch := make(chan struct{})
			// p := GetFreePort(t)
			for i := 0; i < tt.count; i++ {
				go func() {
					defer wg.Done()
					<-ch
					p := GetFreePort(t)

					// NOTE: simply net.Listen() and net.Close() was not yielding expected test results.
					l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", p))
					if !assert.NoError(t, err) {
						return
					}
					s := httptest.Server{
						Listener: l,
						Config: &http.Server{
							Addr: fmt.Sprintf("localhost:%d", p),
							Handler: http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
							}),
							ReadHeaderTimeout: 10 * time.Second,
						},
					}
					s.Start()
					defer s.Close()

					resp, err := http.Get(s.URL)
					if !assert.NoError(t, err) {
						return
					}
					assert.NoError(t, resp.Body.Close())
				}()
			}
			close(ch)
			wg.Wait()
		})
	}
}
