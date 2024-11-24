package bcmetrics

import (
	"golang.org/x/sys/unix"
	"testing"
	"time"
)

const unit string = "ns_cpu/op"

func now() time.Duration {
	var ts unix.Timespec
	if err := unix.ClockGettime(unix.CLOCK_PROCESS_CPUTIME_ID, &ts); err != nil {
		panic(err)
	}
	return time.Duration(ts.Nano())
}

type Metric interface {
	Start()
	Stop()
	Reset()
	Report(b *testing.B)
}

type CpuTime struct {
	on       bool
	start    time.Duration
	duration time.Duration
}

func (c *CpuTime) Start() {
	if !c.on {
		c.on = true
		c.start = now()
	}
}

func (c *CpuTime) Stop() {
	if c.on {
		c.duration += now() - c.start
		c.on = false
	}
}

func (c *CpuTime) Reset() {
	if c.on {
		c.start = now()
	}
	c.duration = 0
}

func (c *CpuTime) Report(b *testing.B) {
	b.ReportMetric(float64(c.duration.Nanoseconds())/float64(b.N), unit)
}

type Collector struct {
	metrics []Metric
}

func (c *Collector) AddMetrics(metrics ...Metric) {
	c.metrics = append(c.metrics, metrics...)
}

func (c *Collector) Start() {
	for _, g := range c.metrics {
		g.Start()
	}
}

func (c *Collector) Stop() {
	for _, g := range c.metrics {
		g.Stop()
	}
}

func (c *Collector) Reset() {
	for _, g := range c.metrics {
		g.Reset()
	}
}

func (c *Collector) Report(b *testing.B) {
	for _, m := range c.metrics {
		m.Report(b)
	}
}
