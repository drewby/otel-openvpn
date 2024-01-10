package pireceiver

import (
	"testing"
)

func TestGetCPUTemperature(t *testing.T) {
	expected := []raspberryPiResult{
		{
			id:        0,
			zone_type: "x86_pkg_temp",
			temp:      49.0,
		},
		{
			id:        1,
			zone_type: "x86_pkg_temp",
			temp:      54.1,
		},
		{
			id:        2,
			zone_type: "x86_pkg_temp",
			temp:      62.3,
		},
	}

	// create a new Raspberry Pi
	r := newRaspberryPi("testdata", nil)

	// get thermal zone temperature information
	zones, err := r.get()
	if err != nil {
		t.Errorf("error getting CPU temperature: %v", err)
	}

	// compare the results
	for i, zone := range zones {
		if zone.zone_type != expected[zone.id].zone_type {
			t.Errorf("expected zone_type %s, got %s", expected[i].zone_type, zone.zone_type)
		}
		if zone.temp != expected[zone.id].temp {
			t.Errorf("expected temp %f, got %f", expected[i].temp, zone.temp)
		}
	}
}
