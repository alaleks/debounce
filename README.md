# debounce

**debounce** is very simple package that implements corresponding pattern.

## installation

```
go get github.com/alaleks/debounce
```

## usage by first call

```
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
}
```

## usage by last call
```
func getHTTP(ctx context.Context) (*http.Response, error) {
	return http.Get("https://google.com")
}

func main() {

	clc := debounce.NewCircuit[*http.Response](1)
	ctx := context.Background()
	dclc := clc.DebounceLastCall(getHTTP)
	for i := 1; i < 11; i++ {
		res, err := dclc(ctx)
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

```