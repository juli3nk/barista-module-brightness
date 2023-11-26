package brightness

import (
	"time"

	"barista.run/bar"
	"barista.run/base/value"
	"barista.run/outputs"
	"barista.run/timing"
)

type Module struct {
	scheduler *timing.Scheduler
	outputFunc value.Value // of func(string) bar.Output
}

func New() *Module {
	m := &Module{scheduler: timing.NewScheduler()}
	m.RefreshInterval(3 * time.Second)

	m.outputFunc.Set(func(in string) bar.Output {
		return outputs.Text(in)
	})

	return m
}

func (m *Module) Output(outputFunc func(int) bar.Output) *Module {
	m.outputFunc.Set(outputFunc)
	return m
}

// RefreshInterval configures the polling frequency for getloadavg.
func (m *Module) RefreshInterval(interval time.Duration) *Module {
	m.scheduler.Every(interval)
	return m
}

func (m *Module) Stream(s bar.Sink) {
	outputFunc := m.outputFunc.Get().(func(int) bar.Output)

	nextOutputFunc, done := m.outputFunc.Subscribe()
	defer done()

	data, err := getBrightness()
	for {
		if s.Error(err) {
			return
		}

		s.Output(outputFunc(data))

		select {
		case <-m.scheduler.C:
			data, err = getBrightness()
		case <-nextOutputFunc:
			outputFunc = m.outputFunc.Get().(func(int) bar.Output)
		}
	}
}
