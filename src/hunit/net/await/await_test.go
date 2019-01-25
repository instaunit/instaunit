package await

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/bww/go-util/debug"
	"github.com/stretchr/testify/assert"
)

func success(rsp http.ResponseWriter, req *http.Request) {
	rsp.WriteHeader(http.StatusOK)
}

func failure(rsp http.ResponseWriter, req *http.Request) {
	rsp.WriteHeader(http.StatusNotFound)
}

func TestMain(m *testing.M) {
	debug.DumpRoutinesOnInterrupt()
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

		go srv.ListenAndServe()
		err := Deps([]string{"http://localhost:9999/"}, time.Minute)
		srv.Shutdown(context.TODO())
		assert.Nil(t, err, fmt.Sprint(err))
	})

	t.Run("Wait2Seconds", func(t *testing.T) {
		srv := &http.Server{
			Addr:         ":9999",
			Handler:      http.HandlerFunc(success),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		wait := time.Second * 2
		go func() {
			time.Sleep(wait)
			srv.ListenAndServe()
		}()

		start := time.Now()
		err := Deps([]string{"http://localhost:9999/"}, time.Minute)
		srv.Shutdown(context.TODO())
		assert.Nil(t, err, fmt.Sprint(err))
		waited := time.Now().Sub(start)
		assert.Equal(t, true, waited < wait+(defaultRetry*2))
	})

	t.Run("Timeout", func(t *testing.T) {
		srv := &http.Server{
			Addr:         ":9999",
			Handler:      http.HandlerFunc(success),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		wait := time.Second * 2
		go func() {
			time.Sleep(wait)
			srv.ListenAndServe()
		}()

		err := Deps([]string{"http://localhost:9999/"}, time.Second)
		srv.Shutdown(context.TODO())
		if assert.NotNil(t, err, "Expected timeout") {
			assert.Equal(t, ErrTimeout, err, fmt.Sprint(err))
		}
	})

}
