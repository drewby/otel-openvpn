package openvpnreceiver

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestGetOpenvpnStatus(t *testing.T) {
	expected := []openvpnStatusResult{
		{
			CommonName:     "iPad",
			RealAddress:    "192.168.140.245",
			RealPort:       62741,
			BytesReceived:  165982,
			BytesSent:      43078,
			ConnectedSince: time.Date(2023, 12, 22, 7, 11, 6, 0, time.UTC),
		},
		{
			CommonName:     "iOS",
			RealAddress:    "192.168.140.245",
			RealPort:       62527,
			BytesReceived:  16836175,
			BytesSent:      556816484,
			ConnectedSince: time.Date(2023, 12, 22, 7, 1, 22, 0, time.UTC),
		},
	}

	f := filepath.Join("testdata", "openvpn.log")
	s := newOpenvpnStatus(f, zap.NewNop())

	status, err := s.get()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(status, expected) {
		t.Errorf("Expected: %+v, but got %+v", expected, status)
	}
}
