package openvpnreceiver

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Common Name,Real Address,Bytes Received,Bytes Sent,Connected Since
type openvpnStatusResult struct {
	CommonName    string
	RealAddress   string
	RealPort      int64
	BytesReceived int64
	BytesSent     int64
	// ConnectedSince is a date and time string to be parsed
	ConnectedSince time.Time
}

type openvpnStatus struct {
	path   string
	logger *zap.Logger
}

func newOpenvpnStatus(path string, logger *zap.Logger) *openvpnStatus {
	return &openvpnStatus{
		path:   path,
		logger: logger,
	}
}

// Get OpenVpn status from /var/log/openvon/status.log (or other path)
func (t *openvpnStatus) get() ([]openvpnStatusResult, error) {
	file, err := os.Open(t.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	statusMap, err := t.parseFile(file)
	if err != nil {
		return nil, err
	}

	status := make([]openvpnStatusResult, 0, len(statusMap))
	for _, stat := range statusMap {
		status = append(status, *stat)
	}

	return status, nil
}

// parseFile parses the OpenVPN status file and returns a map of
// common name to openvpnStatusResult
func (t *openvpnStatus) parseFile(file *os.File) (map[string]*openvpnStatusResult, error) {
	status := make(map[string]*openvpnStatusResult)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if isIgnoreLine(line) {
			continue
		}

		if isStopLine(line) {
			break
		}

		result, err := parseLine(line)
		if err != nil {
			t.logger.Warn(err.Error(), zap.String("line", line))
			continue
		}

		status[result.CommonName] = result
	}

	return status, scanner.Err()
}

func isIgnoreLine(line string) bool {
	switch {
	case
		strings.HasPrefix(line, "OpenVPN"),
		strings.HasPrefix(line, "Updated"),
		strings.HasPrefix(line, "Common Name"):
		return true
	}
	return false
}

func isStopLine(line string) bool {
	return strings.HasPrefix(line, "ROUTING TABLE")
}

func parseLine(line string) (*openvpnStatusResult, error) {
	fields := strings.Split(line, ",")
	if len(fields) != 5 {
		return nil, errors.New("Invalid client line")
	}

	realAddress, realPort, err := parseAddressAndPort(fields[1])
	if err != nil {
		return nil, err
	}

	bytesReceived, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return nil, errors.New("Invalid bytes received")
	}

	bytesSent, err := strconv.ParseInt(fields[3], 10, 64)
	if err != nil {
		return nil, errors.New("Invalid bytes sent")
	}

	connectedSince, err := time.Parse("2006-01-02 15:04:05", fields[4])
	if err != nil {
		return nil, errors.New("Invalid connected since")
	}

	return &openvpnStatusResult{
		CommonName:     fields[0],
		RealAddress:    realAddress,
		RealPort:       realPort,
		BytesReceived:  bytesReceived,
		BytesSent:      bytesSent,
		ConnectedSince: connectedSince,
	}, nil
}

func parseAddressAndPort(input string) (string, int64, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return "", 0, errors.New("Invalid real address and port")
	}

	port, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return input, 0, err
	}

	return parts[0], port, nil
}
