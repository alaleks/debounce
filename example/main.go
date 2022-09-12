package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alaleks/debounce"
)

func getHTTP(ctx context.Context) (*http.Response, error) {
	return http.Get("https://google.com")
}

func main() {

	cfc := debounce.NewCircuit[*http.Response](2)
	ctx := context.Background()
	dcfc := cfc.DebounceFirstCall(getHTTP)
	for i := 1; i < 11; i++ {
		res, err := dcfc(ctx)
		if err != nil {
			fmt.Println(i, err, cfc.Caching)
		} else {
			fmt.Println(i, res.Status, cfc.Caching)
		}
	}

	clc := debounce.NewCircuit[*http.Response](1)
	ctx2 := context.Background()
	dclc := clc.DebounceLastCall(getHTTP)
	for i := 1; i < 11; i++ {
		res, err := dclc(ctx2)
		if err != nil {
			fmt.Println(i, err, clc.Caching)
		} else {
			fmt.Println(i, res, clc.Caching)
		}
		if i == 4 {
			time.Sleep(1 * time.Second)
		}

	}
}
