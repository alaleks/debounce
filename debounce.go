package debounce

import (
	"context"
	"sync"
	"time"
)

// cache status
const (
	caching uint8 = 1 << iota
	noCaching
)

type CircuitType interface {
	any
}

type CircuitFunc[T CircuitType] func(ctx context.Context) (T, error)

// circuit: first call
type Circuit[T CircuitType] struct {
	Caching       uint8 // 1 - cache, 2 - isn`t cache
	Threshold     time.Duration
	mtx           *sync.Mutex
	thresholdTime time.Time
	circuitLc
}

// circuit: last call
type circuitLc struct {
	ticker time.Ticker
	once   sync.Once
}

// create a new circuit
func NewCircuit[T CircuitType](threshold uint) Circuit[T] {
	c := Circuit[T]{}
	c.Threshold = time.Duration(threshold) * time.Second
	c.mtx = &sync.Mutex{}
	c.ticker = *time.NewTicker(100 * time.Millisecond) // default ticker is 100 milliseconds
	return c
}

// cnahge time duration for ticker
func (c *Circuit[T]) ChangeTicker(tickerMillisecond uint) {
	c.ticker = *time.NewTicker(time.Duration(tickerMillisecond) * time.Millisecond)
}

// method debounce (first call)
// return any function with signature func(ctx context.Context) (T, error)
func (c *Circuit[T]) DebounceFirstCall(cfn CircuitFunc[T]) CircuitFunc[T] {
	var res T
	var err error

	return func(ctx context.Context) (T, error) {
		c.mtx.Lock()
		defer c.mtx.Unlock()

		if time.Now().Before(c.thresholdTime.Add(c.Threshold)) {
			c.Caching = caching
			return res, err
		}
		c.thresholdTime = time.Now()
		res, err = cfn(ctx)
		c.Caching = noCaching
		return res, err
	}
}

// method debounce (last call)
// return any function with signature func(ctx context.Context) (T, error)
func (c *Circuit[T]) DebounceLastCall(cfn CircuitFunc[T]) CircuitFunc[T] {
	var res T
	var err error

	return func(ctx context.Context) (T, error) {
		c.mtx.Lock()
		defer c.mtx.Unlock()

		c.once.Do(func() {
			go func() {
				defer func() {
					c.mtx.Lock()
					c.ticker.Stop()
					c.once = sync.Once{}
					c.mtx.Unlock()
				}()
				for {
					select {
					case <-c.ticker.C:
						c.Caching = caching
						c.mtx.Lock()
						if time.Now().After(c.thresholdTime.Add(c.Threshold)) {
							res, err = cfn(ctx)
							c.Caching = noCaching
							c.thresholdTime = time.Now()
							c.mtx.Unlock()
							return
						}
						c.mtx.Unlock()
					case <-ctx.Done():
						c.mtx.Lock()
						err = ctx.Err()
						c.mtx.Unlock()
						return
					}
				}
			}()
		})

		return res, err
	}
}
