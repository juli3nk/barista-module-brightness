package brightness

import (
  "fmt"
  "io/ioutil"
  "math"
  "os"
  "runtime"
  "strconv"
  "strings"
)

func getKernel() (string, error) {
  var cpu string
  if runtime.GOARCH == "amd64" {
    cpu = "amd"
	}
  if runtime.GOARCH == "386" {
    cpu = "intel"
  }

  files, err := os.ReadDir("/sys/class/backlight")
  if err != nil {
    return "", err
  }

  for _, file := range files {
    if strings.Contains(file.Name(), cpu) {
      return file.Name(), nil
    }
  }

  return "", nil
}

func readValue(kernel, name string) (int, error) {
  file := fmt.Sprintf("/sys/class/backlight/%s/%s", kernel, name)

  bytes, err := ioutil.ReadFile(file)
  if err != nil {
    return 0, err
  }

  dat := strings.TrimSpace(string(bytes))

  n, err := strconv.Atoi(dat)
  if err != nil {
    return 0, err
  }
  return n, nil
}

// Fraction returns the brightness as a fraction of the maximum value.
func fraction(max, bri int) float64 {
  if max == 0 {
    return 0
  }

  return float64(bri) / float64(max)
}

// Percent returns the brightness in percent of the maximum value.
func percent(max, bri int) int {
  return int(math.Round(fraction(max, bri) * 100.0))
}

// Get updates the Bri and Max values after reading the respective files.
func getBrightness() (int, error) {
  kernel, err := getKernel()
  if err != nil {
    return -1, err
  }

  max, err := readValue(kernel, "max_brightness")
  if err != nil {
    return -1, err
  }

  bri, err := readValue(kernel, "actual_brightness")
  if err != nil {
    return -1, err
  }

  value := percent(max, bri)

  return value, nil
}
