package pireceiver

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// current thermal zone id, type, and current temperature
type raspberryPiResult struct {
	id        int
	zone_type string
	temp      float64
}

type raspberryPi struct {
	// path t the raspberry pi's /sys/class/
	path   string
	logger *zap.Logger
}

func newRaspberryPi(path string, logger *zap.Logger) *raspberryPi {
	return &raspberryPi{
		path:   path,
		logger: logger,
	}
}

// Get Raspberry Pi thermal zone information from /sys/class/thermal/thermal_zone*
func (t *raspberryPi) get() ([]raspberryPiResult, error) {
	// get a list of thermal zones
	zoneList, err := t.getThermalZoneList()
	if err != nil {
		return nil, err
	}

	// get the temperature for each thermal zone
	zoneMap, err := t.getThermalZoneMap(zoneList)
	if err != nil {
		return nil, err
	}

	// convert the map to a slice
	zoneSlice := make([]raspberryPiResult, 0)
	for _, zone := range zoneMap {
		zoneSlice = append(zoneSlice, *zone)
	}

	return zoneSlice, nil
}

// getThermalZoneList returns a list of thermal zones
func (t *raspberryPi) getThermalZoneList() ([]int, error) {
	// get directories in /sys/class/thermal/
	dir, err := os.Open(t.path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// filter out directories that don't start with thermal_zone
	zoneList := make([]int, 0)
	for {
		// get the next directory
		fileInfos, err := dir.Readdir(1)
		// if EOF, break out of the loop
		if err != nil {
			if err.Error() != "EOF" {
				return nil, err
			}
			break
		}

		// check if there are no more directories
		if len(fileInfos) == 0 {
			break
		}

		// check if the directory starts with thermal_zone
		if !strings.HasPrefix(fileInfos[0].Name(), "thermal_zone") {
			continue
		}

		// convert the directory name to an int
		zoneID, err := strconv.Atoi(fileInfos[0].Name()[len("thermal_zone"):])
		if err != nil {
			return nil, err
		}

		zoneList = append(zoneList, zoneID)
	}

	return zoneList, nil
}

// getThermalZoneTemp returns a map of thermal zone id to thermal zone type and temperature
func (t *raspberryPi) getThermalZoneMap(zoneList []int) (map[int]*raspberryPiResult, error) {
	zoneMap := make(map[int]*raspberryPiResult)

	for _, id := range zoneList {
		zoneType, err := t.getThermalZoneType(id)
		if err != nil {
			return nil, err
		}

		zoneTemp, err := t.getThermalZoneTemp(id)
		if err != nil {
			return nil, err
		}

		zoneMap[id] = &raspberryPiResult{
			id:        id,
			zone_type: zoneType,
			temp:      zoneTemp,
		}
	}

	return zoneMap, nil
}

func (t *raspberryPi) getThermalZoneType(id int) (string, error) {
	// Construct the file path using filepath.Join for safety
	filePath := filepath.Join(t.path, fmt.Sprintf("thermal_zone%d", id), "type")

	// Read the entire file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read thermal zone type from %s: %w", filePath, err)
	}

	// Return the string content, trimming any new line character
	return strings.TrimSpace(string(data)), nil
}

func (t *raspberryPi) getThermalZoneTemp(id int) (float64, error) {
	// Construct the file path using filepath.Join for safety
	filePath := filepath.Join(t.path, fmt.Sprintf("thermal_zone%d", id), "temp")

	// Read the entire file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read thermal zone temperature from %s: %w", filePath, err)
	}

	// Convert the string data to an integer
	temp, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, fmt.Errorf("failed to convert temperature data to integer: %w", err)
	}

	// Convert the int to degrees Celsius and round to one decimal place
	tempFloat := float64(temp) / 1000.0
	return math.Round(tempFloat*10) / 10, nil
}
