package debounce

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func getHTTP(ctx context.Context) (*http.Response, error) {
	return http.Get("https://google.com")
}

func TestDebounceFirstCall(t *testing.T) {
	cfc := NewCircuit[*http.Response](2)
	ctx := context.Background()
	dcfc := cfc.DebounceFirstCall(getHTTP)
	for i := 1; i < 11; i++ {
		_, _ = dcfc(ctx)
		switch i {
		case 1:
			if cfc.Caching != noCaching {
				t.Errorf("cfc.Caching should be 2 because this is first call")
			}
		case 2:
			if cfc.Caching != caching {
				t.Errorf("cfc.Caching should be 1 because this is result from caching")
			}
		}
	}
}

func TestDebounceLastCall(t *testing.T) {
	clc := NewCircuit[*http.Response](2)
	ctx2 := context.Background()
	dclc := clc.DebounceLastCall(getHTTP)
	for i := 1; i < 11; i++ {

		_, _ = dclc(ctx2)

		switch i {
		case 5:
			if clc.Caching != noCaching {
				t.Errorf("cfc.Caching should be 2")
			}
		case 1:
			if clc.Caching != 0 {
				t.Errorf("cfc.Caching should be 0")
			}

		}
		if i == 4 {
			time.Sleep(2 * time.Second)
		}
	}
}
