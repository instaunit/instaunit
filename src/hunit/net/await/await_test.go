package await

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	// "github.com/bww/go-util/debug"
	"github.com/stretchr/testify/assert"
)

func success(rsp http.ResponseWriter, req *http.Request) {
	rsp.WriteHeader(http.StatusOK)
}

func failure(rsp http.ResponseWriter, req *http.Request) {
	rsp.WriteHeader(http.StatusNotFound)
}

func TestMain(m *testing.M) {
	// debug.DumpRoutinesOnInterrupt()
	os.Exit(m.Run())
}

func TestAwaitHTTP(t *testing.T) {

	t.Run("Available", func(t *testing.T) {
		srv := &http.Server{
			Addr:         ":9999",
			Handler:      http.HandlerFunc(success),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		defer srv.Shutdown(context.TODO())

		go srv.ListenAndServe()
		err := Await(context.Background(), []string{"http://localhost:9999/"}, time.Minute)
		assert.Nil(t, err, fmt.Sprint(err))
	})

	t.Run("Wait2Seconds", func(t *testing.T) {
		srv := &http.Server{
			Addr:         ":9999",
			Handler:      http.HandlerFunc(success),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		defer srv.Shutdown(context.TODO())

		wait := time.Second * 2
		go func() {
			time.Sleep(wait)
			srv.ListenAndServe()
		}()

		start := time.Now()
		err := Await(context.Background(), []string{"http://localhost:9999/"}, time.Minute)
		assert.Nil(t, err, fmt.Sprint(err))
		waited := time.Now().Sub(start)
		assert.Equal(t, true, waited < wait+(defaultRetry*2))
	})

	t.Run("Timeout", func(t *testing.T) {
		err := Await(context.Background(), []string{"http://localhost:9999/"}, time.Second)
		if assert.NotNil(t, err, "Expected timeout") {
			assert.Equal(t, ErrTimeout, err, fmt.Sprint(err))
		}
	})

	t.Run("Cancel", func(t *testing.T) {
		srv := &http.Server{
			Addr:         ":9999",
			Handler:      http.HandlerFunc(success),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		defer srv.Shutdown(context.TODO())

		wait := time.Second * 2
		cxt, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(wait)
			srv.ListenAndServe()
		}()

		start := time.Now()
		go func() {
			time.Sleep(wait / 2)
			cancel()
		}()

		err := Await(cxt, []string{"http://localhost:9999/"}, time.Minute)
		assert.Nil(t, err, fmt.Sprint(err))
		waited := time.Now().Sub(start)
		assert.Equal(t, true, waited < (wait/2)+defaultRetry)
	})

}

func TestAwaitSocket(t *testing.T) {

	t.Run("Available", func(t *testing.T) {
		srv, err := net.Listen("tcp", ":9999")
		assert.Nil(t, err, fmt.Sprint(err))
		defer srv.Close()

		go func() {
			for {
				_, err := srv.Accept()
				if err != nil {
					return
				}
			}
		}()

		err = Await(context.Background(), []string{"tcp4://localhost:9999/"}, time.Minute)
		assert.Nil(t, err, fmt.Sprint(err))
	})

	t.Run("Wait2Seconds", func(t *testing.T) {
		var srv net.Listener

		wait := time.Second * 2
		go func() {
			time.Sleep(wait)
			var err error
			srv, err = net.Listen("tcp", ":9999")
			assert.Nil(t, err, fmt.Sprint(err))
			for {
				_, err := srv.Accept()
				if err != nil {
					return
				}
			}
		}()

		start := time.Now()
		err := Await(context.Background(), []string{"tcp4://localhost:9999/"}, time.Minute)
		assert.Nil(t, err, fmt.Sprint(err))
		waited := time.Now().Sub(start)
		assert.Equal(t, true, waited < wait+(defaultRetry*2))

		srv.Close()
	})

	t.Run("Timeout", func(t *testing.T) {
		err := Await(context.Background(), []string{"tcp4://localhost:9999/"}, time.Second)
		if assert.NotNil(t, err, "Expected timeout") {
			assert.Equal(t, ErrTimeout, err, fmt.Sprint(err))
		}
	})

	t.Run("Cancel", func(t *testing.T) {
		cxt, cancel := context.WithCancel(context.Background())
		wait := time.Second

		start := time.Now()
		go func() {
			time.Sleep(wait / 2)
			cancel()
		}()

		err := Await(cxt, []string{"http://localhost:9999/"}, time.Minute)
		assert.Nil(t, err, fmt.Sprint(err))
		waited := time.Now().Sub(start)
		assert.Equal(t, true, waited < (wait/2)+defaultRetry)
	})

}

func TestAwaitFile(t *testing.T) {

	cleanup := func(dir string) {

	}

	t.Run("Available", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "await_")
		assert.Nil(t, err, fmt.Sprint(err))
		defer cleanup(dir)

		p := path.Join(dir, "test1")
		_, err = os.OpenFile(p, os.O_RDONLY|os.O_CREATE, 0666)
		assert.Nil(t, err, fmt.Sprint(err))

		err = Await(context.Background(), []string{fmt.Sprintf("file://%s", p)}, time.Minute)
		assert.Nil(t, err, fmt.Sprint(err))
	})

	t.Run("Wait2Seconds", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "await_")
		assert.Nil(t, err, fmt.Sprint(err))
		defer cleanup(dir)

		wait := time.Second * 2
		p := path.Join(dir, "test2")
		go func() {
			time.Sleep(wait)
			_, err := os.OpenFile(p, os.O_RDONLY|os.O_CREATE, 0666)
			assert.Nil(t, err, fmt.Sprint(err))
		}()

		start := time.Now()
		err = Await(context.Background(), []string{fmt.Sprintf("file://%s", p)}, time.Minute)
		assert.Nil(t, err, fmt.Sprint(err))
		waited := time.Now().Sub(start)
		assert.Equal(t, true, waited < wait+(defaultRetry*2))
	})

	t.Run("Timeout", func(t *testing.T) {
		err := Await(context.Background(), []string{"file:///does/not/exist"}, time.Second)
		if assert.NotNil(t, err, "Expected timeout") {
			assert.Equal(t, ErrTimeout, err, fmt.Sprint(err))
		}
	})

	t.Run("Cancel", func(t *testing.T) {
		cxt, cancel := context.WithCancel(context.Background())
		wait := time.Second

		start := time.Now()
		go func() {
			time.Sleep(wait / 2)
			cancel()
		}()

		err := Await(cxt, []string{"file:///does/not/exist"}, time.Minute)
		assert.Nil(t, err, fmt.Sprint(err))
		waited := time.Now().Sub(start)
		assert.Equal(t, true, waited < (wait/2)+defaultRetry)
	})

}
