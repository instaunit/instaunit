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

	"github.com/bww/go-util/debug"
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
	defer cancel() // make sure we stop waiting on return

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

	if timeout > 0 {
		select {
		case <-deps:
			break
		case err := <-errs:
			return err
		case <-time.After(timeout):
			return ErrTimeout
		}
	} else {
		select {
		case <-deps:
			break
		case err := <-errs:
			return err
		}
	}

	return nil
}

func waitForFile(cxt context.Context, wg *sync.WaitGroup, path string, retry time.Duration) {
	defer wg.Done()
	for {
		select {
		case <-cxt.Done():
			return
		default:
			// ... continue
		}
		_, err := os.Stat(path)
		if err == nil {
			return
		} else if os.IsNotExist(err) {
			log("await: No such file: %v", path)
			time.Sleep(retry)
		} else { // something else went wrong; just retry?
			log("await: Could not stat: %v: %v", path, err)
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
		select {
		case <-cxt.Done():
			return
		default:
			// ... continue
		}

		req, err := http.NewRequest("GET", endpoint.String(), nil)
		if err != nil {
			time.Sleep(retry)
		}

		rsp, err := client.Do(req.WithContext(cxt))
		if err != nil {
			log("await: Could not connect: %v", endpoint)
			time.Sleep(retry)
		} else if rsp.StatusCode >= 200 && rsp.StatusCode < 300 {
			return // ok, connected
		} else {
			log("await: Unexpected status: %v -> %v", endpoint, rsp.Status)
			time.Sleep(retry)
		}
	}
}

func waitForSocket(cxt context.Context, wg *sync.WaitGroup, scheme, addr string, retry time.Duration) {
	defer wg.Done()
	dialer := net.Dialer{Timeout: retry}
	for {
		select {
		case <-cxt.Done():
			return
		default:
			// ... continue
		}

		conn, err := dialer.DialContext(cxt, scheme, addr)
		if err != nil {
			time.Sleep(retry)
			log("await: Could not connect: %v", addr)
		} else if conn != nil {
			return
		}
	}
}

func log(f string, a ...interface{}) (int, error) {
	if debug.VERBOSE {
		if l := len(f); l < 1 || f[l-1] != '\n' {
			return fmt.Printf(f+"\n", a...)
		} else {
			return fmt.Printf(f, a...)
		}
	} else {
		return 0, nil
	}
}
