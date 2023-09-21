package brightness

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
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

func getBrightnessType(filename string) (int, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return -1, err
	}

	value, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return -1, err
	}

	return value, nil
}

func getBrightness() (int, error) {
	backlight_path := "/sys/devices/pci0000:00/0000:00:08.1/0000:e3:00.0/backlight/amdgpu_bl0"
	//backlight_path := "/sys/devices/pci0000:00/0000:00:02.0/drm/card0/card0-eDP-1/intel_backlight"

	fMaxBrightness := fmt.Sprintf("%s/max_brightness", backlight_path)
	fCurrentBrightness := fmt.Sprintf("%s/actual_brightness", backlight_path)

	max, err := getBrightnessType(fMaxBrightness)
	if err != nil {
		return -1, err
	}

	current, err := getBrightnessType(fCurrentBrightness)
	if err != nil {
		return -1, err
	}

	value := float64(current) / float64(max) * 100

	return int(math.Ceil(value)), nil
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
