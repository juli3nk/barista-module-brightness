package brightness

import (
	"fmt"
	"path/filepath"

	"github.com/barista-run/barista/bar"
	"github.com/barista-run/barista/base/value"
	"github.com/barista-run/barista/outputs"
	"github.com/fsnotify/fsnotify"
)

type Module struct {
	outputFunc  value.Value // of func(int) bar.Output
	brightnessPath string
	maxBrightness  int
}

func New() *Module {
	m := &Module{}

	path, err := getBacklightPath()
	if err != nil {
		fmt.Println("Erreur : impossible de détecter le périphérique de rétroéclairage.")
		return nil
	}
	m.brightnessPath = path

	m.maxBrightness, err = readIntFromFile(filepath.Join(m.brightnessPath, "max_brightness"))
	if err != nil {
		fmt.Println("Erreur : impossible de lire max_brightness.")
		return nil
	}

	m.outputFunc.Set(func(in int) bar.Output {
		return outputs.Textf("%d", in)
	})

	return m
}

func (m *Module) Output(outputFunc func(int) bar.Output) *Module {
	m.outputFunc.Set(outputFunc)
	return m
}

func (m *Module) Stream(s bar.Sink) {
	outputFunc := m.outputFunc.Get().(func(int) bar.Output)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		s.Error(err)
		return
	}
	defer watcher.Close()

	brightnessFile := filepath.Join(m.brightnessPath, "actual_brightness")
	err = watcher.Add(brightnessFile)
	if err != nil {
		s.Error(err)
		return
	}

	data, err := getBrightness(m.brightnessPath, m.maxBrightness)
	if err != nil {
		s.Error(err)
		return
	}
	s.Output(outputFunc(data))

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				data, err = getBrightness(m.brightnessPath, m.maxBrightness)
				if err != nil {
					s.Error(err)
					return
				}
				s.Output(outputFunc(data))
			}
		case err := <-watcher.Errors:
			s.Error(err)
			return
		}
	}
}
