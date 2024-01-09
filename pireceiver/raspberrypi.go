package pireceiver

import (
	"os"
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

type rasperryPi struct {
	// path t the raspberry pi's /sys/class/
	path   string
	logger *zap.Logger
}

func newRaspberryPi(path string, logger *zap.Logger) *rasperryPi {
	return &rasperryPi{
		path:   path,
		logger: logger,
	}
}

// Get Raspberry Pi thermal zone information from /sys/class/thermal/thermal_zone*
func (t *rasperryPi) get() ([]raspberryPiResult, error) {
	// loop over /sys/class/thermal/thermal_zone* and get the current temperature
	// id is the suffix of thermal_zone*
	// type is the first line of /sys/class/thermal/thermal_zone*/type
	// temp is the first line of /sys/class/thermal/thermal_zone*/temp
	// temp is in millidegrees Celsius, convert to degrees Celsius round to one decimal place
	// https://www.kernel.org/doc/Documentation/thermal/sysfs-api.txt

	// get a list of thermal zones
	zoneList, err := t.getThermalZoneList()
	if err != nil {
		return nil, err
	}

	// get the temperature for each thermal zone
	zoneMap, err := t.getThermalZoneTemp(zoneList)
	if err != nil {
		return nil, err
	}

	// convert the map to a slice
}

// getThermalZoneList returns a list of thermal zones
func (t *rasperryPi) getThermalZoneList() ([]int, error) {
	// get directories in /sys/class/thermal/
	// filter out directories that don't start with thermal_zone
	// convert the remaining directories to ints
	// return the list of ints

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
		if err != nil {
			return nil, err
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
func (t *rasperryPi) getThermalZoneTemp(zoneList []int) (map[int]*raspberryPiResult, error) {
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
