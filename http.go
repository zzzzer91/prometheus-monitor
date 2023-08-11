package monitor

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(port uint16, g prometheus.Gatherer, opts ...ConfigOption) {
	c := newHttpConfig()
	c.apply(opts...)

	// By default, pprof does not track information about blocks and mutexes.
	// If you want to see this information, you need to enable it manually.
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	mux := http.NewServeMux()
	mux.Handle(c.path, promhttp.HandlerFor(g, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}))
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
}
