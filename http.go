package monitor

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(g prometheus.Gatherer, opts ...httpConfigOption) {
	c := newHttpConfig()
	c.apply(opts...)

	// 默认情况下 pprof 是不追踪block和mutex的信息的，如果想要看这两个信息，需要手动开启
	runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪，block
	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪，mutex

	http.Handle(c.path, promhttp.HandlerFor(g, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}))
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", c.port), nil); err != nil {
			panic(err)
		}
	}()
}
