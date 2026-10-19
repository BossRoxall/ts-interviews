// Command loadgen sends POST /v1/check requests to the rate limiter, each from
// a different key, and prints a running tally of response codes.
//
//	go run ./tools/loadgen -url http://localhost:8080/v1/check -keys 200000
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	url := flag.String("url", "http://localhost:8080/v1/check", "check endpoint")
	keys := flag.Int("keys", 200000, "number of distinct keys to send (one request each)")
	concurrency := flag.Int("concurrency", 32, "concurrent requests")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := &http.Client{Timeout: 5 * time.Second}
	var (
		next   atomic.Int64
		mu     sync.Mutex
		counts = map[string]int{}
		wg     sync.WaitGroup
	)
	record := func(code string) {
		mu.Lock()
		counts[code]++
		mu.Unlock()
	}

	start := time.Now()
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				mu.Lock()
				fmt.Printf("%5.0fs  sent=%d  %s\n", time.Since(start).Seconds(), next.Load(), summary(counts))
				mu.Unlock()
			case <-done:
				return
			}
		}
	}()

	for range *concurrency {
		wg.Go(func() {
			for {
				i := next.Add(1)
				if i > int64(*keys) || ctx.Err() != nil {
					return
				}
				body := fmt.Sprintf(`{"key": "loadgen-user-%012d-%s"}`, i, start.Format("150405"))
				req, _ := http.NewRequestWithContext(ctx, http.MethodPost, *url, bytes.NewBufferString(body))
				resp, err := client.Do(req)
				if err != nil {
					record("error")
					continue
				}
				resp.Body.Close()
				record(fmt.Sprint(resp.StatusCode))
			}
		})
	}
	wg.Wait()
	close(done)
	fmt.Printf("done in %s  %s\n", time.Since(start).Round(time.Second), summary(counts))
}

func summary(counts map[string]int) string {
	codes := make([]string, 0, len(counts))
	for code := range counts {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	parts := make([]string, 0, len(codes))
	for _, code := range codes {
		parts = append(parts, fmt.Sprintf("%s=%d", code, counts[code]))
	}
	return strings.Join(parts, "  ")
}
