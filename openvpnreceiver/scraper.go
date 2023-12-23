package openvpnreceiver

import (
	"context"
	"time"

	"github.com/drewby/openvpnreceiver/internal/metadata"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

const (
	receive  = "receive"
	transmit = "transmit"
)

type scraper struct {
	logger         *zap.Logger              // Logger to log events
	metricsBuilder *metadata.MetricsBuilder // MetricsBuilder to build metrics
	openvpnStatus  *openvpnStatus           // openvpnStatus to get stats from /var/log/openvpn/status.log
}

// newScraper is a constructor function which returns a new scraper instance
func newScraper(metricsBuilder *metadata.MetricsBuilder, path string, logger *zap.Logger) *scraper {
	return &scraper{
		logger:         logger,
		metricsBuilder: metricsBuilder,
		openvpnStatus:  newOpenvpnStatus(path, logger),
	}
}

func (s *scraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
	s.logger.Debug("Scraping OpenVPN status from path", zap.String("path", s.openvpnStatus.path))

	status, err := s.openvpnStatus.get()
	if err != nil {
		return pmetric.NewMetrics(), err
	}

	now := pcommon.NewTimestampFromTime(time.Now())

	// start connection count
	var connectionCount int64 = 0

	// Iterate over the stats and add them to the metrics builder
	for _, stat := range status {

		s.metricsBuilder.RecordOpenvpnBytesDataPoint(now, stat.BytesReceived, receive, stat.CommonName, stat.RealAddress, stat.RealPort)
		s.metricsBuilder.RecordOpenvpnBytesDataPoint(now, stat.BytesSent, transmit, stat.CommonName, stat.RealAddress, stat.RealPort)

		connectionCount++
	}

	s.metricsBuilder.RecordOpenvpnConnectionsDataPoint(now, connectionCount)

	metrics := s.metricsBuilder.Emit()

	return metrics, nil
}
