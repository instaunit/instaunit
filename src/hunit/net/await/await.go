package await

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const defaultRetry = time.Second

var (
	ErrTimeout     = fmt.Errorf("Timed out")
	ErrUnsupported = fmt.Errorf("Unsupported scheme")
)

func Await(cxt context.Context, ustr []string, timeout time.Duration) error {
	deps := make(chan struct{})
	errs := make(chan error, 1)

	urls := make([]*url.URL, len(ustr))
	for i, e := range ustr {
		u, err := url.Parse(e)
		if err != nil {
			return err
		}
		urls[i] = u
	}

	var cancel context.CancelFunc
	cxt, cancel = context.WithCancel(cxt)
	defer cancel()

	wg := &sync.WaitGroup{}
	go func() {
		for _, u := range urls {
			switch u.Scheme {
			case "file":
				wg.Add(1)
				go waitForFile(cxt, wg, u.Path, defaultRetry)
			case "tcp", "tcp4", "tcp6":
				wg.Add(1)
				go waitForSocket(cxt, wg, u.Scheme, u.Host, defaultRetry)
			case "unix":
				wg.Add(1)
				go waitForSocket(cxt, wg, u.Scheme, u.Path, defaultRetry)
			case "http", "https":
				wg.Add(1)
				go waitForHTTP(cxt, wg, u, defaultRetry)
			default:
				errs <- ErrUnsupported
				break
			}
		}
		wg.Wait()
		close(deps)
	}()

	select {
	case <-deps:
		break
	case err := <-errs:
		return err
	case <-time.After(timeout):
		return ErrTimeout
	}

	return nil
}

func waitForFile(cxt context.Context, wg *sync.WaitGroup, path string, retry time.Duration) {
	defer wg.Done()
	for {
		_, err := os.Stat(path)
		if err == nil {
			return
		} else if os.IsNotExist(err) {
			time.Sleep(retry)
		} else { // something else went wrong; just retry?
			time.Sleep(retry)
		}
	}
}

func waitForHTTP(cxt context.Context, wg *sync.WaitGroup, endpoint *url.URL, retry time.Duration) {
	defer wg.Done()
	client := &http.Client{
		Timeout: retry,
	}

	for {
		req, err := http.NewRequest("GET", endpoint.String(), nil)
		if err != nil {
			time.Sleep(retry)
		}

		resp, err := client.Do(req)
		if err != nil { // something else went wrong; just retry?
			time.Sleep(retry)
		} else if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return
		} else {
			time.Sleep(retry)
		}
	}
}

func waitForSocket(cxt context.Context, wg *sync.WaitGroup, scheme, addr string, retry time.Duration) {
	defer wg.Done()
	for {
		conn, err := net.DialTimeout(scheme, addr, retry)
		if err != nil {
			time.Sleep(retry)
		}
		if conn != nil {
			return
		}
	}
}
