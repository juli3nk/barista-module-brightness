package brightness

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getBacklightPath() (string, error) {
	basePath := "/sys/class/backlight/"

	entries, err := os.ReadDir(basePath)
	if err != nil || len(entries) == 0 {
		return "", fmt.Errorf("no backlight device detected")
	}

	return filepath.Join(basePath, entries[0].Name()), nil
}

func readIntFromFile(filename string) (int, error) {
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

func getBrightness(brightnessPath string, maxBrightness int) (int, error) {
	currentBrightness, err := readIntFromFile(filepath.Join(brightnessPath, "actual_brightness"))
	if err != nil {
		return -1, err
	}

	return int(math.Round(float64(currentBrightness) / float64(maxBrightness) * 100)), nil
}
